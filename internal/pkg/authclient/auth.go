package authclient

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

// For demonstration purposes
const SecretServiceCode = "asdWQEfsdmkfmsdlgeruitEEFW12345!fwemofgwerg"

type UserData struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type UserResponse struct {
	Success bool     `json:"success"`
	Data    UserData `json:"data"`
}

var c *fasthttp.HostClient

func Init(host string) {
	c = &fasthttp.HostClient{
		Addr: host,
	}
}

func GetUserByToken(token string) (UserData, bool) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("http://" + c.Addr + "/get_service_user_info")
	req.Header.Set(fasthttp.HeaderAuthorization, token)
	req.Header.Set("X-Service-Key", SecretServiceCode)
	req.Header.SetHost(c.Addr)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := c.Do(req, resp)
	if err != nil || resp.StatusCode() != http.StatusOK {
		return UserData{}, false
	}

	if resp.StatusCode() != http.StatusOK {
		return UserData{}, false
	}

	var userResponse UserResponse
	if err := json.Unmarshal(resp.Body(), &userResponse); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return UserData{}, false
	}

	return UserData{
		ID:   userResponse.Data.ID,
		Role: userResponse.Data.Role,
	}, true
}
