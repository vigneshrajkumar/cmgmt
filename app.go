package main

import (
	"cmgmt/datastore"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

var store *datastore.Store

func setEnv() {
	os.Setenv("HOST", "127.0.0.1")
	os.Setenv("PORT", "6001")
	os.Setenv("DATABASE", "membership")
	os.Setenv("USER", "postgres")
	os.Setenv("PASSWORD", "Stonebraker")
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

	r := mux.NewRouter()

	r.HandleFunc("/search", handleGETSearch)

	r.HandleFunc("/util/families", handleGETFamilyList)
	r.HandleFunc("/util/member/{mid}/family", handleGETMemberFamily)

	r.HandleFunc("/family/members", handleFamilyMembers)
	r.HandleFunc("/members", handleMembers)
	r.HandleFunc("/members/{key}", handleMemberByID)
	r.HandleFunc("/families", handleFamilies)
	r.HandleFunc("/families/{key}", handleFamiliesByID)
	http.Handle("/", r)

	fs := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	log.Println("Listening on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
