package main

import (
    "time"
    "github.com/dgrijalva/jwt-go"
)

type Contact struct {
    ID          int    `json:"id"`
    UserID      int    `json:"user_id"`
    ContactID   int    `json:"contact_id"`
    ContactName string `json:"contact_name"`
}


type UserUpdateRequest struct {
    Firstname string `json:"firstname"`
    Lastname  string `json:"lastname"`
    Bio       string `json:"bio"`
}

type UserRegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChatRequest struct {
	Participants []int `json:"participants"` // IDs of the users participating in the chat
}

type Message struct {
    ID        int       `json:"id"`
    Sender    int       `json:"sender"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Firstname string `json:"firstname"`
    Lastname  string `json:"lastname"`
    Phone     string `json:"phone"`
    Bio       string `json:"bio"`
}

type ChatSummary struct {
	ID           int    `json:"id"`
	Participants string `json:"participants"` // Assumes participants are a comma-separated string
}

type Claims struct {
	UserID int `json:"userId"`
	jwt.StandardClaims
}
