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

func MainPage(w http.ResponseWriter, r *http.Request) {
	ip := getIpFromRequest(r)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{"ip": ip}); err != nil {
		log.Println("There was an error with header X-Forwarded-For parsing:", err)
	}
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
