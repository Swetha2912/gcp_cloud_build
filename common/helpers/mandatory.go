package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"sync"
	"tericai/common/govalidator"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/streadway/amqp"
)

// UserTeam -- user's team
type UserTeam struct {
	Name   string        `json:"name"`
	Role   string        `json:"role"`
	TeamID bson.ObjectId `json:"team_id"`
}

// UserBlock -- user info
type UserBlock struct {
	ID    bson.ObjectId `json:"_id"`
	Name  string        `json:"name"`
	Email string        `json:"email"`
	Role  string        `json:"role"`
	Phone string        `json:"phone"`
}

// RequestedUser object
type RequestedUser struct {
	Info    UserBlock           `json:"user"`
	Modules map[string][]string `json:"modules"`
	Teams   []UserTeam          `json:"teams"`
	// Role   string
	// Phone  string
}

// Request object
type Request struct {
	Body          Payload
	User          RequestedUser
	CorrelationID string
	ReplyTo       string
	RoutingKey    string
	Module        string // module name ex : user
	Action        string // Function ex : login
}

// RPC object -- used inside fn.Execute
type RPC struct {
	InChannel  chan bool
	Body       map[string]interface{}
	BodyKey    string
	Pattern    string
	Name       string
	OutChannel chan bool
	Timeout    time.Duration
}

// MongoCall object -- used inside fn.Execute
type MongoCall struct {
	InChannel  chan bool
	Collection *mgo.Collection
	Select     bson.M
	Filter     bson.M
	Name       string
	OutChannel chan bool
}

// CallBack object -- used inside fn.Execute
type CallBack struct {
	InChannel  chan bool
	Function   interface{}
	OutChannel chan bool
}

// Response returned
type Response struct {
	Msg       string      `json:"msg,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	HTTPCode  int         `json:"http_code,omitempty"`
	Error     string      `json:"error,omitempty"`
	Exception error       `json:"ex,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
}

// PaginateResult with docs structure
type PaginateResult struct {
	TotalCount  int                      `json:"total"`
	Limit       int                      `json:"limit"`
	CurrentPage int                      `json:"current_page"`
	Docs        []map[string]interface{} `json:"docs"`
}

//Paginate function for page limit
func Paginate(collection *mgo.Collection, query bson.M, selectFilter bson.M, _page interface{}, _limit interface{}, sort string, pRes *PaginateResult) error {

	var page int
	var limit int

	var results []map[string]interface{}

	if _page == nil {
		page = 1
	} else {
		switch _page.(type) {
		case int:
			page = _page.(int)
			break
		case float64:
			page = int(_page.(float64))
			break
		case string:
			page, _ = strconv.Atoi(_page.(string))
			break
		}
	}

	if _limit == nil {
		limit = 10
	} else {

		switch _limit.(type) {
		case int:
			limit = _limit.(int)
			break
		case float64:
			limit = int(_limit.(float64))
			break
		case string:
			limit, _ = strconv.Atoi(_limit.(string))
			break
		}
	}

	count, _ := collection.Find(query).Count()
	pRes.TotalCount = count
	pRes.Limit = limit
	pRes.CurrentPage = page

	mongoQuery := collection.Find(query).Select(selectFilter).Skip(limit * (page - 1)).Limit(limit)
	if sort != "" {
		mongoQuery = mongoQuery.Sort(sort)
	}

	err := mongoQuery.All(&results)
	pRes.Docs = results

	if results == nil {
		pRes.Docs = make([]map[string]interface{}, 0)
	}

	return err
}

//Sorting function for images sorting
func Sorting(collection *mgo.Collection, query bson.M, sort string, _page interface{}, _limit interface{}, pRes *PaginateResult) error {

	var page int
	var limit int

	var results []map[string]interface{}

	if _page == nil {
		page = 1
	} else {
		switch _page.(type) {
		case int:
			page = _page.(int)
			break
		case float64:
			page = int(_page.(float64))
			break
		case string:
			page, _ = strconv.Atoi(_page.(string))
			break
		}
	}

	if _limit == nil {
		limit = 10
	} else {

		switch _limit.(type) {
		case int:
			limit = _limit.(int)
			break
		case float64:
			limit = int(_limit.(float64))
			break
		case string:
			limit, _ = strconv.Atoi(_limit.(string))
			break
		}
	}

	count, _ := collection.Find(query).Count()
	pRes.TotalCount = count
	pRes.Limit = limit
	pRes.CurrentPage = page
	err := collection.Find(query).Sort(sort).Skip(limit * (page - 1)).Limit(limit).All(&results)
	pRes.Docs = results

	if results == nil {
		pRes.Docs = make([]map[string]interface{}, 0)
	}

	return err
}

