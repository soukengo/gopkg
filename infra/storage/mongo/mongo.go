package mongo

import (
	"context"
	"errors"
	"github.com/soukengo/gopkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	mgoopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	FieldId = "_id"
)

type Client struct {
	opts   *options
	config *Config
	client *mongo.Client
}

func NewClient(conf *Config) *Client {
	return NewClientWithOptions(conf)
}
func NewClientWithOptions(conf *Config, opts ...Option) *Client {
	c := new(Client)
	opt := &options{
		logger: log.GetLogger(),
	}
	c.opts = opt.apply(opts...)
	c.config = conf
	ctx := context.Background()
	clientOptions := mgoopts.Client().ApplyURI(c.config.Address)
	if len(c.config.Username) > 0 {
		credential := mgoopts.Credential{AuthSource: c.config.AuthSource, Username: c.config.Username, Password: c.config.Password}
		clientOptions = clientOptions.SetAuth(credential)
	}
	timeout := 10 * time.Second
	clientOptions.ConnectTimeout = &timeout
	m := newMonitor(c.opts)
	m.WithLog()
	m.WithTrace()
	clientOptions.Monitor = m.Entry()
	// 连接到MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	c.client = client
	return c
}

func (c *Client) Close() error {
	return c.client.Disconnect(context.Background())
}

func (c *Client) Database(storage string, opts ...*mgoopts.DatabaseOptions) *mongo.Database {
	return c.client.Database(storage, opts...)
}
func (c *Client) Collection(name string, opts ...*mgoopts.CollectionOptions) *mongo.Collection {
	return c.Database(c.config.Database).Collection(name, opts...)
}

func (c *Client) Connect(ctx context.Context) error {
	return c.client.Connect(ctx)
}
func IsErrNoDocuments(source error) bool {
	return errors.Is(source, mongo.ErrNoDocuments)
}
