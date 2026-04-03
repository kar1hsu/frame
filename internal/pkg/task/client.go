package task

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// Client wraps asynq.Client as the task producer.
type Client struct {
	client *asynq.Client
}

func NewClient(redisAddr, password string, db int) *Client {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: password,
		DB:       db,
	})
	return &Client{client: client}
}

func (c *Client) Close() error {
	return c.client.Close()
}

// Enqueue sends a task to the default queue for immediate processing.
func (c *Client) Enqueue(typeName string, payload interface{}, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal task payload: %w", err)
	}
	task := asynq.NewTask(typeName, data)
	return c.client.Enqueue(task, opts...)
}

// EnqueueDelay sends a task to be processed after the given delay.
func (c *Client) EnqueueDelay(typeName string, payload interface{}, delay time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.ProcessIn(delay))
	return c.Enqueue(typeName, payload, opts...)
}

// EnqueueAt sends a task to be processed at the given time.
func (c *Client) EnqueueAt(typeName string, payload interface{}, at time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.ProcessAt(at))
	return c.Enqueue(typeName, payload, opts...)
}

// EnqueueUnique sends a task with deduplication within the given TTL.
func (c *Client) EnqueueUnique(typeName string, payload interface{}, uniqueTTL time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.Unique(uniqueTTL))
	return c.Enqueue(typeName, payload, opts...)
}

// EnqueueToQueue sends a task to a specific named queue.
func (c *Client) EnqueueToQueue(typeName string, payload interface{}, queue string, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.Queue(queue))
	return c.Enqueue(typeName, payload, opts...)
}
