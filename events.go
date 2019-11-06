package main

import (
	"fmt"
	fn "tericai/common/helpers"
	md "tericai/common/models"

	"github.com/globalsign/mgo/bson"
)

func eventModelSyncStatus(req fn.Request) {
	defer fn.RecoverPanic()

	// push it to sockets and store inside logs
	modelID := req.Body["model_id"].(string)
	status := req.Body["status"].(string)

	//get nodel
	var model md.TrainedModel
	err := collections["Model"].FindId(bson.ObjectIdHex(modelID)).One(&model)

	if err != nil {
		fmt.Println(" Log -- cannot get model")
		return
	}
	if status != "complete" {
		fmt.Println(" Log -- ", req.Body["data"])
	}

	model.Status = status

	//update model status
	err2 := collections["Model"].Update(bson.M{"_id": model.ID}, bson.M{"$set": model})
	if err2 != nil {
		resp := fn.Response{Error: "cannot edit model", Exception: err2, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}
}
