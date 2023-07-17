package handlers

import (
	"database/sql"
	"encoding/json"
	"go-htmx/database"
	"net/http"
)

type Token struct {
	Token string `json:"token"`
}

type AdminData struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Role      string `json:"role"`
}

type UserData struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Role      string `json:"role"`
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	var token Token
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbConn := database.ConnectDB()
	defer dbConn.Close()

	// Prepare the SQL statement to check the role based on the token
	roleQuery := "SELECT role FROM users WHERE token = ?"
	roleStmt, err := dbConn.Prepare(roleQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer roleStmt.Close()

	// Execute the role query
	row := roleStmt.QueryRow(token.Token)

	var role string

	err = row.Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if role != "admin" {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Prepare the SQL statement to retrieve all user data
	query := "SELECT email, firstname, lastname, role FROM users"
	usersStmt, err := dbConn.Prepare(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer usersStmt.Close()

	// Execute the user data query
	usersRows, err := usersStmt.Query()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer usersRows.Close()

	var adminData AdminData
	var userData []UserData

	for usersRows.Next() {
		var email, firstname, lastname, role string

		err := usersRows.Scan(&email, &firstname, &lastname, &role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if role == "admin" {
			adminData = AdminData{
				Email:     email,
				Firstname: firstname,
				Lastname:  lastname,
			}
		} else {
			user := UserData{
				Email:     email,
				Firstname: firstname,
				Lastname:  lastname,
				Role:      role,
			}
			userData = append(userData, user)
		}
	}

	if len(adminData.Email) == 0 {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(map[string]interface{}{
		"admin": adminData,
		"users": userData,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response content type
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.Write(jsonData)

	// tmpl, err := template.ParseFiles("templates/dashboard.html")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// err = tmpl.Execute(w, nil)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}
