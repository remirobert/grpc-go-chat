package client

import (
	pb "grpc-go-chat/chat"
	"log"
	"context"
	"strconv"
)

var user = pb.User{Username:"remi", Id:"1234"}

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

	if err != nil {
		log.Print("get error = ", err)
		return
	}

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
			registerMessage := pb.ChatMessage{User:&user, Type:pb.ChatMessage_USER_LEAVE}
			stream.Send(&registerMessage)
			return
		}
	}
}

/*func main() {
	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	serviceClient := pb.NewChatServiceClient(conn)
	sendMessage(serviceClient)
}*/
