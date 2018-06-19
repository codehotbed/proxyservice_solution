package queue

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	rabbit "github.com/streadway/amqp"
)

var (
	// AMQPPrefix rabbit url prefix
	AMQPPrefix     = "amqp://"
	replyQueueName = "reply_rpc"
)

// Queue Holds basic queue properties
type Queue struct {
	host     string
	port     string
	user     string
	password string
}

// NewQueue Queue parametrized constructor
func NewQueue(host, port, user, pwd string) *Queue {
	return &Queue{
		host:     host,
		password: pwd,
		port:     port,
		user:     user,
	}
}

// CallRPC Performs Remote Procedure Call to given queue
func (queue *Queue) CallRPC(queueName string, body *[]byte) (interface{}, error) {
	log.Printf("[RABBITMQ] Setting up connection with host: %v on port: %v ..", queue.host, queue.port)
	connection, err := rabbit.Dial(AMQPPrefix + queue.user + ":" + queue.password + "@" + queue.host + ":" + queue.port)
	if err != nil {
		log.Printf("[RABBITMQ] Error while setting up connection to a server. " + err.Error())
		return nil, err
	}
	defer connection.Close()
	log.Printf("[RABBITMQ] Connection established.")

	log.Printf("[RABBITMQ] Creating a channel..")
	queueChan, err := connection.Channel()
	if err != nil {
		log.Printf("[RABBITMQ] Error while creating a channel. " + err.Error())
		return nil, err
	}
	defer queueChan.Close()
	log.Printf("[RABBITMQ] Channel created.")

	log.Printf("[RABBITMQ] Declaring queue ..")
	q, err := queueChan.QueueDeclare(replyQueueName, false, false, false, false, nil)
	if err != nil {
		log.Printf("[RABBITMQ] Error while declaring queue. " + err.Error())
		return nil, err
	}
	log.Printf("[RABBITMQ] Queue declared.")

	msgs, err := queueChan.Consume(q.Name, "new_consumer", true, false, false, false, nil)
	if err != nil {
		log.Printf("[RABBITMQ] Error while reading from queue. " + err.Error())
		return nil, err
	}

	correlationID := strconv.FormatInt(time.Now().UnixNano(), 10)

	publishing := rabbit.Publishing{
		ContentType:   "application/json",
		ReplyTo:       replyQueueName,
		CorrelationId: correlationID,
	}

	if body != nil {
		publishing.Body = *body
	}

	var waitGroup sync.WaitGroup

	log.Printf("[RABBITMQ] Publishing on queue: %v ..", queueName)
	err = queueChan.Publish("", queueName, false, false, publishing)
	if err != nil {
		return nil, nil
	}
	log.Printf("[RABBITMQ] Message has been succesfuly published.")

	var response interface{}
	var receivedMessage []byte

	waitGroup.Add(1)
	go func() {
		log.Printf("[RABBITMQ] Reading from queue: [%v]", q.Name)
		for msg := range msgs {
			if msg.CorrelationId == correlationID {
				receivedMessage = msg.Body
				log.Printf("Succesfully received message: %v from queue: %v", receivedMessage, queueName)
				waitGroup.Done()
				return
			}
		}
	}()
	waitGroup.Wait()

	err = json.Unmarshal(receivedMessage, &response)
	if err != nil {
		log.Printf("[RABBITMQ] Error while unmarshalling message queue response" + err.Error())
		return nil, err
	}

	log.Printf("Succesfully unmarshalled queue: %v response: %v", queueName, response)
	return response, nil
}
