package main

import (
	"net/url"
	"strings"
	"time"

	govalidator "tericai/common/govalidator"
	fn "tericai/common/helpers"
	md "tericai/common/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/globalsign/mgo/bson"
)

//CRUD
func credentialList(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"page":  []string{"numeric_between:1,999999"},
		"limit": []string{"numeric_between:1,999999"},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {
		var results fn.PaginateResult

		//get all te credentials
		err := fn.Paginate(collections["Credential"], nil, nil, req.Body["page"], req.Body["limit"], "", &results)
		if err != nil {
			resp := fn.Response{Error: "cannot get credentials list", Exception: err, HTTPCode: 500}
			fn.ReplyBack(req, resp)
			return
		}

		//limiting the visibility of access secret
		for i, cred := range results.Docs {
			accessSecretLength := len(cred["access_secret"].(string))
			start := cred["access_secret"].(string)[:5]
			end := cred["access_secret"].(string)[accessSecretLength-3:]
			cred["access_secret_masked"] = start + "*******" + end
			results.Docs[i] = cred
		}

		resp := fn.Response{Msg: "credentials list", Data: results, HTTPCode: 200}
		fn.ReplyBack(req, resp)
	}

}

func credentialAdd(req fn.Request) {

	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"name":          []string{"required", "unique:Credential"},
		"access_key":    []string{"required", "min:10"},
		"access_secret": []string{"required", "min:10"},
		"region":        []string{"required", "min:5"},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {

		c := md.CredentialModel{
			ID:           bson.NewObjectId(),
			Name:         req.Body["name"].(string),
			AccessKey:    req.Body["access_key"].(string),
			AccessSecret: req.Body["access_secret"].(string),
			Region:       req.Body["region"].(string),
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := collections["Credential"].Insert(c)
		if err != nil {
			resp := fn.Response{Error: "cannot insert credentials", Exception: err, HTTPCode: 500}
			fn.ReplyBack(req, resp)
			return
		}

		resp := fn.Response{Msg: "credentials inserted", Data: c, HTTPCode: 201}
		fn.ReplyBack(req, resp)
	}
}

func credentialEdit(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"credential_id": []string{"required", "objectID", "exists:Credential,_id"},
		"is_active":     []string{"bool"},
		"access_key":    []string{"required", "min:10"},
		"access_secret": []string{"required", "min:10"},
		"region":        []string{"required", "min:5"},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {

		var credential md.CredentialModel

		//get credential details
		credentialID := req.Body["credential_id"].(string)
		err := collections["Credential"].FindId(bson.ObjectIdHex(credentialID)).One(&credential)
		if err != nil {
			resp := fn.Response{Error: "cannot edit credentials", Exception: err, HTTPCode: 500}
			fn.ReplyBack(req, resp)
			return
		}

		//optional payload
		if req.Body["access_key"] != nil {
			credential.AccessKey = req.Body["access_key"].(string)
		}
		if req.Body["access_secret"] != nil {
			credential.AccessSecret = req.Body["access_secret"].(string)
		}
		if req.Body["is_active"] != nil {
			credential.IsActive = req.Body["is_active"].(bool)
		}
		if req.Body["region"] != nil {
			credential.Region = req.Body["region"].(string)
		}
		if req.Body["name"] != nil {
			credential.Name = req.Body["name"].(string)
		}

		credential.UpdatedAt = time.Now()

		//updating
		err2 := collections["Credential"].Update(bson.M{"_id": credential.ID}, bson.M{"$set": credential})
		if err2 != nil {
			resp := fn.Response{Error: "cannot edit credentials", Exception: err2, HTTPCode: 500}
			fn.ReplyBack(req, resp)
			return
		}

		resp := fn.Response{Msg: "edited credentials", HTTPCode: 201}
		fn.ReplyBack(req, resp)
	}
}

func credentialDelete(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"credential_id": []string{"required", "objectID", "exists:Credential,_id"},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {

		credentialID := req.Body["credential_id"].(string)
		err := collections["Credential"].Remove(bson.M{"_id": bson.ObjectIdHex(credentialID)})

		if err != nil {
			resp := fn.Response{Msg: "cannot delete", Exception: err, HTTPCode: 500}
			fn.ReplyBack(req, resp)
		} else {
			resp := fn.Response{Msg: "deleted successfully", HTTPCode: 200}
			fn.ReplyBack(req, resp)
		}
	}
}

func credentialView(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"credential_id": []string{"required", "objectID", "exists:Credential,_id"},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {

		var credential md.CredentialModel
		credentialID := req.Body["credential_id"].(string)

		err := collections["Credential"].FindId(bson.ObjectIdHex(credentialID)).One(&credential)

		if err != nil {
			resp := fn.Response{Exception: err, Error: "not found", Data: nil, HTTPCode: 404}
			fn.ReplyBack(req, resp)
		} else {
			resp := fn.Response{Data: credential, HTTPCode: 200}
			fn.ReplyBack(req, resp)
		}
	}
}

// list of all instances in aws
func credentialListDevices(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"credential_id": []string{"required", "objectID", "exists:Credential,_id"},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {

		var credential md.CredentialModel
		credentialID := req.Body["credential_id"].(string)

		//get credential details
		err := collections["Credential"].FindId(bson.ObjectIdHex(credentialID)).One(&credential)
		if err != nil {
			resp := fn.Response{Exception: err, Error: "not found", Data: nil, HTTPCode: 404}
			fn.ReplyBack(req, resp)
		} else {

			accesskey := credential.AccessKey
			accesssecret := credential.AccessSecret

			//connection to aws
			sess, _ := session.NewSession(&aws.Config{
				Region:      aws.String(credential.Region),
				Credentials: credentials.NewStaticCredentials(accesskey, accesssecret, ""),
			})

			svc := ec2.New(sess)
			params := &ec2.DescribeInstancesInput{
				Filters: []*ec2.Filter{
					&ec2.Filter{
						Name: aws.String("instance-state-name"),
						Values: []*string{
							aws.String("running"),
							aws.String("pending"),
							aws.String("stopped"),
						},
					},
				},
			}

			allInstances := make([]fn.Payload, 0)

			//list all ec2 instances in aws
			response, _ := svc.DescribeInstances(params)
			for idx := range response.Reservations {
				for _, inst := range response.Reservations[idx].Instances {
					name := "None"
					for _, keys := range inst.Tags {
						if *keys.Key == "Name" {
							name = url.QueryEscape(*keys.Value)
						}
					}
					name = strings.Replace(name, "+", " ", -1)

					instance := fn.Payload{
						"instance_id": inst.InstanceId,
						"name":        name,
						"public_ip":   inst.PublicIpAddress,
					}

					allInstances = append(allInstances, instance)

				}
			}

			resp := fn.Response{Data: allInstances, HTTPCode: 200}
			fn.ReplyBack(req, resp)
		}

	}

}
