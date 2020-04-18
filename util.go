package main

import (
	"cmgmt/datastore"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
)

func handleGETSearch(w http.ResponseWriter, r *http.Request) {

	filters := make([]*datastore.Filter, 0)
	for key, val := range r.URL.Query() {
		fmt.Println(key, "=>", val)

		switch key {
		case "first-name-starts-with":
			filters = append(filters, datastore.NewFilter("first_name", "ILIKE", val[0]+"%"))
		case "first-name-ends-with":
			filters = append(filters, datastore.NewFilter("first_name", "ILIKE", "%"+val[0]))
		case "first-name-is":
			filters = append(filters, datastore.NewFilter("first_name", "=", val[0]))
		default:
			fmt.Println("unhandled fitler processing")
		}

	}

	mems, err := store.GetMembers(filters...)
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
