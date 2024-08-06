package handler

import (
	"crud/internal/domain"
	"crud/internal/pkg/authclient"
	"crud/internal/service"
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
)

type RequestHandler func(*fasthttp.RequestCtx)

var routeHandlers = map[string]map[string]RequestHandler{
	"/get": {
		fasthttp.MethodGet: GetHandler,
	},
	"/get_all": {
		fasthttp.MethodGet: GetAllHandler,
	},
	"/delete": {
		fasthttp.MethodDelete: DeleteHandler,
	},
	"/post": {
		fasthttp.MethodPost: PostHandler,
	},
}

func ServerHandler(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, "*")
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodPost)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodGet)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodDelete)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowHeaders, fasthttp.HeaderContentType)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowHeaders, fasthttp.HeaderAuthorization)

	if ctx.IsOptions() {
		return
	}

	path := string(ctx.Path())
	method := string(ctx.Method())

	if path != "/get_all" {
		token := ctx.Request.Header.Peek(fasthttp.HeaderAuthorization)
		log.Println(string(token) == "", !authclient.ValidateToken(string(token)), string(token) == "" || !authclient.ValidateToken(string(token)))
		if string(token) == "" || !authclient.ValidateToken(string(token)) {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			log.Println("Get request", string(ctx.Method()), string(token), "error", fasthttp.StatusUnauthorized)
			return
		}
	}

	if methodHandlers, ok := routeHandlers[path]; ok {
		if handler, ok := methodHandlers[method]; ok {
			handler(ctx)
			return
		}
	}

	ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)

}

func GetHandler(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().Peek("id")
	if len(id) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	rec, err := service.Get(string(id))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	marshal, err := json.Marshal(rec)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if _, err = ctx.Write(marshal); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetAllHandler(ctx *fasthttp.RequestCtx) {
	recipes, err := service.GetAll()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	marshal, err := json.Marshal(recipes)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if _, err = ctx.Write(marshal); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func DeleteHandler(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().Peek("id")
	if len(id) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if err := service.Delete(string(id)); err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func PostHandler(ctx *fasthttp.RequestCtx) {
	var rec domain.Recipe
	log.Println(string(ctx.PostBody()))
	if err := json.Unmarshal(ctx.PostBody(), &rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if err := service.AddOrUpd(&rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	resp := IdResponse{ID: rec.ID}

	marshal, err := json.Marshal(resp)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if _, err = ctx.Write(marshal); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
