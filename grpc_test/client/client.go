package main

import (
	"awesomeProject/grpc_stream/proto"
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)

	// server stream mode
	res, _ := c.GetStream(context.Background(), &proto.StreamReqData{Data: "Muke"})
	for {
		a, err := res.Recv()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(a.Data)
	}

	// client stream mode
	putS, _ := c.PostStream(context.Background())
	for i := 0; i < 10; i++ {
		_ = putS.Send(&proto.StreamReqData{
			Data: fmt.Sprintf("Watler ping %d time", i),
		})
		time.Sleep(time.Second)
	}

	// all stream mode
	allStr, _ := c.AllStream(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			data, _ := allStr.Recv()
			fmt.Println("Stream from Server: " + data.Data)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			allStr.Send(&proto.StreamReqData{
				Data: "data from Client",
			})
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
}
