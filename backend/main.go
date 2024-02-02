package main

import (
    "net/http"
    "database/sql"
    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("mysql", "root:xCjUmBYEXZ1Kptj8rJsdAT6I@tcp(kilimanjaro.liara.cloud:32794)/MessengerApp")
    if err != nil {
        log.Fatal(err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }
}

func main() {
    log.Println("Initializing database...")
    initDB()
    log.Println("Database initialized.")
    router := mux.NewRouter()

    // Public endpoints
    router.HandleFunc("/api/register", registerHandler).Methods("POST")
    router.HandleFunc("/api/login", loginHandler).Methods("POST")

    // Subrouter for protected routes
    protected := router.PathPrefix("/api").Subrouter()
    protected.Use(jwtMiddleware) // Apply jwtMiddleware to all routes handled by "protected"
    
    // Apply the middleware to routes that require authentication
    protected.HandleFunc("/users/{user_id}", getUserHandler).Methods("GET")
    protected.HandleFunc("/users/{user_id}", updateUserHandler).Methods("PATCH")
    protected.HandleFunc("/users/{user_id}", deleteUserHandler).Methods("DELETE")
    protected.HandleFunc("/users", searchUserHandler).Methods("GET")
    protected.HandleFunc("/users/{user_id}/contacts", getContactsHandler).Methods("GET")
    protected.HandleFunc("/users/{user_id}/contacts", addContactHandler).Methods("POST")
    protected.HandleFunc("/users/{user_id}/contacts/{contact_id}", removeContactHandler).Methods("DELETE")
    protected.HandleFunc("/chats", createChatHandler).Methods("POST")
    protected.HandleFunc("/chats", listChatsHandler).Methods("GET")
    protected.HandleFunc("/chats/{chat_id}", getChatHandler).Methods("GET") // Assuming you want to protect this route

    // WebSocket route, assuming it needs authentication
    router.Handle("/ws", jwtMiddleware(http.HandlerFunc(wsHandler)))

    // Start the server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal("ListenAndServe error: ", err)
    }
}
