package spot

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	spotv1 "spot-sdk-go/proto-gen/spot_v1"
)

type Client struct {
	conn   *grpc.ClientConn
	client spotv1.SpotClient
}

func NewClientWithAddress(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := spotv1.NewSpotClient(conn)
	return &Client{conn: conn, client: client}, nil
}

// NewClient creates a new Client, reading the address from the environment variable if not provided.
func NewClient() (*Client, error) {
	address := os.Getenv("SPOT_GRPC_ADDRESS")
	if address == "" {
		return nil, fmt.Errorf("environment variable SPOT_GRPC_ADDRESS not set")
	}
	return NewClientWithAddress(address)
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) ListLayers() (*spotv1.Layers, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return c.client.ListLayers(ctx, &emptypb.Empty{})
}

func (c *Client) Call(input *spotv1.Input) (*spotv1.Output, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return c.client.Call(ctx, input)
}

func (c *Client) CallStreamOutput(input *spotv1.Input) ([]*spotv1.Output, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := c.client.CallStreamOutput(ctx, input)
	if err != nil {
		return nil, err
	}

	var outputs []*spotv1.Output
	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		outputs = append(outputs, resp)
	}
	return outputs, nil
}
