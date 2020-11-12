package controllers

import "github.com/bayuajik2403/golang-backend-crud-api/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Products routes
	s.Router.HandleFunc("/product", middlewares.SetMiddlewareJSON(s.CreateProduct)).Methods("POST")
	s.Router.HandleFunc("/product", middlewares.SetMiddlewareJSON(s.GetProducts)).Methods("GET")
	s.Router.HandleFunc("/product/{id}", middlewares.SetMiddlewareJSON(s.GetProduct)).Methods("GET")
	s.Router.HandleFunc("/product/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateProduct))).Methods("PUT")
	s.Router.HandleFunc("/product/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteProduct)).Methods("DELETE")
	s.Router.HandleFunc("/product/find/{id}", middlewares.SetMiddlewareJSON(s.GetProductByUser)).Methods("GET")

	//Transactions routes
	s.Router.HandleFunc("/transaction", middlewares.SetMiddlewareJSON(s.CreateTransaction)).Methods("POST")
	s.Router.HandleFunc("/transaction", middlewares.SetMiddlewareJSON(s.GetTransactions)).Methods("GET")
	s.Router.HandleFunc("/transaction/{id}", middlewares.SetMiddlewareJSON(s.GetTransaction)).Methods("GET")
	s.Router.HandleFunc("/transaction/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateTransaction))).Methods("PUT")
	s.Router.HandleFunc("/transaction/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteTransaction)).Methods("DELETE")
	s.Router.HandleFunc("/transaction/find/{id}", middlewares.SetMiddlewareJSON(s.GetTransactionByUser)).Methods("GET")

}
