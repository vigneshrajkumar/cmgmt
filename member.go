package main

import (
	"cmgmt/datastore"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"net/http"
)

func handleMemberByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("GET /member/id")
		handleGETMemberByID(w, r)

	case "DELETE":
		fmt.Println("DELETE /member/id")
		handleDELETEMemberByID(w, r)

	default:
		fmt.Println("Route not handled ", r.Method)
	}
}

func handleMembers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("GET /members")
		handleGETMembers(w, r)

	case "POST":
		fmt.Println("POST /members")
		handlePOSTMembers(w, r)

	default:
		fmt.Println("Route not handled ", r.Method)
	}
}

func handleGETMembers(w http.ResponseWriter, r *http.Request) {

	mems, err := store.GetMembers()
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(mems)

	js, err := json.Marshal(mems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handleGETMemberByID(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)
	splitPath := strings.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(splitPath[len(splitPath)-1])
	if err != nil {
		fmt.Println(err)
	}

	mems, err := store.GetMemberByID(int64(id))
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(mems)

	js, err := json.Marshal(mems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func handleDELETEMemberByID(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)

	splitPath := strings.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(splitPath[len(splitPath)-1])
	if err != nil {
		fmt.Println(err)
	}

	if err := store.DeleteMember(int64(id)); err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusCreated)
	js, err := json.Marshal(map[string]string{"status": "deleted"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createMemberWithoutIDFromRequest(r *http.Request) (*datastore.Member, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	birthday := r.FormValue("birthday")
	gender := r.FormValue("gender")
	fID := r.FormValue("fID")
	if fID == "" {
		fID = "-1"
	}

	parsedBd, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return nil, err
	}

	fidStr, err := strconv.Atoi(fID)
	if err != nil {
		return nil, err
	}

	return datastore.NewMember(store.NextID(), firstName, lastName, parsedBd, gender, int64(fidStr)), nil
}

func handlePOSTMembers(w http.ResponseWriter, r *http.Request) {

	m, err := createMemberWithoutIDFromRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m.ID = store.NextID()

	fmt.Println(m)

	err = store.AddMember(m)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	js, err := json.Marshal(map[string]string{"status": "created"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
