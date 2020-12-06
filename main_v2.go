package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/gorilla/mux"
    "database/sql"
    _ "github.com/lib/pq"
)

const (
  host     = "localhost"
  port     = 5432
  user     = "nicholas.pampe"
  password = "your-password"
  dbname   = "postgres"
)

type User struct {
    Id string `json:"Id"`
    FirstName string `json:"first_name"`
    LastName string `json:"last_name"`
    Messages map[string]string `json:"messages"` // A collection that maps a users messages to the message Table [key:userId]MessagesId
}
var Users []User

// A single entry in the Messages table contains all the messages between two users
type Message struct {
    Id string `json:"Id"`
    Messages []string `json:"messages"`
}
var Messages []Message


func homePage(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome to the HomePage!")
  fmt.Println("Endpoint Hit: homePage")
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Endpoint Hit: returnAllUsers")
  json.NewEncoder(w).Encode(Users)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := vars["id"]

  for _, user := range Users {
    if user.Id == key {
      json.NewEncoder(w).Encode(user)
    }
  }
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
  fmt.Println("createNewUser")

  reqBody, _ := ioutil.ReadAll(r.Body)

  var user User
  json.Unmarshal(reqBody, &user)
  // update our global Users array to include our new User
  Users = append(Users, user)

  json.NewEncoder(w).Encode(user)

  // Should be working... It's not ;(
  /**
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  sqlStatement := `
    INSERT INTO users (first_name, last_name)
    VALUES ($1, $2)`
  _, err = db.Exec(sqlStatement, "Nico", "Pampe")
  if err != nil {
    panic(err)
  }
  */
}

func getAllMessage(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Endpoint Hit: getAllMessage")
  json.NewEncoder(w).Encode(Messages)

}

func getUserMessages(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Endpoint Hit: getUserMessages")
  vars := mux.Vars(r)
  key := vars["id"]

  var userMessages []string

  for _, user := range Users {
    if user.Id == key {
      for _, messageId := range user.Messages {
        var message = getMessages(messageId)
        userMessages = append(userMessages, message...)
      }
    }
  }

  json.NewEncoder(w).Encode(userMessages)
}

func getMessages(messageId string) []string {
  for _, message := range Messages {
    if message.Id == messageId {
      return message.Messages
    }
  }
  return nil
}

// TODO: this function is very sloopy. There should be a cleaner soluntion to finding the "messages" between two users.
func message(w http.ResponseWriter, r *http.Request) {
  type MessageReq struct {
    UserA string
    UserB string
    Message string
  }
  var messageReq MessageReq

  // parse the req into a struct
  err := json.NewDecoder(r.Body).Decode(&messageReq)
  if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  }
  // fmt.Fprintf(w, "Person: %+v", message)

  fmt.Println(messageReq)
  fmt.Println(messageReq.UserA)


  // first, find UserA
  for _, user := range Users {
    if user.Id == messageReq.UserA {
      // Next check if UserB is in UserA's message history
      for recipientId, messageId := range user.Messages {
        fmt.Println(recipientId)
        fmt.Println(messageReq.UserB)
        if recipientId == messageReq.UserB {
          fmt.Println("they are equal!")
          // simply append the new message to the list
          appendMessage(messageId, messageReq.Message)
          break
        }
      }
    }
  }
}

// TODO: this should handle errors
func appendMessage(key string, m string)  {
  for idx, message := range Messages {
    if message.Id == key {
      Messages[idx].Messages = append(message.Messages, m)
    }
  }
}


func handleRequests() {
  // creates a new instance of a mux router
  myRouter := mux.NewRouter().StrictSlash(true)
  myRouter.HandleFunc("/", homePage)
  myRouter.HandleFunc("/all", returnAllUsers)
  myRouter.HandleFunc("/user", createNewUser).Methods("POST")
  myRouter.HandleFunc("/user/{id}", returnSingleUser)

  // sending messages
  myRouter.HandleFunc("/message", message).Methods("PUT")
  myRouter.HandleFunc("/message", getAllMessage)
  myRouter.HandleFunc("/message/all", getAllMessage)
  myRouter.HandleFunc("/message/{id}", getUserMessages)

  log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func handleDB()  {
  // Open up our database connection.
  // I've set up a postgress database on my local machine using `brew services start postgresql`.
  // The database is called postgres
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)


  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }

  fmt.Println("Successfully connected!")
  fmt.Println("Init Table (Normally, there would be a persistent table. For this demo, a fresh user table will be created)")

  sqlStatement := `
    CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, first_name TEXT, last_name TEXT, messagesMap JSONB)`
  _, err = db.Exec(sqlStatement)
  if err != nil {
    panic(err)
  }
}

func main() {
  fmt.Println("Rest API v2.0 - Mux Routers")
  Users = []User{
    User{Id: "1", FirstName: "Nico", LastName: "Pampe", Messages: map[string]string{"2":"abc"}},
    User{Id: "2", FirstName: "Bilbo", LastName: "Baggins", Messages: map[string]string{"1":"abc"}},
  }
  Messages = []Message{
    Message{Id: "abc", Messages: []string{"Hello bilbo!", "Hello Nico! Have you seen my precious", "I'm afraid not"} },
  }
  handleDB()
  handleRequests()
}
