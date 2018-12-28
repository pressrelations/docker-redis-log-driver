package redis

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

type Config struct {
	Server         string
	Sentinels      []string
	MasterName     string
	Password       string
	Database       int
	List           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ConnectTimeout time.Duration
}

type Connection struct {
	client *redis.Client
	config *Config
}

func Connect(c *Config) (*Connection, error) {
	var client *redis.Client

	if c.Server != "" {
		client = redis.NewClient(&redis.Options{
			Addr:            c.Server,
			Password:        c.Password,
			DB:              c.Database,
			DialTimeout:     c.ConnectTimeout,
			ReadTimeout:     c.ReadTimeout,
			WriteTimeout:    c.WriteTimeout,
			PoolSize:        5,
			MinIdleConns:    1,
			MaxRetries:      10,
			MinRetryBackoff: 250 * time.Millisecond,
			MaxRetryBackoff: 1000 * time.Millisecond,
		})

		err := client.Ping().Err()

		if err != nil {
			logrus.WithError(err).Errorf("error connecting to Redis server: %q", c.Server)
			client.Close()
			return nil, err
		}

		logrus.Debugf("connection established to Redis server: %q", c.Server)
	} else {
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:      c.MasterName,
			SentinelAddrs:   c.Sentinels,
			Password:        c.Password,
			DB:              c.Database,
			DialTimeout:     c.ConnectTimeout,
			ReadTimeout:     c.ReadTimeout,
			WriteTimeout:    c.WriteTimeout,
			PoolSize:        5,
			MinIdleConns:    1,
			MaxRetries:      10,
			MinRetryBackoff: 250 * time.Millisecond,
			MaxRetryBackoff: 1000 * time.Millisecond,
		})

		err := client.Ping().Err()

		if err != nil {
			logrus.WithError(err).Errorf("error connecting to Redis server via Sentinel: %q", c.Sentinels)
			client.Close()
			return nil, err
		}
		logrus.Debugf("connection established to Redis server via Sentinel: %q", c.Sentinels)
	}

	return &Connection{
		client: client,
		config: c,
	}, nil
}

func (c *Connection) Disconnect() {
	if c.client != nil {
		c.client.Close()
		logrus.Debugf("disconnected from Redis server")
	}

	return
}

func AppendToList(str string, c *Connection) error {
	err := c.client.RPush(c.config.List, str).Err()

	if err != nil {
		return err
	}

	return nil
}
