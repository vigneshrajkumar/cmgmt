package main

import (
	"cmgmt/datastore"
	"fmt"
	"html/template"

	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func handlePOSTLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println("1")
	valid, err := store.ValidateUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(valid)
	if !valid {
		tmpl := template.Must(template.ParseFiles("./public/landing.html"))
		tmpl.Execute(w, datastore.ErrorData{ErrorMessage: "Invalid login credentials"})
		return
	}

	sessionToken, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error while generating session token.", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Println(sessionToken)

	if err := store.UpdateSessionToken(username, sessionToken.String()); err != nil {
		fmt.Println("error while updating session token.", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: time.Now().Add(30 * time.Minute),
	})

	fmt.Println("auth success... redirecting to members...")
	http.Redirect(w, r, "/members", http.StatusSeeOther)
	return
}

// func handleLogin(w http.ResponseWriter, r *http.Request) {

// 	switch r.Method {
// 	case "GET":
// 		fmt.Println("GET /")
// 		handleGETLogin(w, r)
// 		return

// 	case "POST":
// 		fmt.Println("POST /")
// 		handlePOSTLogin(w, r)
// 		return
// 	}
// }

func handleLogout(w http.ResponseWriter, r *http.Request) {
	// You logout by setting negative time
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "expired",
		Expires: time.Now().Add(-7 * 24 * time.Hour),
	})
	fmt.Println("logout success... redirecting to landing page...")
	http.Redirect(w, r, "/", http.StatusOK)
}
