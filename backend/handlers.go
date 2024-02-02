package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
    "strconv"
    "strings"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

type contextKey string

var userIDKey = contextKey("userID")

func wsHandler(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	tokenString := params.Get("token")

	if tokenString == "" {
		http.Error(w, "Token is missing", http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("Received: %s\n", message)

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			break
		}
	}
}

func jwtMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header is required", http.StatusUnauthorized)
            return
        }
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            // Ensure algorithm is HMAC and matches what you expect (e.g., "HS256")
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, http.ErrAbortHandler
            }
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Optionally log or handle the UserID for auditing or further checks
        // log.Printf("Authenticated UserID: %d", claims.UserID)

        ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
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

    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

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

    // Modify the response to include the user ID along with the token
    response := map[string]interface{}{
        "token":  tokenString,
        "userID": userID, // Include the user ID in the response
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
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
    userIDStr := vars["user_id"]
    userID, err := strconv.Atoi(userIDStr) // Convert userID to int
    if err != nil {
        log.Printf("Error converting userID to int: %v", err)
        http.Error(w, "Invalid user ID format", http.StatusBadRequest)
        return
    }

    var user User
    err = db.QueryRow("SELECT id, username, firstname, lastname, phone, bio FROM Account WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Phone, &user.Bio)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("No user found with ID: %d", userID)
            http.NotFound(w, r)
            return
        }
        log.Printf("Error querying user from database: %v", err)
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user); err != nil {
        log.Printf("Error encoding user to JSON: %v", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
    }
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

func createChatHandler(w http.ResponseWriter, r *http.Request) {
	var chatRequest ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&chatRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert a new chat record into the Chat table
	// Adjust this as per your Chat table schema, e.g., if you're storing a chat name or other metadata
	res, err := tx.Exec("INSERT INTO Chat (created_at) VALUES (NOW())")
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create chat: "+err.Error(), http.StatusInternalServerError)
		return
	}
	chatID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert participants into the ChatUser table
	for _, participantID := range chatRequest.Participants {
		_, err := tx.Exec("INSERT INTO ChatUser (chat_id, user_id) VALUES (?, ?)", chatID, participantID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to add participant to chat: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Transaction commit failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond to client
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"chat_id": chatID})
}

func getUserIDFromContext(ctx context.Context) (int, bool) {
	// The ctx.Value returns an interface{} and a boolean indicating if the key was found.
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}

func listChatsHandler(w http.ResponseWriter, r *http.Request) {
	// Corrected to handle both return values from getUserIDFromContext
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to retrieve user ID from context", http.StatusInternalServerError)
		return
	}

	query := `
    SELECT c.id, GROUP_CONCAT(u.username SEPARATOR ', ') AS participants
    FROM Chat c
    JOIN ChatUser cu ON c.id = cu.chat_id
    JOIN Account u ON cu.user_id = u.id
    WHERE cu.user_id = ?
    GROUP BY c.id
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		http.Error(w, "Server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var chats []ChatSummary
	for rows.Next() {
		var chat ChatSummary
		if err := rows.Scan(&chat.ID, &chat.Participants); err != nil {
			http.Error(w, "Server error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		chats = append(chats, chat)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

func userIsPartOfChat(userID int, chatID string) bool {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM ChatUser WHERE user_id = ? AND chat_id = ?)"
    err := db.QueryRow(query, userID, chatID).Scan(&exists)
    return err == nil && exists
}

// Ensure you've imported "strconv" at the top of your file.

func getChatHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    chatID := vars["chat_id"]

    userID, ok := getUserIDFromContext(r.Context())
    if !ok {
        http.Error(w, "Failed to retrieve user ID from context", http.StatusUnauthorized)
        return
    }

    // Verify the user is part of the chat
    if !userIsPartOfChat(userID, chatID) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    query := "SELECT m.id, m.sender, m.content, m.created_at FROM Message m WHERE m.chat_id = ?"
    rows, err := db.Query(query, chatID)
    if err != nil {
        http.Error(w, "Server error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var message Message
        if err := rows.Scan(&message.ID, &message.Sender, &message.Content, &message.CreatedAt); err != nil {
            http.Error(w, "Server error: "+err.Error(), http.StatusInternalServerError)
            return
        }
        messages = append(messages, message)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)
}