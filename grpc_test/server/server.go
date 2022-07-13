package main

import (
	"awesomeProject/grpc_stream/proto"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const PROT = ":50052"

type server struct {
	proto.UnimplementedGreeterServer
}

func (s *server) GetStream(req *proto.StreamReqData, res proto.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++
		_ = res.Send(&proto.StreamResData{
			Data: fmt.Sprintf("%v", time.Now().Unix()),
		})
		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}
	return status.Errorf(codes.OK, "method GetStream ended")
}

func (s *server) PostStream(cliStr proto.Greeter_PostStreamServer) error {
	for {
		a, err := cliStr.Recv()
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println(a.Data)
		}
	}
	return status.Errorf(codes.OK, "method GetStream ended")
}

func (s *server) AllStream(allStr proto.Greeter_AllStreamServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			data, _ := allStr.Recv()
			fmt.Println("Stream from client: " + data.Data)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			allStr.Send(&proto.StreamResData{
				Data: "data from Server",
			})
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
	return status.Errorf(codes.OK, "method GetStream ended")
}

func main() {
	// instantiate server obj
	lis, err := net.Listen("tcp", PROT)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}
