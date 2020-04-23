package main

import (
	"cmgmt/datastore"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"net/http"
)

func handleMemberByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("GET /member/id")
		handleGETMemberByID(w, r)

	case "DELETE":
		log.Println("DELETE /member/id")
		handleDELETEMemberByID(w, r)

	default:
		log.Println("Route not handled ", r.Method)
	}
}

func handleMembers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("GET /members")
		handleGETMembers(w, r)

	case "POST":
		log.Println("POST /members")
		handlePOSTMembers(w, r)

	default:
		log.Println("Route not handled ", r.Method)
	}
}

func handleGETMembers(w http.ResponseWriter, r *http.Request) {
	log.Println("handleGETMembers()")

	mems, err := store.GetMembers()
	if err != nil {
		log.Println("getMembers() error ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(mems)

	js, err := json.Marshal(mems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handleGETMemberByID(w http.ResponseWriter, r *http.Request) {
	splitPath := strings.Split(r.URL.Path, "/")
	id, err := strconv.ParseFloat(splitPath[len(splitPath)-1], 64)
	if err != nil {
		log.Println(err)
	}

	mem, err := store.GetMemberByID(id)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(mem)

	prof, err := store.ResolveProfession(mem.ProfessionID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(mem)

	payload := make(map[string]interface{})
	payload["member"] = mem
	payload["resolvedProfession"] = prof

	js, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func handleDELETEMemberByID(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL.Path)

	splitPath := strings.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(splitPath[len(splitPath)-1])
	if err != nil {
		log.Println(err)
	}

	if err := store.DeleteMember(int64(id)); err != nil {
		log.Println(err)
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

func handlePOSTMembers(w http.ResponseWriter, r *http.Request) {
	log.Println("handlePOSTMembers()")

	decoder := json.NewDecoder(r.Body)
	var m datastore.Member
	err := decoder.Decode(&m)
	if err != nil {
		log.Println(err)
	}
	// log.Println("FirstName : ", m.FirstName)

	log.Println("FirstName : ", m.FirstName)
	log.Println("LastName : ", m.LastName)

	m.ID = float64(int(store.NextID()))

	log.Println(m)

	err = store.AddMember(&m)
	if err != nil {
		log.Println(err)
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
