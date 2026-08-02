package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/hibiken/asynq"
)

type Client struct {
	client *asynq.Client
	cfg    *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
		client: asynq.NewClient(
			asynq.RedisClientOpt{
				Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
				Password: cfg.Redis.Pass,
				DB:       cfg.Redis.DB,
			},
		),
	}
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) EnqueueEmail(job EmailJob) error {
	return c.enqueue(
		TaskEmail,
		QueueEmail,
		job,
	)
}

func (c *Client) EnqueueWhatsApp(job WhatsAppJob) error {
	return c.enqueue(
		TaskWhatsApp,
		QueueWhatsApp,
		job,
	)
}

func (c *Client) enqueue(
	taskName string,
	queue string,
	payload any,
) error {

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", taskName, err)
	}

	task := asynq.NewTask(taskName, body)

	_, err = c.client.Enqueue(
		task,
		asynq.Queue(queue),
		asynq.MaxRetry(c.cfg.Worker.MaxRetry),
		asynq.Timeout(
			time.Duration(c.cfg.Worker.TimeoutSecs)*time.Second,
		),
		asynq.Unique(5*time.Minute),
	)

	if err != nil {
		return fmt.Errorf("enqueue %s: %w", taskName, err)
	}

	return nil
}
