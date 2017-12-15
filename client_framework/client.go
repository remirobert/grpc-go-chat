package client_framework

import (
	pb "grpc-go-chat/chat"
	"google.golang.org/grpc"
	"context"
	"log"
	"github.com/golang/protobuf/proto"
)

type ClientCallback interface {
	GetNewMessage(message []byte)
}

type Client struct {
	callback ClientCallback
	conn   *grpc.ClientConn
	serv   pb.ChatServiceClient
	stream pb.ChatService_StreamClient
	user pb.User
}

func (c *Client) JoinChat(username string) {
	if c.serv == nil || c.stream == nil {
		return
	}
	user := pb.User{
		Username:username,
		Id:"123",
	}
	c.user = user
	msg := pb.ChatMessage{
		Type: pb.ChatMessage_USER_JOIN,
		User: &user,
	}
	c.stream.Send(&msg)
}

func (c *Client) readMessage() error {
	for {
		msg, err := c.stream.Recv()
		if err != nil {
			log.Print("err : ", err)
			return err
		}
		data, err := proto.Marshal(msg)
		if err != nil {
			return err
		}
		c.callback.GetNewMessage(data)
		log.Print("get msg : ", msg)
	}
}

func (c *Client) StartStream() error {
	stream, err := c.serv.Stream(context.Background())

	if err == nil {
		return err
	}
	c.stream = stream

	go c.readMessage()
	return nil
}

func (c *Client) SendMessages(msgBytes []byte) error {
	message := &pb.Message{}
	err := proto.Unmarshal(msgBytes, message)
	if err != nil {
		return err
	}
	if c.stream == nil {
		return nil
	}
	msg := pb.ChatMessage{
		User:&c.user,
		Type:pb.ChatMessage_USER_CHAT,
		Message:message,
	}
	c.stream.Send(&msg)
	return nil
}

func (c *Client) Connect(target string) error {
	log.Print("connecting to the server : ", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	log.Print(conn)
	log.Print(err)
	if err != nil {
		return err
	}
	c.conn = conn
	c.serv = pb.NewChatServiceClient(conn)
	return c.StartStream()
}

func (c *Client) Disconnect() error {
	return c.conn.Close()
}

func NewClient(callback ClientCallback) *Client {
	return &Client{
		callback:callback,
	}
}
