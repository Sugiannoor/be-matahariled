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
	Error  map[string][]string `json:"errors"`
}

type GeneralResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type DataTableResponse struct {
	CurrentPage  int           `json:"current_page"`
	FirstPageURL string        `json:"first_page_url"`
	From         int           `json:"from"`
	LastPage     int           `json:"last_page"`
	LastPageURL  string        `json:"last_page_url"`
	NextPageURL  string        `json:"next_page_url"`
	PrevPageURL  string        `json:"prev_page_url"`
	To           int           `json:"to"`
	Total        int           `json:"total"`
	Data         []interface{} `json:"data"` // Gunakan interface{} sebagai tipe data untuk data yang bervariasi
}
