package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/gorilla/mux"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

// Files Functions

func readCurrentDir() {
	file, err := os.Open(".")
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		fmt.Println(name)
	}
}

// Routes

// sleep before serving file
func sleepFile(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	seconds := -1
	var err error
	if val, ok := pathParams["seconds"]; ok {
		seconds, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Must be an integer"}`))
			return
		}
	}

	filename := pathParams["filename"]

	// check if POST, read parameters and execute a callback
	if r.Method == http.MethodPost {
		var m map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		callbackURL := m["url"]
		body := m["body"]

		bodyString, _ := json.Marshal(body)
        urlString := fmt.Sprintf("%v", callbackURL)

		fmt.Println("url    : ", callbackURL,
			"\nbody   : ", body)

		go callback(string(urlString), string(bodyString), seconds)
	} else {
		log.Printf("Sleeping %d seconds.\n", seconds)
		time.Sleep(time.Second * time.Duration(seconds))
	}

	http.ServeFile(w, r, filename)
}

func callback(url string, body string, seconds int) {
	log.Printf("Waiting %d seconds.\n", seconds)
	time.Sleep(time.Second * time.Duration(seconds))
	log.Printf("Callback is being called:\n")
	log.Printf("URL:\n%s\n", url)
	log.Printf("Body:\n%s\n", body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(respBody))
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "GET called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "POST called"}`))
}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "PUT called"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "DELETE called"}`))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "Not implemented."}`))
}

func main() {
	flag.Usage = usage
	port := flag.Int("port", 9195, "Listen on defined port.")
	staticDir := flag.String("dir", "static", "Static files directory")
	flag.Parse()

	sPort := ":" + strconv.Itoa(*port)
	ln, err := net.Listen("tcp", sPort)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't listen on port %q: %s", *port, err)
		os.Exit(1)
	}

	err = ln.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't stop listening on port %q: %s", *port, err)
		os.Exit(1)
	}

	//readCurrentDir()

	log.Printf("Start listening on port %s...\n", sPort)
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(*staticDir))))
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", get).Methods(http.MethodGet)
	api.HandleFunc("/", post).Methods(http.MethodPost)
	api.HandleFunc("/", notFound)
	r.HandleFunc("/sleep/{seconds}/{filename}", sleepFile).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/", notFound)

	srv := &http.Server{
		Handler: r,
		Addr:    sPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
