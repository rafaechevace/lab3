package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// In-memory stores and synchronization primitives

// userStore maps username -> password hash.
var userStore = make(map[string]string)

// tokenStore maps token -> TokenInfo.
var tokenStore = make(map[string]TokenInfo)

// Mutexes to protect concurrent access to the maps.
var userMutex = &sync.Mutex{}
var tokenMutex = &sync.Mutex{}

// Filename used to persist users to disk.
const userFile = "users.json"

// UserCredentials represents the JSON body for signup/login requests.
type UserCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenInfo stores the username associated with a token and its expiry time.
type TokenInfo struct {
	Username  string
	ExpiresAt time.Time
}

// init loads persisted users on package initialization.
func init() {
	loadUsers()
}

// Signup registers a new user, persists it, and returns an access token.
func Signup(c *gin.Context) {
	var creds UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing fields or invalid format"})
		return
	}

	userMutex.Lock()
	defer userMutex.Unlock()

	// Check if username already exists.
	if _, exists := userStore[creds.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	// Hash the password before storing.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process password"})
		return
	}

	// Save the new user and persist to disk.
	userStore[creds.Username] = string(hashedPassword)
	if err := saveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save user"})
		return
	}

	// Generate and return an access token.
	accessToken := generateToken(creds.Username)
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

// Login authenticates a user and returns a new access token on success.
func Login(c *gin.Context) {
	var creds UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing fields or invalid format"})
		return
	}

	userMutex.Lock()
	defer userMutex.Unlock()

	// Retrieve stored hash for the username.
	hashedPassword, ok := userStore[creds.Username]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Compare provided password with stored hash.
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate and return a new access token.
	accessToken := generateToken(creds.Username)
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

// generateToken creates a cryptographically random token and stores it with an expiry.
func generateToken(username string) string {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	// Generate 32 random bytes for the token.
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		log.Printf("error generating token: %v", err)
		// Fallback to a less secure token if crypto randomness fails.
		return "fallback_token_" + time.Now().String()
	}

	token := hex.EncodeToString(randomBytes)

	// Store token with a 5-minute expiration.
	tokenStore[token] = TokenInfo{
		Username:  username,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	return token
}

// saveUsers persists the userStore map to a JSON file.
func saveUsers() error {
	data, err := json.MarshalIndent(userStore, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(userFile, data, 0644)
}

// loadUsers loads users from the JSON file into memory.
func loadUsers() {
	userMutex.Lock()
	defer userMutex.Unlock()

	data, err := os.ReadFile(userFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("file %s does not exist. It will be created on first registration.", userFile)
			return
		}
		log.Printf("error reading users file: %v", err)
		return
	}

	if err := json.Unmarshal(data, &userStore); err != nil {
		log.Printf("error decoding users file: %v", err)
	}
}
