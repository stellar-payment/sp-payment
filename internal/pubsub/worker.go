package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/inconst"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/service"
)

type EventPubSub struct {
	logger       zerolog.Logger
	redis        *redis.Client
	service      service.Service
	secureRoutes []string
}

type NewEventPubSubParams struct {
	Logger       zerolog.Logger
	Redis        *redis.Client
	Service      service.Service
	SecureRoutes []string
}

func NewEventPubSub(params *NewEventPubSubParams) *EventPubSub {
	return &EventPubSub{
		logger:       params.Logger,
		redis:        params.Redis,
		service:      params.Service,
		secureRoutes: params.SecureRoutes,
	}
}

func (pb *EventPubSub) Listen() {
	ctx := context.Background()

	subscriber := pb.redis.Subscribe(ctx,
		inconst.TOPIC_REQUEST_SECURE_ROUTE,
		inconst.TOPIC_CREATE_CUSTOMER,
		inconst.TOPIC_DELETE_CUSTOMER,
		inconst.TOPIC_CREATE_MERCHANT,
		inconst.TOPIC_DELETE_MERCHANT,
		inconst.TOPIC_CREATE_TRX,
	)

	data := fmt.Sprintf("%s,%s", "payment", strings.Join(pb.secureRoutes, ","))
	pb.redis.Publish(context.Background(), inconst.TOPIC_BROADCAST_SECURE_ROUTE, data)

	defer subscriber.Close()
	for msg := range subscriber.Channel() {
		switch msg.Channel {
		case inconst.TOPIC_CREATE_CUSTOMER:
			data := &indto.EventCustomer{}
			if err := json.Unmarshal([]byte(msg.Payload), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Msg("failed to marshal payload")
				continue
			}

			if err := pb.service.HandleCreateCustomer(context.Background(), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Send()
				continue
			}
		case inconst.TOPIC_DELETE_CUSTOMER:
			data := &indto.EventCustomer{}
			if err := json.Unmarshal([]byte(msg.Payload), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Msg("failed to marshal payload")
				continue
			}

			if err := pb.service.HandleDeleteCustomer(context.Background(), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Send()
				continue
			}
		case inconst.TOPIC_CREATE_MERCHANT:
			data := &indto.EventMerchant{}
			if err := json.Unmarshal([]byte(msg.Payload), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Msg("failed to marshal payload")
				continue
			}

			if err := pb.service.HandleCreateMerchant(context.Background(), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Send()
				continue
			}
		case inconst.TOPIC_DELETE_MERCHANT:
			data := &indto.EventMerchant{}
			if err := json.Unmarshal([]byte(msg.Payload), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Msg("failed to marshal payload")
				continue
			}

			if err := pb.service.HandleDeleteMerchant(context.Background(), data); err != nil {
				pb.logger.Warn().Err(err).Str("channel", msg.Channel).Send()
				continue
			}
		case inconst.TOPIC_CREATE_TRX:
		case inconst.TOPIC_REQUEST_SECURE_ROUTE:
			data := fmt.Sprintf("%s,%s", "payment", strings.Join(pb.secureRoutes, ","))
			pb.redis.Publish(context.Background(), inconst.TOPIC_BROADCAST_SECURE_ROUTE, data)
		}
	}
}
