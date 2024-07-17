package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/JanSzymanski/todostorelib"
)

func errorWritter(w http.ResponseWriter, err error, err_msg string) {
	if err != nil {
		fmt.Fprint(w, err_msg)
	}
}

func main() {
	fmt.Println("Building REST API server for the ToDo application")
	todostore := todostorelib.NewTodoStore("Jan's Todo vault")
	todostore.AddTodo("Some new temporary todo")
	todostore.AddTodo("Yet another temporary todo")

	mux := http.NewServeMux()

	//GET
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from todo")
	})
	mux.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		todos := todostore.GetTodos(0, 20)
		b, err := json.MarshalIndent(todos, "", "    ")
		errorWritter(w, err, "Error in marshaling list of todos")
		w.Write([]byte(b))
		fmt.Println("Remote user requested a list of Todos")
	})
	mux.HandleFunc("GET /todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id_str := r.PathValue("id")
		id, err := strconv.Atoi(id_str)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Error in reading 'id': %q. Please provide a valid id number", id_str)
			return
		}
		fmt.Printf("Remote user requested details on todo id: %d\n", id)
		r.Header.Set("Content-Type", "application/json")
		todo, err := todostore.GetTodo(id)
		if err != nil {
			w.WriteHeader(204)
			fmt.Fprintf(w, "No todo with id: %d in the store.", id)
			return
		}
		b, err := json.MarshalIndent(todo, "", "    ")
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Error with marshaling todo of is: %d", id)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(b))
	})
	//POST
	mux.HandleFunc("POST /addtodo", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Cannot read the request body")
		}
		fmt.Printf("Creating user with message: %q\n", body)
		todostore.AddTodo(string(body))
		w.WriteHeader(201)
		fmt.Fprint(w, "Todo created")
	})
	//PATCH
	mux.HandleFunc("PATCH /chmtodo", func(w http.ResponseWriter, r *http.Request) {
		bodyDecoder := json.NewDecoder(r.Body)
		body := struct {
			Id  int    `json:"id"`
			Msg string `json:"msg"`
		}{}
		err := bodyDecoder.Decode(&body)
		if err != nil {
			fmt.Println("Cannot decode request body")
			fmt.Fprintf(w, "Todo: %q not edited.", body.Id)
			return
		}
		fmt.Println(body)
		todostore.ChangeTodoMessagge(body.Id, body.Msg)
		fmt.Fprintf(w, "Todo: %d message edited successfully.", body.Id)
	})
	mux.HandleFunc("PATCH /chstodo", func(w http.ResponseWriter, r *http.Request) {
		bodyDecoder := json.NewDecoder(r.Body)
		body := struct {
			Id   int    `json:"id"`
			Stat string `json:"stat"`
		}{}
		err := bodyDecoder.Decode(&body)
		if err != nil {
			fmt.Println("Cannot decode request body")
			fmt.Fprintf(w, "Todo: %q not edited.", body.Id)
			return
		}
		fmt.Println(body)
		todostore.ChangeTodoStatus(body.Id, todostorelib.Todostatus(body.Stat))
		fmt.Fprintf(w, "Todo: %d status changed to: %q.", body.Id, body.Stat)
	})
	//DELETE
	mux.HandleFunc("DELETE /deltodo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id_str := r.PathValue("id")
		id, err := strconv.Atoi(id_str)
		fmt.Println("Remote user requested to delete todo with id: ", id_str)
		if err != nil {
			fmt.Println("And that request ends with error.")
			fmt.Fprintf(w, "Cannot convert %q to int. Please provide a valid id.", id_str)
			return
		}
		err = todostore.DeleteTodo(id)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Fprint(w, err.Error())
		}
		fmt.Printf("Todo id: %d has been deleted.", id)
		fmt.Fprintf(w, "Todo id: %d has been deleted.", id)
	})

	serv := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	if err := serv.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
