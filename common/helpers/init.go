package helpers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/streadway/amqp"
)

// Ch - RabbitMQ channel
var Ch *amqp.Channel

// Connection - RabbitMQ Connection
var Connection *amqp.Connection

// Mutex -->
var Mutex sync.RWMutex

// AppMode --> #0 = prod | #1 = uat |  #2 = test | #3 dev
var AppMode int

// AppType --> cloud or enterprise
var AppType string

// Connections -- stores DB connections (name -- *Conn)
var Connections = make(map[string]*mgo.Collection)

// Replies - ( correleationID - chan Response )
var Replies = make(map[string]chan Response)

// SetAppMode -- will set the app into dev or test or UAT mode
func SetAppMode(mode string) {
	var err error
	AppMode, err = strconv.Atoi(mode)
	if err != nil {
		AppMode = 3
	}
}

// SetAppType -- will set app type cloud or enterprise
func SetAppType(appType string) {
	if appType == "cloud" || appType == "enterprise" {
		AppType = appType
	} else {
		AppType = "enterprise"
	}
}

// PrintAppType -- prints current mode of app
func PrintAppType() {
	fmt.Println("Started Teric.ai - " + AppType + " edition")
}

// PrintAppMode -- will print app mode
func PrintAppMode() {
	if AppMode == 3 {
		fmt.Println("Log -- Running in DEV mode")
	} else if AppMode == 2 {
		fmt.Println("Log -- Running in TEST mode")
	} else if AppMode == 1 {
		fmt.Println("Log -- Running in UAT mode")
	} else if AppMode == 0 {
		fmt.Println("Log -- Running in PRODUCTION mode")
	} else {
		fmt.Println("Log -- Invalid App mode, please set with SetAppMode(int)")
	}
}

// InitMsgBroker initializes RabbitMQ
func InitMsgBroker(rabbitConn string) {

	Log("connecting to "+rabbitConn, 0)

	Connection, err := amqp.Dial(rabbitConn)
	FailOnError(err, "Cannot connect to RabbitMQ")
	if err != nil {
		// try again after some time for 10 times and finally fail after 10 - 15 sec
	}

	Ch, err = Connection.Channel()
	FailOnError(err, "Cannot get channel from RabbitMQ")
	if err != nil {
		// try again after some time for 10 times and finally fail after 10 - 15 sec
	}
}

// NewQueue - create a queue and bind accross a pattern
func NewQueue(exchange string, pattern string) (<-chan Request, error) {
	q, err := Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	err = Ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	FailOnError(err, "Failed to declare an exchange")

	err = Ch.QueueBind(
		q.Name,   // queue name
		pattern,  // routing key
		exchange, // exchange
		false,
		nil)
	FailOnError(err, "Failed to bind *.* queue")

	err = Ch.QueueBind(
		q.Name,       // queue name
		pattern+".*", // routing key
		exchange,     // exchange
		false,
		nil)
	FailOnError(err, "Failed to bind *.*.* queue")

	msgs, err := Ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	FailOnError(err, "Cannot consume messages chan")

	requests := make(chan Request)

	go func() {
		defer RecoverPanic()
		Log("listening to msgs ...", 0)
		for d := range msgs {

			Log("Log --> rawRequest "+d.RoutingKey, 0)
			var reqBody Payload
			err := ConvertMap(d.Body, &reqBody)
			req := Request{RoutingKey: d.RoutingKey, ReplyTo: d.ReplyTo, CorrelationID: d.CorrelationId, Body: reqBody}

			if err != nil {
				// kill the request with error response
				ReplyBack(req, Response{
					Error: "cannot parse given data",
					Data:  string(d.Body),
				})
				return
			}

			routingBits := strings.Split(d.RoutingKey, ".")

			if len(routingBits) == 3 {
				req.Module = routingBits[1]
				req.Action = routingBits[2]
			} else if len(routingBits) == 2 {
				req.Module = routingBits[0]
				req.Action = routingBits[1]
			} else {
				fmt.Println("Log -- ", req)
				Log("Invalid routing bits "+d.RoutingKey, 0)
				os.Exit(2)
			}

			Log("got new request : "+req.Module+" > "+req.Action, 0)

			requests <- req
		}
	}()

	Log("listening pattern : "+pattern, 0)

	return requests, err

}

// GetReplies -- get replies sent directy to queue
func GetReplies(exchange string, queueName string) error {

	fmt.Println(" Log -- listening replies on", exchange, queueName)

	q, err := Ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := Ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	FailOnError(err, "Cannot consume messages chan")

	go func() {
		defer RecoverPanic()
		Log("listening to replies ...", 0)

		for d := range msgs {

			// got reply
			Log("got reply "+d.RoutingKey, 0)
			var resp Response
			err := ConvertResponse(d.Body, &resp)

			if err != nil {
				fmt.Println(" Log -- err", string(d.Body))
				Log("Log -- cannot parse data recieved in reply", 0)
			}

			Mutex.RLock()
			fmt.Println(" Log -- Replies[d.CorrelationId]", Replies[d.CorrelationId])
			if Replies[d.CorrelationId] != nil {
				replyChannel := Replies[d.CorrelationId]
				replyChannel <- resp
				delete(Replies, d.CorrelationId)
			}
			Mutex.RUnlock()

		}
	}()

	return err

}

// GetEvents -- get replies sent directy to queue
func GetEvents(exchange string) (<-chan Request, error) {

	// get events
	err := Ch.ExchangeDeclare(
		"events", // name
		"fanout", // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	FailOnError(err, "Failed to declare an events exchange")

	// declare queue to recieve events
	eventQ, err := Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare a events queue")

	// bind queue with exhange
	err = Ch.QueueBind(
		eventQ.Name, // queue name
		"",          // routing key
		"events",    // exchange
		false,
		nil)
	FailOnError(err, "Failed to bind a events queue")

	// get messages from queue -- handler to recieve msg
	events, err := Ch.Consume(
		eventQ.Name, // queue
		"",          // consumer
		true,        // auto ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // args
	)

	requests := make(chan Request)

	go func() {
		defer RecoverPanic()
		Log("listening to events ...", 0)
		for d := range events {

			// Log("Log --> rawRequest "+d.RoutingKey, 0)
			var reqBody Payload
			err := ConvertMap(d.Body, &reqBody)
			req := Request{RoutingKey: d.RoutingKey, ReplyTo: d.ReplyTo, CorrelationID: d.CorrelationId, Body: reqBody}

			if err != nil {
				// kill the request with error response
				ReplyBack(req, Response{
					Error: "cannot parse given data",
					Data:  string(d.Body),
				})
				return
			}

			routingBits := strings.Split(d.RoutingKey, ".")

			if len(routingBits) == 3 {
				req.Module = routingBits[1]
				req.Action = routingBits[2]
			} else if len(routingBits) == 2 {
				req.Module = routingBits[0]
				req.Action = routingBits[1]
			} else {
				fmt.Println("Log -- ", req)
				Log("Invalid routing bits "+d.RoutingKey, 0)
				os.Exit(2)
			}

			// Log("got new event : "+req.Module+" > "+req.Action, 0)

			requests <- req
		}
	}()

	return requests, err

}
