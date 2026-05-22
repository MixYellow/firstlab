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


