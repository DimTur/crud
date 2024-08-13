package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

var pc *fasthttp.PipelineClient

type Ing struct {
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

type Recipe struct {
	ID          string `json:"id"`
	AuthorID    string `json:"user_id"`
	Name        string `json:"name"`
	Ingredients []Ing  `json:"ingredients"`
	Temperature int    `json:"temperature"`
}

type RecipeListResponse struct {
	Recipes []Recipe `json:"recipes"`
	Total   int      `json:"total"`
}

func GetAllRecipes(total int, limit int) ([]Recipe, error) {
	var allRecipes []Recipe
	pageCount := (total + limit - 1) / limit

	stopCh := make(chan struct {
		recipes []Recipe
		err     error
	}, pageCount)

	for i := 0; i < pageCount; i++ {
		go func(page int) {
			recipes, err := doReq(pc, page, limit)
			stopCh <- struct {
				recipes []Recipe
				err     error
			}{recipes, err}
		}(i)
	}

	for i := 0; i < pageCount; i++ {
		select {
		case result := <-stopCh:
			if result.err != nil {
				return nil, result.err
			}
			allRecipes = append(allRecipes, result.recipes...)
		case <-time.After(3 * time.Second):
			fmt.Println("timeout")
		}
	}

	return allRecipes, nil
}

func doReq(pc *fasthttp.PipelineClient, page int, limit int) ([]Recipe, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("http://" + pc.Addr + "/get_all?page=" + strconv.Itoa(page+1) + "&limit=" + strconv.Itoa(limit))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := pc.DoTimeout(req, resp, time.Second)
	if err != nil || resp.StatusCode() != http.StatusOK {
		return nil, errors.New("ERROR")
	}

	var response RecipeListResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, err
	}

	return response.Recipes, nil
}
