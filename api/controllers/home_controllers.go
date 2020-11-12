package controllers

import (
	"net/http"

	"github.com/bayuajik2403/golang-backend-crud-api/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome")

}
