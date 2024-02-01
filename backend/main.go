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

    // User endpoints
    router.HandleFunc("/api/register", registerHandler).Methods("POST")
    router.HandleFunc("/api/login", loginHandler).Methods("POST")
    router.HandleFunc("/api/users/{user_id}", getUserHandler).Methods("GET")
    router.HandleFunc("/api/users/{user_id}", updateUserHandler).Methods("PATCH")
    router.HandleFunc("/api/users/{user_id}", deleteUserHandler).Methods("DELETE")
    router.HandleFunc("/api/users", searchUserHandler).Methods("GET")
    router.HandleFunc("/api/users/{user_id}/contacts", getContactsHandler).Methods("GET")
    router.HandleFunc("/api/users/{user_id}/contacts", addContactHandler).Methods("POST")
    router.HandleFunc("/api/users/{user_id}/contacts/{contact_id}", removeContactHandler).Methods("DELETE")


    // Start the server
    http.ListenAndServe(":8080", router)

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
