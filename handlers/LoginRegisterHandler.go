package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"go-htmx/database"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email           string
	Password        string
	ConfirmPassword string
	FirstName       string
	LastName        string
	Role            string
	Token           string
}

func GenerateToken(email string) string {
	tokenBytes := make([]byte, 32)
	_, _ = rand.Read(tokenBytes)

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	tokenData := email + timestamp + string(tokenBytes)

	token := base64.StdEncoding.EncodeToString([]byte(tokenData))
	return token
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func PrintJSONResponse(data interface{}, statusCode int) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Printf("Response (Status: %d): %s\n", statusCode, string(jsonData))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		db := database.ConnectDB()
		defer db.Close()

		query := "SELECT email, password, role, token FROM users WHERE email=?"
		var user User

		err := db.QueryRow(query, email).Scan(&user.Email, &user.Password, &user.Role, &user.Token)
		if err != nil {
			log.Println("Invalid email or password:", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			log.Println("Invalid email or password:", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if user.Token == "" {
			log.Println("Token not found for the user:", user.Email)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		switch user.Role {
		case "admin":
			setTokenAndRedirect(w, r, user.Token, "/dashboard")
		case "user":
			setTokenAndRedirect(w, r, user.Token, "/home")
		default:
			http.Error(w, "Invalid role", http.StatusInternalServerError)
		}

	} else {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func setTokenAndRedirect(w http.ResponseWriter, r *http.Request, token string, redirectURL string) {
	if token == "" {
		log.Println("Token is empty for the user")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	script := fmt.Sprintf(`<script>localStorage.setItem('token', '%s'); window.location.href = '%s';</script>`, token, redirectURL)
	fmt.Fprint(w, script)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		role := "user" // Set the default role to "user"

		if password != confirmPassword {
			log.Println("Password and Confirm Password do not match")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Generate bcrypt hash of the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error generating bcrypt hash:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		db := database.ConnectDB()
		defer db.Close()

		query := "SELECT COUNT(*) FROM users WHERE email=?"
		var count int
		err = db.QueryRow(query, email).Scan(&count)
		if err != nil {
			log.Println("Error checking user existence:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			log.Println("User already exists with the provided email")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		token := GenerateToken(email)

		insertQuery := "INSERT INTO users (email, password, firstname, lastname, role, token) VALUES (?, ?, ?, ?, ?, ?)"
		_, err = db.Exec(insertQuery, email, hashedPassword, firstName, lastName, role, token)
		if err != nil {
			log.Println("Error inserting user:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
