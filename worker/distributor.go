package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendEmail(
		ctx context.Context,
		payload *PayloadSendEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	cli := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: cli,
	}
}
