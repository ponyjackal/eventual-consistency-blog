package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type CacheManager struct {
	publishedPostReader		*kafka.Reader
	rdb						*redis.Client
}

func NewCacheManager() (*CacheManager, func()) {
	redisURL := os.Getenv("REDIS_URL")
	kafkaBrokerURL := os.Getenv("KAFKA_BROKER")
	kafkaUsername := os.Getenv("KAFKA_USERNAME")
	kafkaPassword := os.Getenv("KAFKA_PASSWORD")

	cm := &CacheManager{}
	mechanism, err := scram.Mechanism(scram.SHA256, kafkaUsername, kafkaPassword)
	if err != nil {
		log.Fatalln(err)
	}

	// setup redis
	opt, _ := redis.ParseURL(redisURL)
	cm.rdb = redis.NewClient(opt)

	// setup kafka
	dialer := &kafka.Dialer{SASLMechanism: mechanism, TLS: &tls.Config{}}
	cm.publishedPostReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokerURL},
		Topic:   "app.publishedPosts",
		GroupID: "service.cacheManager",
		Dialer:  dialer,
	})

	return cm, func() {
		cm.publishedPostReader.Close()
		cm.rdb.Close()
	}
}

func (c *CacheManager) Run() {
	for {
		publishedPost, err := c.publishedPostReader.FetchMessage(context.Background())
		if err != nil {
			if errors.Is(err, io.EOF) {
			   return
			}
			log.Fatalln(err)
		}

		var published PublishedPostMessage
		if err := json.Unmarshal(publishedPost.Value, &published); err != nil {
			log.Printf("decoding published post error: %s\n", err.Error())
         	continue
		}

		b, _ := json.Marshal(published.Post)
		c.rdb.Set(context.Background(), "post:"+published.Slug, b, 0)
		c.publishedPostReader.CommitMessages(context.Background(), publishedPost)
		log.Printf("the %s post has been saved in Redis\n", published.ID)
	}
}