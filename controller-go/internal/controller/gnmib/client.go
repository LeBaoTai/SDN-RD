package gnmib

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/api"
	"github.com/openconfig/gnmic/pkg/api/target"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn     *grpc.ClientConn
	gnmiC    gnmi.GNMIClient
	username string
	password string
	timeout  time.Duration
}

type TargetClient struct {
	Target *target.Target
	Client *ygnmi.Client
}

type Cfg struct {
	Address    string
	Username   string
	Port       string
	Password   string
	SkipVerify bool
	Timeout    time.Duration
}

func CreateTargetConnection(cfg Cfg, ctx context.Context) (*TargetClient, error) {
	target, err := NewTarget(cfg, ctx)
	if err != nil {
		return nil, err
	}
	client, err := ygnmi.NewClient(target.Client)
	return &TargetClient{
		Target: target,
		Client: client,
	}, nil
}

func NewTarget(cfg Cfg, ctx context.Context) (*target.Target, error) {
	target, err := api.NewTarget(
		api.Address(cfg.Address+":"+cfg.Port),
		api.SkipVerify(cfg.SkipVerify),
		api.Username(cfg.Username),
		api.Password(cfg.Password),
		api.Timeout(cfg.Timeout),
	)
	if err != nil {
		return nil, err
	}

	if err := target.CreateGNMIClient(ctx); err != nil {
		return nil, err
	}
	return target, err
}

func (t *TargetClient) Close() {
	t.Close()
}

func New(cfg Cfg) (*Client, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.SkipVerify,
		NextProtos:         []string{"h2"},
	}
	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", cfg.Address, err)
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		conn:     conn,
		gnmiC:    gnmi.NewGNMIClient(conn),
		username: cfg.Username,
		password: cfg.Password,
		timeout:  timeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) ctxWithAuth() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	md := metadata.Pairs("username", c.username, "password", c.password)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, cancel
}
