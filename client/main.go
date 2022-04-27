package main

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/fukaraca/gRPC-hello-Mars/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("dialing failed", err)
	}
	defer conn.Close()
	client := pb.NewMessagingServiceClient(conn)

	fmt.Println("Type your name:")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	readerFunc := func() {
		for {
			fmt.Println(">>")
			msg, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln("msg couldn't be read", err)
			}
			if msg == "clear" {
				CallClear()
			}
			temp := &pb.MessageRequest{
				Text:        msg,
				Sender:      name,
				SendingTime: timestamppb.Now(),
			}
			newReq := &pb.CreateMessageRequest{Request: temp}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			resp, err := client.SendMessage(ctx, newReq)
			if err != nil {
				log.Fatalln("sending message failed from client", err)
			}
			if resp.Response.Read == true {
				log.Println("message read")
			}

		}
	}
	go readerFunc()
	chann := make(chan struct{})
	<-chann

}

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

//clear terminal
func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		log.Fatalln("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
