package entities

import (
	"encoding/json"
	"log"
)

/*Response a entity to display a message and a status*/
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

/*NewResponse returns a instanceof new response*/
func NewResponse(s int, m string) *Response {
	return &Response{
		Status:  s,
		Message: m,
	}
}

/*ToJSON converts the response to JSON*/
func (r *Response) ToJSON() []byte {
	bytes, err := json.Marshal(r)
	if err != nil {
		log.Println("couldnt marshall response json")
	}
	return bytes
}
