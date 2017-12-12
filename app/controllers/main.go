package controllers

import "bzppx-codepub/app/business"

type MainController struct {
	BaseController
}

// return json success
func (this *MainController) Success() {
	message := "ok"
	data := map[string]string{
		"name": "api",
	}

	this.jsonSuccess(message, data)
}

// return json error
func (this *MainController) Error() {

	this.jsonError("ok", map[string]string{
		"name": "test",
		"pass": "test",
	})
}

// return json data
func (this *MainController) Index() {
	message := "ok"
	data := map[string]string{
		"name": "api",
	}

	// call business
	business.Main.Index()

	this.jsonSuccess(message, data)
}
