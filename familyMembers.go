package main

import (
	"cmgmt/datastore"
	"encoding/json"
	"fmt"
	"net/http"
)

func handleFamilyMembers(w http.ResponseWriter, r *http.Request) {
	memberID := store.NextID()
	familyID := store.NextID()

	f := datastore.NewFamily(familyID, memberID)

	m, err := createMemberWithoutIDFromRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m.ID = memberID
	m.FamilyID = familyID

	// TODO: The below code should be in a transaction

	if err := store.AddMember(m); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := store.AddFamily(f); err != nil {
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
