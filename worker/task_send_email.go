package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendEmail = "task:send_email"

type PayloadSendEmail struct {
	Username string `json:"username"`
}

func (rtd *RedisTaskDistributor) DistributeTaskSendEmail(
	ctx context.Context,
	payload *PayloadSendEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal the payload: %w", err)
	}
	task := asynq.NewTask(TaskSendEmail, jsonPayload, opts...)
	taskInfo, err := rtd.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to  enqueue the task: %w", err)
	}

	log.Info().Str("type", task.Type()).
		Bytes("payload", taskInfo.Payload).
		Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).
		Msg("enqueued the task")

	return nil
}

func (rtp *RedisTaskProcessor) ProcessTaskSendEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal the payload: %w", asynq.SkipRetry)
	}

	user, err := rtp.store.GetUser(ctx, payload.Username)
	if err != nil {
		// Allow retry
		// if errors.Is(err, pgx.ErrNoRows) {
		// 	return fmt.Errorf("user don't exits: %w", asynq.SkipRetry)
		// }
		return fmt.Errorf("failed to get user: %w", err)
	}

	log.Info().Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed the task")

	return nil
}
