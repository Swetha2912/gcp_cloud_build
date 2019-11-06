package main

import (
	"fmt"
	helpers "tericai/common/helpers"
)

func callMethod(callbackFn string, req helpers.Request) {
	fmt.Print(callbackFn)

	if callbackFn == "credential.list" {
		credentialList(req)
	} else if callbackFn == "credential.add" {
		credentialAdd(req)
	} else if callbackFn == "credential.delete" {
		credentialDelete(req)
	} else if callbackFn == "credential.edit" {
		credentialEdit(req)
	} else if callbackFn == "credential.view" {
		credentialView(req)
	} else if callbackFn == "credential.listCredentialDevices" {
		credentialListDevices(req)
	}

}
