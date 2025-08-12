package main

import (
	
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/IBM/sarama"
)

type Comment struct {
	Text      string    `form:"text" json:"text"`
	
}

func main() {
	app := fiber.New()
	api := app.Group("api/v1")
	api.Post("/createComment",createComment)
	app.Listen(":3000")
}


func ConnectProducer(brokerUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokerUrl, config)
	if err != nil {
		println("Error connecting to Kafka:", err.Error())
		return nil, err
	}
	println("Connected to Kafka successfully")
	return conn, nil
	
}

func pushMessageToQueue(topic string, message []byte) {
	
	brokerUrl := []string{"localhost:9092"} 
	producer, err := ConnectProducer(brokerUrl)
	if err != nil {
		println("Failed to connect to producer:", err.Error())
		return	
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)

	if err != nil {
		println("Failed to send message:", err.Error())	
		return 
	}

	println("Message sent successfully to topic:", topic, "Partition:", partition, "Offset:", offset)
}

func createComment(c *fiber.Ctx) error {
	var comment Comment
	if err := c.BodyParser(&comment); err != nil {
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}
	
	cmtInBytes, err := json.Marshal(comment)

	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": true,
			"message": "Internal server error",
		})	
	}

	pushMessageToQueue("comments", cmtInBytes)

	c.Status(200).JSON(&fiber.Map{
		"success": true,
		"message": "Comment created successfully",
	})

	return nil
}