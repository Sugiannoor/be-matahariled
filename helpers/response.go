package helpers

import "Matahariled/models"

type ResponseGetAll struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
type ResponseGetSingle struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   models.User `json:"data"`
}

type ResponseMassage struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
type ResponseError struct {
	Code   int                 `json:"code"`
	Status string              `json:"status"`
	Error  map[string][]string `json:"error"`
}
