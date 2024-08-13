package handler

import (
	"crud/internal/domain"
	"crud/internal/pkg/authclient"
	"crud/internal/service"
	"encoding/json"
	"log"
	"strconv"

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
		if token == nil || string(token) == "" {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			log.Println("No token provided")
			return
		}

		userInfo, valid := authclient.GetUserByToken(string(token))
		if !valid {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			log.Println("Invalid token")
			return
		}

		ctx.Request.Header.Set("X-User-ID", userInfo.ID)
		ctx.Request.Header.Set("X-User-Role", userInfo.Role)
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
	pageStr := string(ctx.QueryArgs().Peek("page"))
	limitStr := string(ctx.QueryArgs().Peek("limit"))
	sortByArg := string(ctx.QueryArgs().Peek("sort_by"))

	page := 1
	limit := 2
	sortBy := "name"

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if sortByArg != "name" {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	if limit > 2 {
		limit = 2
	}

	recipes, err := service.GetAll(page, limit, sortBy)
	if err != nil {
		if err.Error() == "page out of range" {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
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

	userID := string(ctx.Request.Header.Peek("X-User-ID"))
	userRole := string(ctx.Request.Header.Peek("X-User-Role"))

	if userID == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("X-User-ID is required in headers")
		return
	}

	if userRole == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("X-User-Role is required in headers")
		return
	}

	rec, err := service.Get(string(id))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	if rec.AuthorID != userID && userRole != "admin" {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		ctx.SetBodyString("You don't have permission to delete this recipe")
		return
	}

	if err := service.Delete(string(id)); err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func PostHandler(ctx *fasthttp.RequestCtx) {
	var input RecipeReq

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	userID := string(ctx.Request.Header.Peek("X-User-ID"))
	if userID == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("X-User-ID is required in headers")
		return
	}

	input.AuthorID = userID

	rec := domain.Recipe{
		ID:          input.ID,
		AuthorID:    input.AuthorID,
		Name:        input.Name,
		Ingredients: input.Ingredients,
		Temperature: input.Temperature,
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
