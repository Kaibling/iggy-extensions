package service

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kaibling/iggy-extensions/models"
	"github.com/kaibling/iggy-extensions/pkg/log"
)

type DiscordBot struct {
	// ctx          context.Context
	session      *discordgo.Session
	log          *log.Logger
	brokerPrefix string
	broker       *NATSClient
}

func NewDiscordClient(brokerPrefix, token string, log *log.Logger, broker *NATSClient) (*DiscordBot, error) {
	session, err := discordgo.New("Bot " + token)
	session.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) {
		log.Info("Bot is ready")
	})

	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	db := &DiscordBot{session, log, brokerPrefix, broker}
	session.AddHandler(db.messageCreate)

	if err = session.Open(); err != nil {
		return nil, fmt.Errorf("error starting server: %w", err)
	}

	return db, nil
}

func (db *DiscordBot) WriteToChannel(channelID, data string) error {
	// connect to discord channel
	_, err := db.session.ChannelMessageSend(channelID, data)
	if err != nil {
		return err
	}

	return nil
}

func (db *DiscordBot) Close() error {
	return db.session.Close()
}

func (db *DiscordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// prevent self messaging
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Privat message
	if m.GuildID == "" {
		sender := m.Author.Username
		chanID := m.ChannelID

		db.log.Debug(fmt.Sprintf("Private Message from %s in channel %s message: %s", sender, chanID, m.Content))

		// send private message to broker
		brokerTopic := db.brokerPrefix + ".private_msg"

		brokerMsg, err := models.Message{
			ChannelID: "",
			Message:   m.Content,
			Timestamp: m.Timestamp,
		}.ToBytes()
		if err != nil {
			db.log.Error(err)

			return
		}

		if err := db.broker.Publish(brokerTopic, brokerMsg); err != nil {
			db.log.Error(err)

			return
		}
	}
}
