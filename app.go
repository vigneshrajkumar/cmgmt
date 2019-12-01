package main

import (
	"cmgmt/datastore"
	"fmt"
	"html/template"
	"os"

	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

var store *datastore.Store

func setEnv() {
	os.Setenv("HOST", "127.0.0.1")
	os.Setenv("PORT", "5432")
	os.Setenv("DATABASE", "membership")
	os.Setenv("USER", "postgres")
	os.Setenv("PASSWORD", "emc")
}

func clearEnv() {
	os.Clearenv()
}

func main() {

	defer clearEnv()
	setEnv()

	fmt.Println("CM+")

	ds, err := datastore.NewStore()
	if err != nil {
		fmt.Println("new ds error.", err)
	}
	store = ds

	if len(os.Args) != 2 {
		fmt.Println("usage: cmgmt [reinit|run]")
		return
	}

	switch os.Args[1] {
	case "reinit":
		fmt.Println("reinit system..")

		err = store.Reset()
		if err != nil {
			fmt.Println(err)
		}
		err = store.Initialize()
		if err != nil {
			fmt.Println("init error", err)
		}

		err = store.EstablishAdminAccess()
		if err != nil {
			fmt.Println("error establishing admin access", err)
		}
	case "run":
		fmt.Println("run system..")
	default:
		fmt.Println("usage: cmgmt [reinit|run]")
		return
	}

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)

	http.HandleFunc("/dashboard/", handleDashboard)

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":3000", nil)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("GET /dashboard")
		handleGETDashboard(w, r)
		return

	default:
		fmt.Println("Route not handled ", r.Method)
	}
}

func validateSession(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("no cookie")
			tmpl := template.Must(template.ParseFiles("./public/error.html"))
			tmpl.Execute(w, datastore.ErrorData{ErrorMessage: "Session Expired, Login Again."})
			return
		}
		fmt.Println("bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handleGETDashboard(w http.ResponseWriter, r *http.Request) {

	validateSession(w, r)

	cookie, _ := r.Cookie("session_token")
	username, err := store.GetUser(cookie.Value)
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("logged in as ", username)
	tmpl := template.Must(template.ParseFiles("./public/dashboard.html"))
	tmpl.Execute(w, nil)
	return
}

func handleGETLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./public/landing.html"))
	tmpl.Execute(w, nil)
	return
}

func handlePOSTLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	valid, err := store.ValidateUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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

	if err := store.UpdateSessionToken(username, sessionToken.String()); err != nil {
		fmt.Println("error while updating session token.", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: time.Now().Add(30 * time.Minute),
	})

	fmt.Println("auth success... redirecting to dashboard...")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	return
}

func handleLogin(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		fmt.Println("GET /")
		handleGETLogin(w, r)
		return

	case "POST":
		fmt.Println("POST /")
		handlePOSTLogin(w, r)
		return
	}
}

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
