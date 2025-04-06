package sso

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	db "main/database"
	"main/helpers/jwtHelper"
	"main/helpers/userHelper"
	"net/http"
	"time"
)

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	var token jwtHelper.TokenPair
	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	v, _ := userHelper.GetUserByToken(token, ctx)
	if v != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	token, err = token.GenerateAndValidateToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := token.ToJson()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	creds := struct {
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	user, err := userHelper.GetUserByPassword(ctx, creds.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	var token jwtHelper.TokenPair
	token.NewTokens(jwtHelper.UserClaims{UserId: user.Id, Role: 1})
	err = user.UpdateToken(ctx, token)

	if err != nil {
		http.Error(w, "Fail update token.", http.StatusBadRequest)
		return
	}

	data, err := token.ToJson()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, data)
}

func main() {
	db.SetupDatabase()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /refresh", refreshHandler)

	mux.HandleFunc("POST /login", loginHandler)

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Fatal(err)
	}
	fmt.Println("server started")
}
