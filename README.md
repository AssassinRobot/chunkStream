# gRPC File Streaming Example

This repository contains a sample application demonstrating how to stream files using gRPC. It includes both a client and a server implementation that allows you to send file chunks over a gRPC connection.

## Usage

1. Clone the repository:
   ```bash
   git clone https://github.com/AssassinRobot/chunkStream.git
   cd chunkStream
   ```
2. Install dependencies
   ```bash
   go mod tidy
   ```
3. Run server and client:
   ```bash
   make run_server
   ```
   open another terminal and press:
   ```bash
   make run_client test2.txt
   ```
   note that test2.txt is that file which is in repository and we send it to server

   
