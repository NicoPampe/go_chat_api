package controller

import (
	"fmt"
  "log"
  "encoding/json"
	"net/http"
	"strconv"

  "github.com/gin-gonic/gin"
  "github.com/swaggo/files"
  "github.com/swaggo/gin-swagger"


	// "github.com/swaggo/swag/example/celler/httputil"
	// "github.com/swaggo/swag/example/celler/model"
)

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

// ShowUser godoc
// @Summary Show a account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} user.User
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} error.HTTPError
// @Failure 500 {object} error.HTTPError
// @Router /accounts/{id} [get]
func (c *Controller) ShowUser(ctx *gin.Context) {
	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		error.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	account, err := user.UserOne(aid)
	if err != nil {
		error.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, account)
}

// ListUsers godoc
// @Summary List accounts
// @Description get accounts
// @Accept  json
// @Produce  json
// @Param q query string false "name search by q"
// @Success 200 {array} user.User
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} error.HTTPError
// @Failure 500 {object} error.HTTPError
// @Router /accounts [get]
func (c *Controller) ListUsers(ctx *gin.Context) {
	q := ctx.Request.URL.Query().Get("q")
	accounts, err := user.UsersAll(q)
	if err != nil {
		error.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
