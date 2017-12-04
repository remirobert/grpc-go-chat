package main

import (
	pb "grpc-go-chat/chat"
	"log"
	"google.golang.org/grpc"
	"context"
	"strconv"
	"os"
)

var user = pb.User{Username:"remi", }

func handleMessageReceived(message pb.ChatMessage) {
	switch message.Type {
	case pb.ChatMessage_USER_JOIN:
		log.Print("New user join : " + message.User.Username)
	case pb.ChatMessage_USER_LEAVE:
		log.Print("User left the chat : " + message.User.Username)
	case pb.ChatMessage_USER_CHAT:
		log.Print("new message : " + message.Message.Content)
	}
}

func readMessage(stream pb.ChatService_StreamClient) error {
	log.Print("subscribe to message")
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Print("err : ", err)
			return err
		}
		handleMessageReceived(*msg)
	}
}

func sendMessage(serviceClient pb.ChatServiceClient) {
	stream, err := serviceClient.Stream(context.Background())

	registerMessage := pb.ChatMessage{User:&user, Type:pb.ChatMessage_USER_JOIN}
	stream.Send(&registerMessage)

	go func() {
		readMessage(stream)
	}()

	if err != nil {
		println(err)
		return
	}

	index := 0
	for {
		messageContent := pb.Message{Content:"message" + strconv.Itoa(index)}
		message := pb.ChatMessage{User:&user, Type: pb.ChatMessage_USER_CHAT, Message:&messageContent}
		stream.Send(&message)
		index = index + 1
		if index == 100 {
			stream.CloseSend()
			os.Exit(0)
		}
	}
}

func main() {
	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	serviceClient := pb.NewChatServiceClient(conn)
	sendMessage(serviceClient)
}
