package main

import (
	"chunk/file"
	uploadpb "chunk/gen"
	"context"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	// 1MB
	chunkSize = 1048576
)

func main() {
	file := file.NewFile()
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)

		}
	}()

	fileName := GetFilePathByArgs()
	if fileName == "" {
		log.Fatalln("File name is required.")
	}

	err := file.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file '%s': %v", fileName, err)
	}

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := uploadpb.NewUploadServiceClient(conn)

	parentCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	md := metadata.Pairs("file_name", file.GetFileName())
	ctx := metadata.NewOutgoingContext(parentCtx, md)

	stream, err := client.Upload(ctx)
	if err != nil {
		log.Fatalf("Failed to create upload stream: %v", err)
	}

	buf := make([]byte, chunkSize)
	batchNumber := 1

	for {
		num, readErr := file.Read(buf)
		if readErr != nil {
			if readErr == io.EOF {
				log.Println("Reached end of file.")
				break
			}
			log.Fatalf("Error reading file: %v", readErr)
		}

		chunk := buf[:num]

		sendErr := stream.Send(&uploadpb.UploadRequest{Chunk: chunk})
		if sendErr != nil {
			log.Fatalf("Error sending chunk #%d: %v", batchNumber, sendErr)
		}

		log.Printf("Sent - batch #%d - size: %d bytes\n", batchNumber, len(chunk))
		batchNumber++
	}

	res, resErr := stream.CloseAndRecv()
	if resErr != nil {
		log.Fatalf("Error closing stream and receiving response: %v", resErr)
	}

	if res.GetSize() != file.GetFileSize() {
		log.Fatalf("Expected file size: %d, but got: %d", file.GetFileSize(), res.GetSize())
	}

	log.Printf("Upload completed - total size: %d bytes - file name: %s\n", res.GetSize(), res.GetName())
}

func GetFilePathByArgs() string {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: go run client/main.go <directory/file>")
	}
	return os.Args[1]
}