// ValidateRespond will validate payload and respond errors directly
func ValidateRespond(req Request, rules govalidator.MapData, Collections map[string]*mgo.Collection) bool {

	messages := govalidator.MapData{}

	opts := govalidator.Options{
		Body:            req.Body,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
		Collections:     Collections,
	}

	v := govalidator.New(opts)
	e := v.Validate()

	validationErrs := map[string]interface{}{"Errors": e}

	if validationErrs != nil && len(e) > 0 {
		resp := Response{HTTPCode: 400, Msg: "validation errors", Errors: validationErrs["Errors"]}
		ReplyBack(req, resp)
		return false
	}
	return true
}

// Validate will validate payload and return errors
func Validate(body map[string]interface{}, rules govalidator.MapData, Collections map[string]*mgo.Collection) (bool, interface{}) {

	messages := govalidator.MapData{}

	opts := govalidator.Options{
		Body:            body,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
		Collections:     Collections,
	}

	v := govalidator.New(opts)
	e := v.Validate()

	validationErrs := map[string]interface{}{"Errors": e}

	if validationErrs != nil && len(e) > 0 {
		return false, validationErrs["Errors"]
	}
	return true, nil
}

// CleanPayload will clean post payload to make compatible with validate function
func CleanPayload(reqBody map[string]interface{}) map[string]interface{} {

	cleanedPayload := make(map[string]interface{})

	for key, value := range reqBody {
		if value != nil && value != "" {
			cleanedPayload[key] = value
		}
	}

	return cleanedPayload
}

// ReplyBack will push response into RabbitMQ adding correlation_id and reply_to
func ReplyBack(req Request, res Response) {

	// ReplyTo := req["ReplyTo"]
	CorrelationID := req.CorrelationID
	ReplyTo := req.ReplyTo

	if res.Exception != nil {
		fmt.Println("Log -- exception", GetString(res.Exception))
	}

	jsonBytesArray, err1 := json.Marshal(res)
	FailOnError(err1, "Cannot decode JSON")

	Log("Responding to "+ReplyTo+" CorrelationID : "+CorrelationID, 0)
	fmt.Println(" Log -- respond payload", res)

	// In Golang there is no sendToQueue() -- Instead use empty exchange
	err2 := Ch.Publish(
		"",      // Exchange global
		ReplyTo, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			CorrelationId: CorrelationID,
			Body:          jsonBytesArray,
		})

	FailOnError(err2, "Failed to publish a message")
}

// RaiseEvent will broadcast msg to everyone
func RaiseEvent(pattern string, reqBody interface{}) {

	// err1 := Ch.ExchangeDeclare(
	// 	"events", // name
	// 	"fanout", // type
	// 	false,    // durable
	// 	false,    // auto-deleted
	// 	false,    // internal
	// 	false,    // no-wait
	// 	nil,      // arguments
	// )

	jsonBytesArray, err1 := json.Marshal(reqBody)
	FailOnError(err1, "Cannot decode JSON")

	Mutex.Lock()
	// In Golang there is no sendToQueue() -- Instead use empty exchange
	err2 := Ch.Publish(
		"events", // Exchange global
		pattern,  // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Body: jsonBytesArray,
		})
	Mutex.Unlock()

	FailOnError(err2, "Failed to raise event")
	fmt.Println("Log -- raised Event", pattern)
}

