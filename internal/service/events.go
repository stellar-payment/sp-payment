package service

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
)

func (s *service) publishEvent(ctx context.Context, channel string, payload interface{}) (err error) {
	logger := zerolog.Ctx(ctx)

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error().Err(err).Str("channel", channel).Msg("failed to marshal payload")
		return
	}

	if err = s.redis.Publish(ctx, channel, string(data)).Err(); err != nil {
		logger.Error().Err(err).Str("channel", channel).Msg("failed to publish message")
		return
	}

	return
}
