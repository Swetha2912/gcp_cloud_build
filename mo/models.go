package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	govalidator "tericai/common/govalidator"
	fn "tericai/common/helpers"
	md "tericai/common/models"

	"github.com/globalsign/mgo/bson"
)

// CRUD
func modelList(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"page":          []string{"numeric_between:1,999999"},
		"limit":         []string{"numeric_between:1,999999"},
		"experiment_id": []string{"objectID"},
		"problem_type":  []string{"in:detection,segmentation,classification"},
	}

	isValid, errors := fn.Validate(req.Body, rules, collections)

	if isValid == false {
		resp := fn.Response{Errors: errors, Error: "validation errors", HTTPCode: 400}
		fn.ReplyBack(req, resp)
		return
	}
	var filter = make(map[string]interface{})

	//filters for searching
	if req.Body["experiment_id"] != nil {
		filter["experiment_id"] = bson.ObjectIdHex(req.Body["experiment_id"].(string))
	}
	if req.Body["problem_type"] != nil {
		filter["problem_type"] = bson.ObjectIdHex(req.Body["problem_type"].(string))
	}

	if req.Body["search"] != nil {
		filter["title"] = bson.M{
			"$regex":   req.Body["search"].(string),
			"$options": "i",
		}
	}

	fmt.Println("Log -- filter", filter)

	var results fn.PaginateResult

	//list of models
	err := fn.Paginate(collections["Model"], filter, nil, req.Body["page"], req.Body["limit"], "", &results)

	if err != nil {
		resp := fn.Response{Error: "cannot get data", Exception: err, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}

	for i, doc := range results.Docs {
		if doc["experiment_id"] != nil {
			success, reply := fn.Execute("project.experiment.viewLite", fn.Payload{"experiment_id": doc["experiment_id"]}, 5)
			if success {
				doc["experiment"] = reply.Data
			}
			fmt.Println("Log -- reply", reply)
			results.Docs[i] = doc
		}
	}

	resp := fn.Response{Msg: "models list", Data: results, HTTPCode: 200}
	fn.ReplyBack(req, resp)

}

