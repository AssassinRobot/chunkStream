package main

import (
	uploadpb "chunk/gen"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	RunServer()
}
const port = "8080"
func RunServer() {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	uploadpb.RegisterUploadServiceServer(grpcServer, &Server{})

	log.Printf("server is going to run on port %v.......",port)
	
	log.Fatalf("Error running server:%v",grpcServer.Serve(lis))
}
