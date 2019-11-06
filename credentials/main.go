package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	// "context"
	// "reflect"

	fn "tericai/common/helpers"

	"github.com/globalsign/mgo"
	"github.com/joho/godotenv"
)

var mutex sync.RWMutex
var collections = make(map[string]*mgo.Collection)

func main() {

	// license check -- CAUTION
	sTime := time.Now()
	licenseValid := fn.VerifyLicense()
	fmt.Println("------------------------------------------------ \n licence checked in ", time.Since(sTime), "\n------------------------------------------------ ")
	if !licenseValid {
		fn.Log("invalid license", 0)
		os.Exit(3)
	}

	fn.RecoverPanic()
	startTime := time.Now()

	// loading .env file
	err := godotenv.Load()
	if err != nil {
		fn.Log("Error loading .env file", 0)
	}

	fn.SetAppMode(os.Getenv("DEBUG_MODE"))

	//MongoDB connection
	credentialSession, err := mgo.Dial("mongodb://" + os.Getenv("CREDENTIAL_DB_HOST"))
	if credentialSession != nil {
		collections["Credential"] = credentialSession.DB(os.Getenv("CREDENTIAL_DB_NAME")).C("credentials")

		credentialSession.SetSocketTimeout(1 * time.Second)
	} else {
		fn.Log("cannot connect to DB exiting ...", 0)
		os.Exit(3)
	}

	// rabbitMQ initialize
	exchange := os.Getenv("RABBIT_MQ_EXCHANGE_NAME")
	pattern := os.Getenv("RABBIT_MQ_PATTERN")
	queueName := os.Getenv("RABBIT_MQ_QUEUE_NAME")
	rabbitConn := "amqp://" + os.Getenv("RABBIT_MQ_USER") + ":" + os.Getenv("RABBIT_MQ_PASSWORD") + "@" + os.Getenv("RABBIT_MQ_HOST")

	fn.InitMsgBroker(rabbitConn)
	patternRequests, err := fn.NewQueue(exchange, pattern)
	fn.GetReplies(exchange, queueName)
	eventRequests, err2 := fn.GetEvents(exchange)

	if err != nil || err2 != nil {
		fmt.Println("Log -- cannot start RabbitMQ ", err)
		os.Exit(3)
	}

	fmt.Println("------------------------------------------------ \n Database + RabbitMQ initialized : ", time.Since(startTime), "\n ------------------------------------------------ ")

	fn.PrintAppMode()

	aclEnabled, err := strconv.ParseBool(os.Getenv("ACL_STATUS"))
	if err != nil {
		aclEnabled = true // default = ACL true
	}

	fn.Log("ACL status "+fn.GetString(aclEnabled), 1)

	// listen to messages
	forever := make(chan bool)

	// listen to requests -- patterns
	go func() {
		for req := range patternRequests {
			go callMethod(req.Module+"."+req.Action, req)
		}
	}()

	go func() {
		for req := range eventRequests {
			go callMethod("event."+req.Module+"."+req.Action, req)
		}
	}()

	<-forever
}