// CheckPermission - will run the endpoint - recieve response and continues further execution
func CheckPermission(resourceName string, resourceType string, resourceID interface{}, permission string, req Request, status chan<- bool) {

	uID, _ := GetUUID()
	CorrelationID := GetString(uID)

	// frame a reqBody to send to checkPermission
	body := make(Payload)
	body["resource_type"] = resourceType
	body["user_id"] = req.User.Info.ID
	body["account_type"] = req.User.Info.Role
	if resourceID != nil {
		body["resource_id"] = bson.ObjectIdHex(resourceID.(string))
	}
	body["resource_type"] = resourceType
	body["resource_name"] = resourceName
	body["permission"] = permission

	fmt.Println("Log -- checkpermissionbody", body)

	jsonBytesArray, _ := json.Marshal(body)

	isReplied := make(chan Response)

	Mutex.Lock()
	Replies[CorrelationID] = isReplied
	Mutex.Unlock()

	fmt.Println("Log -- rabbitmq_name", os.Getenv("RABBIT_MQ_QUEUE_NAME"))

	// // In Golang there is no sendToQueue() -- Instead use empty exchange
	err2 := Ch.Publish(
		os.Getenv("RABBIT_MQ_EXCHANGE_NAME"), // Exchange global
		"acl.checkPermission",                // routing key
		false,                                // mandatory
		false,                                // immediate
		amqp.Publishing{
			ReplyTo:       os.Getenv("RABBIT_MQ_QUEUE_NAME"),
			CorrelationId: CorrelationID,
			Body:          jsonBytesArray,
		})

	if err2 != nil {
		fmt.Println("Log -- cannot publish", err2)
	}

	select {
	case reply := <-isReplied:
		fmt.Println("Log -- got ACL reply", reply)
		if reply.HTTPCode != 200 {
			resp := Response{Msg: "access denied", HTTPCode: 403, Error: reply.Error, Errors: reply.Errors}
			ReplyBack(req, resp)
			status <- false
		} else {
			status <- true
		}
	case <-time.After(3 * time.Second):
		resp := Response{Msg: "ACL timeout", HTTPCode: 500}
		ReplyBack(req, resp)
		status <- false
	}
}

// Execute - will run the endpoint - recieve response and continues further execution
func Execute(pattern string, res interface{}, timeout time.Duration) (bool, Response) {

	if timeout <= 5 {
		timeout = 5
	}

	Log("* calling "+pattern+" send reply to "+os.Getenv("RABBIT_MQ_QUEUE_NAME"), 0)

	// ReplyTo := req["ReplyTo"]
	uID, _ := GetUUID()
	// ReplyTo := "dataset"
	CorrelationID := fmt.Sprint(uID)

	jsonBytesArray, err1 := json.Marshal(res)
	FailOnError(err1, "Cannot decode JSON")

	replyBackFlag := make(chan Response)

	Mutex.Lock()
	Replies[CorrelationID] = replyBackFlag
	Mutex.Unlock()

	// // In Golang there is no sendToQueue() -- Instead use empty exchange
	err2 := Ch.Publish(
		os.Getenv("RABBIT_MQ_EXCHANGE_NAME"), // Exchange global
		pattern,                              // routing key
		false,                                // mandatory
		false,                                // immediate
		amqp.Publishing{
			ReplyTo:       os.Getenv("RABBIT_MQ_QUEUE_NAME"),
			CorrelationId: CorrelationID,
			Body:          jsonBytesArray,
		})

	if err2 != nil {
		return false, Response{Error: "cannot publish message", Exception: err2, HTTPCode: 500}
	}

	select {
	case reply := <-replyBackFlag:
		if reply.HTTPCode != 200 {
			return false, reply
		}

		return true, reply

	case <-time.After(timeout * time.Second):
		return false, Response{Error: "timeout", HTTPCode: 500}
	}

}

