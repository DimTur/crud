package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

var c *fasthttp.HostClient

type TotalResponse struct {
	Total int `json:"total"`
}

func GetMenusCount() (TotalResponse, bool) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("http://" + c.Addr + "/get_all")
	req.Header.SetHost(c.Addr)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := c.Do(req, resp)
	if err != nil || resp.StatusCode() != http.StatusOK {
		return TotalResponse{}, false
	}

	var totalResponse TotalResponse
	if err := json.Unmarshal(resp.Body(), &totalResponse); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return TotalResponse{}, false
	}

	return TotalResponse{
		Total: totalResponse.Total,
	}, true
}
