package controllers

import (
	"net/http"

	"github.com/bayuajik2403/go-crud-api/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome")

}
