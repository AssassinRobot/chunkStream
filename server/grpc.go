package main

import (
	"chunk/file"
	uploadpb "chunk/gen"
	"io"
	"log"
	"path/filepath"
	"runtime"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	uploadpb.UnimplementedUploadServiceServer
}

const uploadDir = "tmp"

func (s *Server) Upload(stream uploadpb.UploadService_UploadServer) error {
	file := file.NewServerFile()
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)

		}
	}()

	var fileName string
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fileName = md["file_name"][0]
		log.Printf("Receiving file: %s\n", fileName)
	}

	var fileSize int64
	fileSize = 0

	for {
		req, err := stream.Recv()
		if err != nil {
			log.Println("error:", err)
			if err == io.EOF {
				break
			}
			return status.Error(codes.Internal, err.Error())
		}

		err = file.SetFile(fileName,GetUploadDir())

		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		chunk := req.GetChunk()
		fileSize += int64(len(chunk))

		log.Printf("received a chunk with size: %d\n", fileSize)

		if err := file.Write(chunk); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	log.Printf("saved file: %s, size: %d\n", fileName, fileSize)

	return stream.SendAndClose(&uploadpb.UploadResponse{Name: fileName, Size: fileSize})
}

func GetUploadDir() string {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Unable to get current file")
	}

	return filepath.Dir(f)+"/"+uploadDir
}
