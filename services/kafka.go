package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/ponyjackal/eventual-consistency-blog/infra/database"
	"github.com/ponyjackal/eventual-consistency-blog/models"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type NewPostMessage struct {
	ID			uuid.UUID		`json:"id"`
	Title		string		`json:"title"`
	Content		string		`json:"content"`
}

type PublishedPostMessage struct {
	models.Post
}

type Publisher struct {
	newPostReader 			*kafka.Reader
	publishedPostWriter		*kafka.Writer
}

func NewPublisher() (*Publisher, func()) {
	kafkaBrokerURL := os.Getenv("KAFKA_BROKER")
	kafkaUsername := os.Getenv("KAFKA_USERNAME")
	kafkaPassword := os.Getenv("KAFKA_PASSWORD")

	p := &Publisher{}
	mechanism, err := scram.Mechanism(scram.SHA256, kafkaUsername, kafkaPassword)
	if err != nil {
		log.Fatalln(err)
	}

	// setup kafka
	dialer := &kafka.Dialer{SASLMechanism: mechanism, TLS: &tls.Config{}}
	p.newPostReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokerURL},
		Topic: "app.newPosts",
		GroupID: "service.publisher",
		Dialer: dialer,
	})
	p.publishedPostWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBrokerURL},
		Topic: "app.publishedPosts",
		Dialer: dialer,
	})

	return p, func() {
		p.newPostReader.Close()
		p.publishedPostWriter.Close()
	}
}

func (p *Publisher) Run() {
	for {
		newPost, err := p.newPostReader.FetchMessage(context.Background())
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Fatalln(err)
		}

		var post NewPostMessage
		if err := json.Unmarshal(newPost.Value, &post); err != nil {
			log.Printf("decoding new post error: %s\n", err.Error())
		}

		postModel := models.Post{
			ID: post.ID,
			Title: post.Title,
			Content: post.Content,
			Slug: slug.Make(post.Title + "-" + time.Now().Format(time.Stamp)),
		}
		if err := database.DB.Create(&postModel).Error; err != nil {
			log.Printf("saving new post in database: %s\n", err.Error())
		}
		p.newPostReader.CommitMessages(context.Background(), newPost)

		b, _ := json.Marshal(PublishedPostMessage{Post: postModel})
		p.publishedPostWriter.WriteMessages(context.Background(), kafka.Message{
			Value: b,
		})
		log.Printf("the %s post has been saved in the database\n", post.ID)

	}
}