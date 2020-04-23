package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"net/http"
)

func handleFamilies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("GET /families")
		handleGETFamilies(w, r)

	// case "POST":
	// 	fmt.Println("POST /families")
	// 	handlePOSTFamilies(w, r)

	default:
		fmt.Println("Route not handled ", r.Method)
	}
}

func handleFamiliesByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("GET /member/id")
		handleGETFamilyByID(w, r)

	// case "DELETE":
	// 	fmt.Println("DELETE /member/id")
	// 	handleDELETEMemberByID(w, r)

	default:
		fmt.Println("Route not handled ", r.Method)
	}
}

func handleGETFamilyByID(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)
	splitPath := strings.Split(r.URL.Path, "/")

	id, err := strconv.ParseFloat(splitPath[len(splitPath)-1], 64)
	if err != nil {
		fmt.Println(err)
	}

	mems, err := store.GetMemberByID(id)
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

func handleGETFamilies(w http.ResponseWriter, r *http.Request) {

	fams, err := store.GetFamilies()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(fams)

	js, err := json.Marshal(fams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
