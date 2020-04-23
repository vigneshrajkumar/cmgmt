package main

import (
	"cmgmt/datastore"
	"flag"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

var store *datastore.Store
var logger *log.Logger

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

	if err := os.Remove("cm.log"); err != nil {
		fmt.Println(err)
	}

	logFile, err := os.OpenFile("cm.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("unable to create logfile")
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetPrefix("cm+ ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Church Management Suite")
	log.Println("-----------------------")

	reinitDB := flag.Bool("reinit", false, "reinitialize database")
	flag.Parse()

	ds, err := datastore.NewStore()
	if err != nil {
		log.Println("new ds error.", err)
	}
	store = ds

	if *reinitDB {
		log.Println("reinitialzing database")

		err = store.Reset()
		if err != nil {
			log.Println(err)
		}
		err = store.Initialize()
		if err != nil {
			log.Println("init error", err)
		}

		err = store.EstablishAdminAccess()
		if err != nil {
			log.Println("error establishing admin access", err)
		}
	}

	r := mux.NewRouter()

	r.HandleFunc("/pdf/member/{mid}", handleGETPDFMember)

	r.HandleFunc("/search", handleGETSearch)

	r.HandleFunc("/util/professions", handleGETProfessionsList)
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
