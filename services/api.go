package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/ponyjackal/eventual-consistency-blog/models"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type API struct {
	rdb				*redis.Client
	newPostWriter	*kafka.Writer
}

func NewAPI() (*API, func()) {
	redisURL := os.Getenv("REDIS_URL")
	kafkaBrokerURL := os.Getenv("KAFKA_BROKER")
	kafkaUsername := os.Getenv("KAFKA_USERNAME")
	kafkaPassword := os.Getenv("KAFKA_PASSWORD")
	
	p := &API{}
	mechanism, err := scram.Mechanism(scram.SHA256, kafkaUsername, kafkaPassword)
	if err != nil {
		log.Fatalln(err)
	}

	// setup Redis
	opt, _ := redis.ParseURL(redisURL)
	p.rdb = redis.NewClient(opt)
	// setup Kafka
	dialer := &kafka.Dialer{SASLMechanism: mechanism, TLS: &tls.Config{}}
	p.newPostWriter = kafka.NewWriter(kafka.WriterConfig{
	   Brokers: []string{kafkaBrokerURL},
	   Topic:   "app.newPosts",
	   Dialer:  dialer})
 
	return p, func() {
	   p.newPostWriter.Close()
	   p.rdb.Close()
	}
}

// NewMessage returns the generated UUID and error
func (a *API) NewMessage(title, content string) (uuid.UUID, error) {
	uuid := uuid.New()
	message := NewPostMessage{
		ID: uuid,
		Title: title,
		Content: content,
	}
	
	b, _ := json.Marshal(message)
	return uuid, a.newPostWriter.WriteMessages(context.Background(), kafka.Message{
		Value: b,
	})
}

func (a *API) GetPost(slug string) (models.Post, error) {
	var p models.Post
	tr := a.rdb.Get(context.Background(), "post:"+slug)
	b, err := tr.Bytes()
	if err != nil {
	   return models.Post{}, err
	}
	json.Unmarshal(b, &p)
	return p, nil
}

