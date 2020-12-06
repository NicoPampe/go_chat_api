package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "hash/fnv"
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
  message_id_req := vars["message_id_req"]

  var userMessages []string

  for _, user := range Users {
    if user.Id == key {
      for _, messageId := range user.Messages {
        // If we have a message id to look for, only get that entry.
        // Otherwise, append all messages
        var message []string
        if message_id_req != "" {
          // TODO: this should be handled in the getMessages func
          if message_id_req == messageId {
            message = getMessages(messageId)
          }
        } else {
          message = getMessages(messageId)
        }
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
  fmt.Println("Endpoint Hit: message")
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
  for usersIdx, user := range Users {
    if user.Id == messageReq.UserA {
      // Next check if UserB is in UserA's message history
      for recipientId, messageId := range user.Messages {
        if recipientId == messageReq.UserB {
          // simply append the new message to the list
          appendMessage(messageId, messageReq.Message)
          return // break out of the function (not the best pratice)
        }
      }

      // If UserB was not in UserA's list, create a new entry
      createMessageBetweenTwoUsers(w, Users[usersIdx], messageReq.UserB, messageReq.Message)
    }
  }
}

// This function creates a new entry in the Messages structure given two user.
// @param userA User     - A refrence to a user in the User db
// @param userBId string - A string to find the refrence of UserB
// @param m string       - the message
func createMessageBetweenTwoUsers(w http.ResponseWriter, userA User, userBId string, m string)  {
  fmt.Println("createMessageBetweenTwoUsers")

  // check if UserB exsits
  var (
    userB User // this should be the correct pratice. I'm stuggling getting the slice ref to work. learned that Go passes by value and it can be tricky getting the ref
    foundUserB bool = false // This is a teriable pattern. Should be able to just check "userB"
    targetUserIndex int // I should be trying to use ref. Failing to do so.
  )

  for usersIdx, user := range Users {
    if user.Id == userBId {
      userB = Users[usersIdx] // get the pointer
      targetUserIndex = usersIdx // ewwwwww I'm sad that this pattern is what I ended up with. Given more time, I would work to get a better understanding of Go's pointers/refrences. Array's/slice's in go are very intresting
      foundUserB = true
      break
    }
  }
  if foundUserB == false {
      http.Error(w, "Error: User " + userBId + " does not exist", http.StatusBadRequest)
      return
  }

  // succes!! let's create that mapping
  // my quick soluntion for an ID is to sha1 the two user ID's. There are more correct ways to do this using a real DB.
  key := genTwoUsersHash(userA.Id, userB.Id)

  if userA.Messages == nil {
    userA.Messages = make(map[string]string)
  }
  if userB.Messages == nil {
    userB.Messages = make(map[string]string)
  }
  userA.Messages[userB.Id] = key
  userB.Messages[userA.Id] = key

  var message Message
  message.Id = key
  message.Messages = []string{m}
  fmt.Println(message)
  // add our newly created entry to the struct (or DB)
  Messages = append(Messages, message)

  // Finally, update the original userB. This code shouldn't be used, Instead the value should be updated by refrence, not by value.
  // I couldn't figure out how to get the value by refrence in the local block.
  Users[targetUserIndex] = userB
}

func genTwoUsersHash(s1 string, s2 string) string {
  h := fnv.New32a()
  h.Write([]byte(s1 + s2 ))
  return fmt.Sprint(h.Sum32())
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
  myRouter.HandleFunc("/user", returnAllUsers)
  myRouter.HandleFunc("/user/{id}", returnSingleUser)

  // sending messages
  myRouter.HandleFunc("/message", message).Methods("PUT")
  myRouter.HandleFunc("/message", getAllMessage)
  myRouter.HandleFunc("/message/{id}", getUserMessages)
  myRouter.HandleFunc("/message/{id}/{message_id_req}", getUserMessages)

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
  // TODO: this should be uncommented to start and wire up to a local postgress DB. Currently doesn't work in a docker containter
  // handleDB()
  handleRequests()
}
