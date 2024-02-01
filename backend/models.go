package main

type Contact struct {
    ID          int    `json:"id"`
    UserID      int    `json:"user_id"`
    ContactID   int    `json:"contact_id"`
    ContactName string `json:"contact_name"`
}

type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Password  string `json:"password,omitempty"`
    Firstname string `json:"firstname"`
    Lastname  string `json:"lastname"`
    Phone     string `json:"phone"`
    Bio       string `json:"bio"`
}

type UserUpdateRequest struct {
    Firstname string `json:"firstname"`
    Lastname  string `json:"lastname"`
    Bio       string `json:"bio"`
}