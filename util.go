package main

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
)

func handleGETFamilyList(w http.ResponseWriter, r *http.Request) {

	fams, err := store.GetFamilyNamesWithID()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	payload := make(map[string]map[int64]string)
	payload["families"] = fams

	js, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handleGETMemberFamily(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fams, err := store.GetMemberFamilyInfo(vars["mid"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// payload := make(map[string][]interface{})
	// payload["families"] = fams

	fmt.Println(fams)

	js, err := json.Marshal(fams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
