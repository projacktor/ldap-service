package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type authorizedUser struct {
	role string
	name string
}

func getIpFromRequest(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return ips[0]
	}
	return ""
}

func AuthorizedPage(w http.ResponseWriter, r *http.Request) {
	response := authorizedUser{
		role: "admin",
		name: "Jonh",
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("There was an error in authorization", err)
	}
}
