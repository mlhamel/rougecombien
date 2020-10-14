package tasks

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/mlhamel/rougecombien/pkg/config"
)

// PubSub is handling everything needed for dealing gcloud pubsub
type PubSub struct {
	cfg       *config.Config
	projectID string
}

// NewPubSub is initiating a default pubsub client
func NewPubSub(cfg *config.Config) *PubSub {
	return &PubSub{projectID: "rougecombien"}
}

// Run is responsible for running the pubsub client
func (ps *PubSub) Run(ctx context.Context) error {
	client, err := pubsub.NewClient(ctx, ps.projectID)

	if err != nil {
		return err
	}

	topic, err := client.CreateTopic(ctx, "topic-name")

	if err != nil {
		return err
	}

	subscription := pubsub.SubscriptionConfig{Topic: topic}

	sub, err := client.CreateSubscription(ctx, "sub-name", subscription)

	if err != nil {
		return err
	}

	return sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message: %s", m.Data)
		m.Ack()
	})
}
