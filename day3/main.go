package main

import (
	"day3/handlers"
	"day3/internal/store"
	"fmt"
	"log"
	"net/http"
)

// main for task 1 and 2
/*
func main() {
//TODO: If wan't to run then create a function to generate random string
 s := store.NewUserStore()

 go func() {
	for {
	id := randString()
	if err := s.Create(store.User{ID: id, Name: "Hanu", Email: "harbinder12@gmail.com"}); err != nil {
			log.Println(err)
		}
	}
}()

 go func ()  {
	for {
	id := randString()
	u, err := s.GetByID(id)
	if err != nil {
		log.Println(err)
 	}else {
		fmt.Printf("%+v", u)
	}
	}
}()
//  data := s.List()
//  fmt.Printf("%+v", data)

  // wait for signal to close
 // signal channels
 sigs := make(chan os.Signal, 1)
 signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

 // Waiting for shutdown sig
 <-sigs
}
*/

func main() {
	// not a good approach to expose Users map
	// st := &store.UserStore{
	// 	Users: make(map[string]store.User),
	// }

	st := store.NewUserStore() // constructor pattern better approach

	userHandler := handlers.NewUserHandler(st)
	//handlers.InitStore(st)

	fmt.Println("Server is running on http://localhost:8080")

	log.Println("Request received")
	// ?? Find Later:  what is the diffrence with creating routes with mutex
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet :
			userHandler.GetAllUsers(w, r)
		case http.MethodPost :
			userHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}	
	})
	
	http.HandleFunc("/users/{id}", userHandler.GetUserbyId)

	// Mistake: same router don't work, pretty obvious though
	// http.HandleFunc("/users", handlers.CreateUser)
	// http.HandleFunc("/users", handlers.GetAllUsers)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Unalble to run server")
	}
}