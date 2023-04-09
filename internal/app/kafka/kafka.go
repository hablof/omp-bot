package kafka

import (
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/hablof/omp-bot/internal/app/router"
	"github.com/hablof/omp-bot/internal/config"
	"github.com/hablof/omp-bot/internal/model"
)

var _ router.CommandSender = &KafkaProducer{}

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// Send implements router.CommandSender
func (kp *KafkaProducer) Send(update tgbotapi.Update) error {

	event := model.TgMsg{
		UserId:            uint64(update.SentFrom().ID),
		Username:          update.SentFrom().String(),
		CallbackQueryData: update.CallbackData(),
		TimeStamp:         time.Now(),
		// MessageText:       update.Message.Text, // nil-safety
	}

	if update.Message != nil {
		event.MessageText = update.Message.Text
	}

	b, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.ByteEncoder(b),
	}

	kp.producer.SendMessage(&msg)

	return nil
}

func NewKafkaProducer(cfg config.Kafka) (*KafkaProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	sp, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: sp,
		topic:    cfg.Topic,
	}, nil
}
