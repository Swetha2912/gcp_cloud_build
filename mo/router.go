package main

import (
	"fmt"
	helpers "tericai/common/helpers"
)

func callMethod(callbackFn string, req helpers.Request) {
	fmt.Print(callbackFn)

	if callbackFn == "model.list" {
		modelList(req)
	} else if callbackFn == "model.add" {
		modelAdd(req)
	} else if callbackFn == "model.addRaw" {
		modelRawAdd(req)
	} else if callbackFn == "model.delete" {
		modelDelete(req)
	} else if callbackFn == "model.edit" {
		modelEdit(req)
	} else if callbackFn == "model.view" {
		modelView(req)
	} else if callbackFn == "model.experimentStats" {
		modelExperimentStats(req)
	} else if callbackFn == "model.uploadModelFile" {
		modelUploadFile(req)
	} else if callbackFn == "model.deleteModelFile" {
		modelDeleteFile(req)
	} else if callbackFn == "event.model.syncStatus" {
		eventModelSyncStatus(req)
	}
}
