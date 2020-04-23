package main

import (
	"cmgmt/datastore"
	"encoding/json"
	"log"
	"net/http"
)

func handleFamilyMembers(w http.ResponseWriter, r *http.Request) {
	log.Println("handleFamilyMembers()")
	log.Println("content type: ", r.Header.Get("Content-type"))

	memberID := store.NextID()
	familyID := store.NextID()

	f := datastore.NewFamily(familyID, memberID)

	decoder := json.NewDecoder(r.Body)
	var m datastore.Member
	err := decoder.Decode(&m)
	if err != nil {
		log.Println("decoding eror", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m.ID = float64(int(memberID))
	m.FamilyID = float64(int(familyID))

	log.Println("Member: ", m)

	// TODO: The below code should be in a transaction

	if err := store.AddMember(&m); err != nil {
		log.Println("Err adding member: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := store.AddFamily(f); err != nil {
		log.Println("Err adding family: ", err)
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