func modelAdd(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		// "experiment_id": []string{"required", "exists:Experiment,_id"},
		"title":        []string{"required", "unique:Model"},
		"source":       []string{"required", "in:trained,uploaded"},
		"problem_type": []string{"in:detection,segmentation,classification"},
	}
	isValid, errors := fn.Validate(req.Body, rules, collections)

	if isValid == false {
		resp := fn.Response{Errors: errors, Error: "validation errors", HTTPCode: 400}
		fn.ReplyBack(req, resp)
		return
	}

	//new record for inserting in table
	m := md.TrainedModel{
		ID:        bson.NewObjectId(),
		Title:     req.Body["title"].(string),
		IsActive:  true,
		IsFav:     false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	//optional payload
	if req.Body["source"] != nil {
		m.Source = req.Body["source"].(string)
	}
	if req.Body["problem_type"] != nil {
		m.ProblemType = req.Body["problem_type"].(string)
	}
	if req.Body["experiment_id"] != nil {
		m.ExperimentID = bson.ObjectIdHex(req.Body["experiment_id"].(string))
	}
	if req.Body["architecture"] != nil {
		m.Architecture = req.Body["architecture"].(string)
	}
	if req.Body["loss"] != nil {
		m.Loss, _ = strconv.ParseFloat((req.Body["loss"].(string)), 64)
	}
	if req.Body["accuracy"] != nil {
		m.Accuracy, _ = strconv.ParseFloat((req.Body["accuracy"].(string)), 64)
	}
	if req.Body["steps"] != nil {
		m.Steps, _ = strconv.Atoi(req.Body["steps"].(string))
	}

	//inserting
	err := collections["Model"].Insert(m)
	if err != nil {
		resp := fn.Response{Error: "cannot insert data", Exception: err, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}

	resp := fn.Response{Msg: "added model", Data: m, HTTPCode: 201}
	fn.ReplyBack(req, resp)

}

func modelEdit(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"model_id":     []string{"required", "objectID", "exists:Model,_id"},
		"is_active":    []string{"bool"},
		"is_fav":       []string{"bool"},
		"title":        []string{"unique:Model"},
		"source":       []string{"in:trained,uploaded"},
		"problem_type": []string{"in:detection,segmentation,classification"},
	}

	isValid, errors := fn.Validate(req.Body, rules, collections)

	if isValid == false {
		resp := fn.Response{Errors: errors, Error: "validation errors", HTTPCode: 400}
		fn.ReplyBack(req, resp)
		return
	}

	// get model from DB
	var model md.TrainedModel
	modelID := req.Body["model_id"].(string)
	err := collections["Model"].FindId(bson.ObjectIdHex(modelID)).One(&model)
	if err != nil {
		resp := fn.Response{Error: "cannot get model", Exception: err, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}

	// optional payload
	if req.Body["source"] != nil {
		model.Source = req.Body["source"].(string)
	}
	if req.Body["title"] != nil {
		model.Title = req.Body["title"].(string)
	}
	if req.Body["architecture"] != nil {
		model.Architecture = req.Body["architecture"].(string)
	}
	if req.Body["problem_type"] != nil {
		model.ProblemType = req.Body["problem_type"].(string)
	}
	if req.Body["loss"] != nil {
		model.Loss, _ = strconv.ParseFloat((req.Body["loss"].(string)), 32)
	}
	if req.Body["accuracy"] != nil {
		model.Accuracy, _ = strconv.ParseFloat((req.Body["accuracy"].(string)), 32)
	}
	if req.Body["steps"] != nil {
		model.Steps, _ = strconv.Atoi(req.Body["steps"].(string))
	}
	if req.Body["is_active"] != nil {
		model.IsActive, _ = strconv.ParseBool((req.Body["is_active"].(string)))
	}
	if req.Body["is_fav"] != nil {
		model.IsFav, _ = req.Body["is_fav"].(bool)
	}
	model.UpdatedAt = time.Now()

	// updating
	err2 := collections["Model"].Update(bson.M{"_id": model.ID}, bson.M{"$set": model})
	if err2 != nil {
		resp := fn.Response{Error: "cannot edit model", Exception: err2, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}
	resp := fn.Response{Msg: "edited model", HTTPCode: 201}
	fn.ReplyBack(req, resp)
}

func modelDelete(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"model_id": []string{"required", "objectID", "exists:Model,_id"},
	}
	isValid, errors := fn.Validate(req.Body, rules, collections)

	if isValid == false {
		resp := fn.Response{Errors: errors, Error: "validation errors", HTTPCode: 400}
		fn.ReplyBack(req, resp)
		return
	}

	// delete model
	modelID := bson.ObjectIdHex(req.Body["model_id"].(string))
	err := collections["Model"].Remove(bson.M{"_id": modelID})
	if err != nil {
		resp := fn.Response{Error: "cannot get model", Exception: err, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}

	resp := fn.Response{Msg: "model deleted", HTTPCode: 200}
	fn.ReplyBack(req, resp)

}

func modelView(req fn.Request) {
	defer fn.RecoverPanic()
	rules := govalidator.MapData{
		"model_id": []string{"required", "objectID", "exists:Model,_id"},
	}

	isValid, errors := fn.Validate(req.Body, rules, collections)

	if isValid == false {
		resp := fn.Response{Errors: errors, Error: "validation errors", HTTPCode: 400}
		fn.ReplyBack(req, resp)
		return
	}

	modelID := req.Body["model_id"].(string)
	var model md.TrainedModel

	//get model
	err := collections["Model"].FindId(bson.ObjectIdHex(modelID)).One(&model)

	if err != nil {
		resp := fn.Response{Error: "cannot get model", Exception: err, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}
	resp := fn.Response{Data: model, HTTPCode: 200}
	fn.ReplyBack(req, resp)

}

func modelUploadFile(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"model_id": []string{"required", "objectID", "exists:Model,_id"},
		"name":     []string{"required"},
		"file":     []string{"required"},
	}

	isValid, errors := fn.Validate(req.Body, rules, collections)

	if isValid == false {
		resp := fn.Response{Errors: errors, Error: "validation errors", HTTPCode: 400}
		fn.ReplyBack(req, resp)
		return
	}

	mdlID := req.Body["model_id"].(string)
	modelID := bson.ObjectIdHex(mdlID)
	name := req.Body["name"].(string)

	//parse file details
	var file md.UploadedFile
	jsonBytes, err := json.Marshal(req.Body["file"])
	err = json.Unmarshal(jsonBytes, &file)
	if err != nil {
		resp := fn.Response{Error: "cannot parse uploaded file data", Exception: err, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}

	tmpPath := file.Path
	defaultCredentialID := os.Getenv("DEFAULT_CREDENTIAL_ID")
	defaultBucketName := os.Getenv("DEFAULT_BUCKET_NAME")

	//get model
	var model md.TrainedModel
	err = collections["Model"].FindId(modelID).One(&model)
	if err != nil {
		resp := fn.Response{Error: "cannot get model", Exception: err, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}

	fileKey := "models/" + mdlID + "/" + name
	var reqBody = make(fn.Payload)
	reqBody["key"] = fileKey
	reqBody["filePath"] = tmpPath
	reqBody["credential_id"] = bson.ObjectIdHex(defaultCredentialID)
	reqBody["bucket_name"] = defaultBucketName

	fmt.Println("i am the reqbody from model", reqBody)

	//send request to s3 to upload file in s3
	success, reply := fn.Execute("databridge.uploadLocal", reqBody, 60)

	if !success {
		reply.Msg = "cannot upload model to s3"
		fn.ReplyBack(req, reply)
		return
	}

	var field = map[string]string{
		"id":        bson.NewObjectId().Hex(),
		"fieldname": name,
		"path":      fileKey,
	}
	model.Files = append(model.Files, field)
	model.UpdatedAt = time.Now()

	//update model with file details
	err2 := collections["Model"].Update(bson.M{"_id": modelID}, bson.M{"$set": model})
	if err2 != nil {
		resp := fn.Response{Error: "cannot upload model file", Exception: err2, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}
	resp := fn.Response{Msg: "uploaded model file", Data: model, HTTPCode: 201}
	fn.ReplyBack(req, resp)
}

func modelDeleteFile(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"model_id": []string{"required", "objectID", "exists:Model,_id"},
		"file_id":  []string{"required", "objectID"},
	}

	isValid, errors := fn.Validate(req.Body, rules, collections)

	if isValid == false {
		resp := fn.Response{Errors: errors, Msg: "validation errors", HTTPCode: 400}
		fn.ReplyBack(req, resp)
		return
	}

	mdlID := req.Body["model_id"].(string)
	modelID := bson.ObjectIdHex(mdlID)

	fID := req.Body["file_id"].(string)
	fileID := bson.ObjectIdHex(fID)

	// get model
	var model md.TrainedModel
	err := collections["Model"].FindId(modelID).One(&model)
	if err != nil {
		resp := fn.Response{Msg: "cannot get model file", Error: fmt.Sprint(err), HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}

	// search for file and remove
	for i, file := range model.Files {
		if file["id"] == fileID.Hex() {

			if len(model.Files) == 1 {
				// model.Files =
				fmt.Println(model.Files)
			} else {
				model.Files = append(model.Files[:i], model.Files[i+1:]...)
			}
			break
		}
	}

	// push updated model
	err2 := collections["Model"].Update(bson.M{"_id": modelID}, bson.M{"$set": model})
	if err2 != nil {
		resp := fn.Response{Msg: "cannot delete model file", Error: fmt.Sprint(err2), HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}
	resp := fn.Response{Msg: "deleted model file", Data: model, HTTPCode: 201}
	fn.ReplyBack(req, resp)

}

//gives models count in experiment list
func modelExperimentStats(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"exp_ids": []string{"required"},
	}

	isValid, errors := fn.Validate(req.Body, rules, collections)

	if isValid == false {
		resp := fn.Response{Errors: errors, Msg: "validation errors", HTTPCode: 400}
		fn.ReplyBack(req, resp)
		return
	}

	var experimentIDs []bson.ObjectId

	//parse all experiment ids
	bytes, _ := json.Marshal(req.Body["exp_ids"])
	err1 := json.Unmarshal(bytes, &experimentIDs)
	if err1 != nil {
		resp := fn.Response{Msg: "cannot parse experimentIDs", HTTPCode: 400}
		fn.ReplyBack(req, resp)
		return
	}

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{"experiment_id": bson.M{"$in": experimentIDs}},
		},
		bson.M{
			"$group": bson.M{
				"_id":    "$experiment_id",
				"models": bson.M{"$sum": 1},
			},
		},
	}

	// get models count
	var allExpModels []map[string]interface{}
	var mappedModelCounts = make(map[string]int)
	collections["Model"].Pipe(pipeline).All(&allExpModels)

	for _, mdl := range allExpModels {
		expID := mdl["_id"].(bson.ObjectId).Hex()
		count := mdl["models"].(int)
		mappedModelCounts[expID] = count
	}

	resp := fn.Response{Data: mappedModelCounts, HTTPCode: 200}
	fn.ReplyBack(req, resp)
}

// called by training service
func modelRawAdd(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"_id": []string{"required"},
	}
	isValid, errors := fn.Validate(req.Body, rules, collections)

	if isValid == false {
		resp := fn.Response{Errors: errors, Error: "validation errors", HTTPCode: 400}
		fn.ReplyBack(req, resp)
		return
	}

	var trainedModel md.TrainedModel

	bytesArray, err := json.Marshal(req.Body)
	err = json.Unmarshal(bytesArray, &trainedModel)
	if err != nil {
		resp := fn.Response{Error: "cannot parse input as json", Exception: err, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}

	//insert tained model record
	err2 := collections["Model"].Insert(trainedModel)
	if err2 != nil {
		resp := fn.Response{Error: "cannot insert data", Exception: err2, HTTPCode: 500}
		fn.ReplyBack(req, resp)
		return
	}

	resp := fn.Response{Msg: "added model", Data: trainedModel, HTTPCode: 200}
	fn.ReplyBack(req, resp)

}
