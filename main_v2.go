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
  fmt.Println(r)


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

func message(w http.ResponseWriter, r *http.Request) {
  fmt.Println("message...")
  r.ParseForm()
  // fmt.Println(r)
  // fmt.Println(r.ParseForm())
  fmt.Printf("MESSAGE => %s\n", r.FormValue("message"))
  // This is working, but I honestly don't understand how the `json.Unmarshal` works.
  // TODO: learn excactly what json.Unmarshal does/mean.


  reqBody, _ := ioutil.ReadAll(r.Body)
  type MessageReq struct {
    UserA string `json:"UserA"`
    UserB string `json:"UserB"`
    Text string `json:"message"`
  }
  var message MessageReq
  json.Unmarshal(reqBody, &message)
  json.NewEncoder(w).Encode(message)


  // foo := reqBody["id_a"]
  fmt.Println(message)

  vars := mux.Vars(r)
  fmt.Println(vars)
  userA := vars["id_a"]
  userB := vars["id_b"]
  // message := vars["message"]

  fmt.Printf("userA %#v\n", userA)
  fmt.Printf("userB %#v\n", userB)
  fmt.Printf("message %#v\n", message)
  // fmt.Println()
  // fmt.Println("userB", userB)
  // fmt.Println("message", message)
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
  myRouter.HandleFunc("/message/{id}", getAllMessage)

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
    User{Id: "1", FirstName: "Nico", LastName: "Pampe"},
    User{Id: "2", FirstName: "Bilbo", LastName: "Baggins"},
  }
  Messages = []Message{
    Message{Id: "1", Messages: []string{"Hello bilbo!", "Hello Nico! Have you seen my precious", "I'm afraid not"} },
  }
  handleDB()
  handleRequests()
}
