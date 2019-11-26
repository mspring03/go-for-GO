package main

import (
	"./func"
	"fmt"
	"net/http"
)

func main() {
	r := net.Router{Handlers: make(map[string]map[string]net.HandlerFunc)}

	r.HandleFunc("GET", "/", net.LogHandler(net.RecoverHandler(net.ParseFormHandler(net.ParseJsonBodyHandler(func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "welcome!", c.Params)
	})))))

	r.HandleFunc("GET", "/about", net.LogHandler(net.RecoverHandler(net.ParseFormHandler(net.ParseJsonBodyHandler(func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "about!", c.Params)
	})))))

	r.HandleFunc("GET", "/users/:user_id", net.LogHandler(net.RecoverHandler(net.ParseFormHandler(net.ParseJsonBodyHandler(func(c *net.Context) {
		if c.Params["user_id"] == "0" {
			panic("id is zero")
		}
		fmt.Fprintf(c.ResponseWriter, "recieve user %v\n", c.Params["user_id"])
	})))))

	r.HandleFunc("GET", "/users/:user_id/addresses/:address_id", net.LogHandler(net.RecoverHandler(net.ParseFormHandler(net.ParseJsonBodyHandler(func(c *net.Context) {
		fmt.Fprintf(c.ResponseWriter, "retrieve user %v's address %v\n", c.Params["user_id"], c.Params["address_id"])
	})))))

	r.HandleFunc("POST", "/users", net.LogHandler(net.RecoverHandler(net.ParseFormHandler(net.ParseJsonBodyHandler(func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "welcome!", c.Params)
	})))))

	r.HandleFunc("POST", "/users/:user_id/addresses", net.LogHandler(net.RecoverHandler(net.ParseFormHandler(net.ParseJsonBodyHandler(func(c *net.Context) {
		fmt.Fprintf(c.ResponseWriter, "create user %v's address", c.Params["user_id"])
	})))))

	http.ListenAndServe(":8080", r)
}