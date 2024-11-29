package service

import (
	"fmt"
	"time"

	"github.com/kaibling/iggy-extensions/config"
	"github.com/kaibling/iggy-extensions/models"
	"github.com/kaibling/iggy-extensions/pkg/log"
	"github.com/nats-io/nats.go"
)

func NewNATSClient(cfg config.Config, l *log.Logger) (*NATSClient, error) {
	nc, err := nats.Connect(cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS server: %w", err)
	}

	return &NATSClient{cfg: cfg, l: l.NewScope("Subscriber"), nc: nc}, nil
}

type NATSClient struct {
	cfg config.Config
	l   *log.Logger
	nc  *nats.Conn
}

func (n *NATSClient) Publish(channelName string, message []byte) error {
	err := n.nc.Publish(channelName, message)
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}

	n.l.Debug(fmt.Sprintf("Message published to %s message: %s", channelName, message))

	return nil
}

func (n *NATSClient) Subscribe(channelName string, db *DiscordBot) error {
	sub, err := n.nc.SubscribeSync(channelName)
	if err != nil {
		return fmt.Errorf("error subscribing to subject: %v", err)
	}

	n.l.Debug(fmt.Sprintf("Subscribed to subject: %s", channelName))

	//Process messages
	for {
		msg, err := sub.NextMsg(1 * time.Minute)
		if err != nil {
			if err == nats.ErrTimeout {
				n.l.Debug("listening timeout: no messages received")
			} else {
				n.l.Info(fmt.Sprintf("mesasge read error occurred: %v", err))
			}
			continue
		}
		n.l.Debug("Received message")
		writeMsg, err := models.ToWriteMessage(msg.Data)
		if err != nil {
			n.l.Error(err)
			continue
		}

		l := n.l.NewScope("subscriber")
		l.AddStringField("request_id", writeMsg.RequestID)
		l.Info(writeMsg.Message)

		// send message
		if err := db.WriteToChannel(writeMsg.ChannelID, writeMsg.Message); err != nil {
			l.Error(err)
		}

	}
}