// BulkExecute - will run multiple RPC calls and bunch all response
func BulkExecute(requests []interface{}) map[string]Response {
	var wg sync.WaitGroup
	var responses = make(map[string]Response, 0)

	wg.Add(len(requests))

	for _, _req := range requests {

		fmt.Println("Log -- ", reflect.TypeOf(_req))

		switch _req.(type) {

		case RPC:
			req := _req.(RPC)
			go func(req RPC) {

				// if no timeout is given use default timeout
				if req.Timeout == 0 {
					req.Timeout = 2
				}

				if req.InChannel != nil {
					previousTaskStatus := <-req.InChannel
					if !previousTaskStatus {
						Mutex.Lock()
						responses[req.Name] = Response{Error: "skipped", HTTPCode: 500}
						Mutex.Unlock()
						fmt.Println("Log -- skipping task since previous did not execute successfully")
						if req.OutChannel != nil {
							req.OutChannel <- false
						}
						wg.Done()
						return
					}
				}

				uID, _ := GetUUID()
				CorrelationID := GetString(uID)

				// if bodykey is present then inject those params into the body
				Mutex.RLock()
				if len(req.BodyKey) > 0 && responses[req.BodyKey].Data != nil {
					req.Body[req.BodyKey] = responses[req.BodyKey].Data.(Payload)
				}
				Mutex.RUnlock()

				jsonBytesArray, err1 := json.Marshal(req.Body)
				if err1 != nil {
					Mutex.Lock()
					responses[req.Name] = Response{
						Exception: err1,
						HTTPCode:  500,
						Error:     "Cannot convert to JSON",
					}
					Mutex.Unlock()
					if req.OutChannel != nil {
						req.OutChannel <- false
					}
					wg.Done()
					return
				}

				isReplied := make(chan Response)

				Mutex.Lock()
				Replies[CorrelationID] = isReplied
				Mutex.Unlock()

				// // In Golang there is no sendToQueue() -- Instead use empty exchange
				Ch.Publish(
					os.Getenv("RABBIT_MQ_EXCHANGE_NAME"), // Exchange global
					req.Pattern,                          // routing key
					false,                                // mandatory
					false,                                // immediate
					amqp.Publishing{
						ReplyTo:       os.Getenv("RABBIT_MQ_QUEUE_NAME"),
						CorrelationId: CorrelationID,
						Body:          jsonBytesArray,
					})

				if req.Timeout <= 5 {
					req.Timeout = 5
				}

				select {
				case resp := <-isReplied:
					if resp.HTTPCode != 200 {
						if req.OutChannel != nil {
							req.OutChannel <- false
						}
					} else {
						if req.OutChannel != nil {
							req.OutChannel <- true
						}
					}
					Mutex.Lock()
					responses[req.Name] = resp
					Mutex.Unlock()
					wg.Done()
				case <-time.After(req.Timeout * time.Second):
					Mutex.Lock()
					responses[req.Name] = Response{Error: "timeout", HTTPCode: 500}
					Mutex.Unlock()
					wg.Done()
					if req.OutChannel != nil {
						req.OutChannel <- false
					}
				}
			}(req)

		case MongoCall:
			req := _req.(MongoCall)

			go func(req MongoCall) {

				if req.InChannel != nil {
					previousTaskStatus := <-req.InChannel
					if !previousTaskStatus {
						Mutex.Lock()
						responses[req.Name] = Response{Error: "skipped", HTTPCode: 500}
						Mutex.Unlock()
						fmt.Println("Log -- skipping task since previous did not execute successfully")
						wg.Done()
						if req.OutChannel != nil {
							req.OutChannel <- false
						}
						return
					}
				}

				var results []interface{}
				query := req.Collection.Find(req.Filter)
				if len(req.Select) > 0 {
					query = query.Select(req.Select)
				}

				err := query.All(&results)
				if err != nil {
					Mutex.Lock()
					responses[req.Name] = Response{Exception: err, Error: "cannot get results", HTTPCode: 500}
					Mutex.Unlock()
					wg.Done()
					if req.OutChannel != nil {
						req.OutChannel <- false
					}
					return
				}

				responses[req.Name] = Response{Data: results, HTTPCode: 200}
				if req.OutChannel != nil {
					req.OutChannel <- true
				}
				wg.Done()

			}(req)
		case CallBack:

			req := _req.(CallBack)

			go func(req CallBack) {
				if req.InChannel != nil {
					previousTaskStatus := <-req.InChannel
					if !previousTaskStatus {
						fmt.Println("Log -- skipping task since previous did not execute successfully")
						if req.OutChannel != nil {
							req.OutChannel <- false
						}
						wg.Done()
						return
					}
				}

				fmt.Println("Log -- passing to func() -->", req, responses)
				// call function
				req.Function.(func(CallBack, map[string]Response))(req, responses)
				wg.Done()

			}(req)
		}

	}

	fmt.Println("Log -- all triggered ", len(requests))

	wg.Wait()
	return responses
}
