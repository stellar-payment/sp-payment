package pubsub

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/inconst"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/service"
)

type EventPubSub struct {
	logger  zerolog.Logger
	redis   *redis.Client
	service service.Service
}

type NewEventPubSubParams struct {
	Logger  zerolog.Logger
	Redis   *redis.Client
	Service service.Service
}

func NewEventPubSub(params *NewEventPubSubParams) *EventPubSub {
	return &EventPubSub{
		logger:  params.Logger,
		redis:   params.Redis,
		service: params.Service,
	}
}

func (pb *EventPubSub) Listen() {
	ctx := context.Background()

	subscriber := pb.redis.Subscribe(ctx,
		inconst.TOPIC_CREATE_CUSTOMER,
		inconst.TOPIC_DELETE_CUSTOMER,
		inconst.TOPIC_CREATE_MERCHANT,
		inconst.TOPIC_DELETE_MERCHANT,
		inconst.TOPIC_CREATE_TRX,
	)

	defer subscriber.Close()
	for msg := range subscriber.Channel() {
		switch msg.Channel {
		case inconst.TOPIC_CREATE_CUSTOMER:
			data := &indto.Customer{}
			if err := json.Unmarshal([]byte(msg.Payload), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Msg("failed to marshal payload")
				continue
			}

			if err := pb.service.HandleCreateCustomer(context.Background(), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Send()
				continue
			}
		case inconst.TOPIC_DELETE_CUSTOMER:
			data := &indto.Customer{}
			if err := json.Unmarshal([]byte(msg.Payload), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Msg("failed to marshal payload")
				continue
			}

			if err := pb.service.HandleDeleteCustomer(context.Background(), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Send()
				continue
			}
		case inconst.TOPIC_CREATE_MERCHANT:
			data := &indto.Merchant{}
			if err := json.Unmarshal([]byte(msg.Payload), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Msg("failed to marshal payload")
				continue
			}

			if err := pb.service.HandleCreateMerchant(context.Background(), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Send()
				continue
			}
		case inconst.TOPIC_DELETE_MERCHANT:
			data := &indto.Merchant{}
			if err := json.Unmarshal([]byte(msg.Payload), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Msg("failed to marshal payload")
				continue
			}

			if err := pb.service.HandleDeleteMerchant(context.Background(), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Send()
				continue
			}
		case inconst.TOPIC_CREATE_TRX:
		}
	}
}
