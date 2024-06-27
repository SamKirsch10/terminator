package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	server   = flag.Bool("server", false, "Set to server mode")
	port     = flag.String("port", "8111", "Port the server will listen on or the client will request")
	authFile = flag.String("file", "", "Location of authorized_keys file to serve.")

	client              = flag.Bool("client", false, "Set to client mode")
	user                = flag.String("user", "sam", "User to apply authorized_keys file in client mode.")
	clientCheckInterval = 60 * time.Second

	authServer = "192.168.1.3"
)

func main() {
	flag.Parse()

	if *server && *client {
		log.Fatal("Only one of -server or -client toggles can be used!")
	}

	if *server {
		if *authFile == "" {
			log.Fatal("-file required when running in server mode.")
		}
		runServer()
	} else if *client {
		runClient()
	}

	log.Fatal("One of -server or -client toggles must be used!")
}

func runServer() {

	if _, err := os.Stat(*authFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal("could not find given authorized_keys file! (or it doesn't exist at given path)")
		}
		log.Fatal("error checking authorized_keys file")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if f, err := os.ReadFile(*authFile); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Printf("failed to read auth file! %v\n", err)
		} else {
			w.Write(f)
		}
	})
	log.Println("listening on port " + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func runClient() {
	targetFile := fmt.Sprintf("/home/%s/.ssh/authorized_keys", *user)
	ticker := time.NewTicker(clientCheckInterval)
	cl := http.DefaultClient
	log.Printf("will check for new file every %s", clientCheckInterval)
	getFile(cl, targetFile)

	for {
		select {
		case <-ticker.C:
			getFile(cl, targetFile)
		}
	}
}

func getFile(cl *http.Client, targetFile string) {
	resp, err := cl.Get(fmt.Sprintf("http://%s:%s", authServer, *port))
	if err != nil {
		log.Printf("got error from requesting new authorized_keys! %v\n", err)
	} else {
		f, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("got error from reading new authorized_keys! %v\n", err)
		} else {
			if err = os.WriteFile(targetFile, f, 0600); err != nil {
				log.Printf("failed to write new authorized keys file to %s: %v", targetFile, err)
			} else {
				os.Chown(targetFile, 1000, 1000)
				log.Println("successfully wrote new file")
			}
		}
	}
}
