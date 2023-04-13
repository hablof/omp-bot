package kafka

import (
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"

	"github.com/hablof/omp-bot/internal/app/router"
	"github.com/hablof/omp-bot/internal/config"
	"github.com/hablof/omp-bot/internal/model"
	"github.com/hablof/omp-bot/internal/service/logistic/mypackage"
)

var _ router.CommandSender = &KafkaProducer{}
var _ mypackage.CacheEventSender = &KafkaProducer{}

type KafkaProducer struct {
	producer        sarama.SyncProducer
	cacheEventTopic string
	tgCommandTopic  string
	maxAttempts     int
}

// SendCacheEvent implements mypackage.CacheEventSender
func (kp *KafkaProducer) SendCacheEvent(event model.CacheEvent) error {
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: kp.cacheEventTopic,
		Value: sarama.ByteEncoder(b),
	}

	_, _, err = kp.producer.SendMessage(&msg)
	if err != nil {
		return err
	}

	return nil
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
		Topic: kp.tgCommandTopic,
		Value: sarama.ByteEncoder(b),
	}

	for i := 0; i < kp.maxAttempts; i++ {
		_, _, err = kp.producer.SendMessage(&msg)
		if err == nil {
			break
		}

		log.Debug().Err(err).Msgf("KafkaProducer.Send: failed attempt #%d to send message", i+1)
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		return err
	}

	return nil
}

func NewKafkaProducer(cfg config.Kafka) (*KafkaProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	var (
		sp  sarama.SyncProducer
		err error
	)
	for i := 0; i < cfg.MaxAttempts; i++ {
		sp, err = sarama.NewSyncProducer(cfg.Brokers, config)

		if err == nil {
			break
		}

		log.Info().Err(err).Msgf("NewKafkaProducer: failed attempt %d/%d to connect to kafka", i+1, cfg.MaxAttempts)
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer:        sp,
		cacheEventTopic: cfg.CacheEventTopic,
		tgCommandTopic:  cfg.TgCommandTopic,
		maxAttempts:     cfg.MaxAttempts,
	}, nil
}
