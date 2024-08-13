package main

import (
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

func main() {
	c = &fasthttp.HostClient{
		Addr: "localhost:8080",
	}

	total, success := GetMenusCount()
	if !success {
		log.Println("Failed to get menus count")
		return
	}

	pc = &fasthttp.PipelineClient{
		Addr:               "localhost:8080",
		MaxConns:           1,
		MaxPendingRequests: total.Total,
		MaxBatchDelay:      5 * time.Millisecond,
	}

	limit := 2
	recipes, _ := GetAllRecipes(total.Total, limit)
	if success {
		log.Printf("Total recipes retrieved: %d", len(recipes))
		for _, recipe := range recipes {
			log.Printf("Recipe ID: %s, Name: %s", recipe.ID, recipe.Name)
		}
	} else {
		log.Println("Failed to get all recipes")
	}
}
