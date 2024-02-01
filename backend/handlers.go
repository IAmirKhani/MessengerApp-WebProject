package main

import (
    "encoding/json"
    "net/http"
    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
    "database/sql"
    "time"
)

type UserRegistrationRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type Claims struct {
    UserID int `json:"userId"`
    jwt.StandardClaims
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    var request UserRegistrationRequest
    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Validate data 
    if request.Username == "" || request.Password == "" {
        http.Error(w, "Username and password are required", http.StatusBadRequest)
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error while hashing password", http.StatusInternalServerError)
        return
    }

    // Insert user into the database
    _, err = db.Exec("INSERT INTO Account (username, password) VALUES (?, ?)", request.Username, string(hashedPassword))
    if err != nil {
        
        http.Error(w, "Error while inserting user into database: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond to client
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}


var jwtKey = []byte("your_secret_key_here") // Replace 

func loginHandler(w http.ResponseWriter, r *http.Request) {
    var creds UserRegistrationRequest 
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var storedPassword string
    var userID int
    err = db.QueryRow("SELECT id, password FROM Account WHERE username = ?", creds.Username).Scan(&userID, &storedPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        } else {
            http.Error(w, "Server error", http.StatusInternalServerError)
        }
        return
    }

    // Compare the stored hashed password with the provided password
    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Generate JWT Token
    expirationTime := time.Now().Add(30 * time.Minute)
    claims := &Claims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)

    if err != nil {
        http.Error(w, "Could not create token", http.StatusInternalServerError)
        return
    }

    // Return the token
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// handlers.go
func getContactsHandler(w http.ResponseWriter, r *http.Request) {

    vars := mux.Vars(r)
    userID := vars["user_id"]


    rows, err := db.Query("SELECT id, user_id, contact_id, contact_name FROM ContactList WHERE user_id = ?", userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var contacts []Contact
    for rows.Next() {
        var contact Contact
        if err := rows.Scan(&contact.ID, &contact.UserID, &contact.ContactID, &contact.ContactName); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        contacts = append(contacts, contact)
    }


    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(contacts)
}

func addContactHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]

    var contact Contact
    if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err := db.Exec("INSERT INTO ContactList (user_id, contact_id, contact_name) VALUES (?, ?, ?)", userID, contact.ContactID, contact.ContactName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Contact added successfully"})
}

func removeContactHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]
    contactID := vars["contact_id"]

    _, err := db.Exec("DELETE FROM ContactList WHERE user_id = ? AND contact_id = ?", userID, contactID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Contact removed successfully"})
}


func getUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]

    var user User
    err := db.QueryRow("SELECT id, username, firstname, lastname, phone, bio FROM Account WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Phone, &user.Bio)
    if err != nil {
        if err == sql.ErrNoRows {
            http.NotFound(w, r)
            return
        }
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}


func updateUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]


    var updateReq UserUpdateRequest
    if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

   
    _, err := db.Exec("UPDATE Account SET firstname = ?, lastname = ?, bio = ? WHERE id = ?", updateReq.Firstname, updateReq.Lastname, updateReq.Bio, userID)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Account updated successfully"})
}


func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]

    _, err := db.Exec("DELETE FROM Account WHERE id = ?", userID)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted successfully"})
}


func searchUserHandler(w http.ResponseWriter, r *http.Request) {
    keyword := r.URL.Query().Get("keyword")

    rows, err := db.Query("SELECT id, username, firstname, lastname FROM Account WHERE username LIKE ? OR firstname LIKE ? OR lastname LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []User 
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname); err != nil {
            http.Error(w, "Server error", http.StatusInternalServerError)
            return
        }
        users = append(users, user)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

