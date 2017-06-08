package redis

import (
	"time"

	"github.com/Sirupsen/logrus"
	redigo "github.com/garyburd/redigo/redis"
)

type Config struct {
	Server         string
	Password       string
	Database       int
	List           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ConnectTimeout time.Duration
}

type Connection struct {
	pool   *redigo.Pool
	config *Config
}

func Connect(c *Config) *Connection {
	pool := &redigo.Pool{
		MaxActive:   1,
		MaxIdle:     1,
		Wait:        true,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redigo.Conn, error) {
			conn, err := redigo.Dial(
				"tcp",
				c.Server,
				redigo.DialReadTimeout(c.ReadTimeout),
				redigo.DialWriteTimeout(c.WriteTimeout),
				redigo.DialConnectTimeout(c.ConnectTimeout),
				redigo.DialDatabase(c.Database),
				redigo.DialPassword(c.Password),
			)

			if err != nil {
				logrus.WithError(err).Errorf("error connecting to Redis server: %q", c.Server)
			} else {
				logrus.Debugf("connection established to Redis server: %q", c.Server)
			}

			return conn, err
		},
	}

	return &Connection{
		pool:   pool,
		config: c,
	}
}

func (c *Connection) Disconnect() {
	if c.pool != nil {
		c.pool.Close()
	}

	return
}

func AppendToList(str string, c *Connection) error {
	conn := c.pool.Get()
	_, err := conn.Do("RPUSH", c.config.List, str)
	conn.Close()

	if err != nil {
		return err
	}

	return nil
}
