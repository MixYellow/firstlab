cat > models/models.go << 'EOF'
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Username  string         `gorm:"unique;not null" json:"username"`
    Password  string         `json:"-"` // не возвращаем в JSON
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    Tasks     []Task         `json:"tasks,omitempty"`
}

type Task struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Title     string         `json:"title"`
    Completed bool           `json:"completed"`
    UserID    uint           `json:"user_id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
EOF

cat > database/database.go << 'EOF'
package database

import (
    "log"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "todo-api/models"
)

var DB *gorm.DB

func InitDB() {
    var err error
    DB, err = gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Автомиграция
    DB.AutoMigrate(&models.User{}, &models.Task{})
    log.Println("Database migrated")
}
EOF



cat > utils/jwt.go << 'EOF'
package utils

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key-change-in-production")

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (uint, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
    if err != nil {
        return 0, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims.UserID, nil
    }
    return 0, errors.New("invalid token")
}
EOF


cat > middleware/auth.go << 'EOF'
package middleware

import (
    "context"
    "net/http"
    "strings"
    "todo-api/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            return
        }

        userID, err := utils.ValidateToken(parts[1])
        if err != nil {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next(w, r.WithContext(ctx))
    }
}
EOF


cat > handlers/auth.go << 'EOF'
package handlers

import (
    "encoding/json"
    "net/http"
    "todo-api/database"
    "todo-api/models"
    "todo-api/utils"
    "golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Хешируем пароль
    hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    user := models.User{
        Username: input.Username,
        Password: string(hashed),
    }
    result := database.DB.Create(&user)
    if result.Error != nil {
        http.Error(w, "Username already exists", http.StatusConflict)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func Login(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    var user models.User
    result := database.DB.Where("username = ?", input.Username).First(&user)
    if result.Error != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": token})
}
EOF


cat > handlers/tasks.go << 'EOF'
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "todo-api/database"
    "todo-api/models"
    "todo-api/middleware"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(uint)

    var input struct {
        Title string `json:"title"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    task := models.Task{
        Title:     input.Title,
        Completed: false,
        UserID:    userID,
    }
    database.DB.Create(&task)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(uint)

    var tasks []models.Task
    database.DB.Where("user_id = ?", userID).Find(&tasks)

    json.NewEncoder(w).Encode(tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(uint)
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var task models.Task
    result := database.DB.Where("id = ? AND user_id = ?", uint(id), userID).First(&task)
    if result.Error != nil {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(uint)
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var task models.Task
    result := database.DB.Where("id = ? AND user_id = ?", uint(id), userID).First(&task)
    if result.Error != nil {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    var input struct {
        Title     string `json:"title"`
        Completed bool   `json:"completed"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if input.Title != "" {
        task.Title = input.Title
    }
    task.Completed = input.Completed

    database.DB.Save(&task)
    json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(uint)
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    result := database.DB.Where("id = ? AND user_id = ?", uint(id), userID).Delete(&models.Task{})
    if result.RowsAffected == 0 {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
EOF


cat > main.go << 'EOF'
package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "todo-api/database"
    "todo-api/handlers"
    "todo-api/middleware"
)

func main() {
    database.InitDB()

    r := mux.NewRouter()

    // Публичные маршруты
    r.HandleFunc("/register", handlers.Register).Methods("POST")
    r.HandleFunc("/login", handlers.Login).Methods("POST")

    // Защищённые маршруты (требуют JWT)
    api := r.PathPrefix("/api").Subrouter()
    api.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            middleware.AuthMiddleware(next.ServeHTTP)(w, r)
        })
    })

    api.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
    api.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
    api.HandleFunc("/tasks/{id}", handlers.GetTask).Methods("GET")
    api.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT")
    api.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
EOF




