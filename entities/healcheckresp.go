package entities

import (
	"encoding/json"
	"log"
)

/*HealthCheckResponse response for a health check*/
type HealthCheckResponse struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

/*NewHealthCheckResponseFromBytes returns a response from a byte string*/
func NewHealthCheckResponseFromBytes(bytes []byte) *HealthCheckResponse {
	var resp *HealthCheckResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		log.Println("could not unmarshall health check json")
	}
	return resp
}

/*IsHealthy returns whether or not a instance is healthy*/
func (hcr *HealthCheckResponse) IsHealthy() bool {
	if hcr.State == "healthy" {
		return true
	}
	return false
}
