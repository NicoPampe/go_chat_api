
# Getting started
To run locally:
* `go` needs to be installed
* `docker` should also be installed

All that should be needed to start the application: from root
```
./start
```
This will build and run a docker image with the go application.
The image is run in the terminal to easily stop the application (control + c).
In a normal production env, it would be run in a detached mode

# Start up our postgress Server (mac)
brew services start postgresql

# Using the go_chat_api
For this project, I've used postman to test the endpoints. "curl" requests can also be used to test endpoints.
Several endpoints for the messages API have been created.
For a demo, two users are initialized for examples.

1st User: {
  Id: "1"
}
2nd User: {
  Id: "2"
}

## home page endpoint
The first endpoint that can be hit is the demo homePage default
```
"/"
```

## "all" endpoint
An endpoint to get all users and their info. The responses is a json blob
```
"/all"
```
Example responses:
[
  {"Id":"1","first_name":"Nico","last_name":"Pampe","messages":{"2":"abc"}},
  {"Id":"2","first_name":"Bilbo","last_name":"Baggins","messages":{"1":"abc"}}
]

## "user" endpoint
And endpoint for CRUD operations on a "user" element.
A User has several elements:
  Id: a unique Id to look up and identify Users
  first_name: a users first name
  last_name: a users last name
  messages: a map of key/value pairs that map to user/messageKey. The key represents the second User in the chat. The value represents a unique Id in the messages tables. All messages are stored in a separate structure.

Create new user.
```
POST "/user"
BODY: {
    "Id": "3",
    "first_name": "Frodo",
    "last_name": "Baggins"
  }
```

Get a single users info.
```
POST "/user/{id}"
```
Example responses:
`http://localhost:10000/user/1`
{"Id":"1","first_name":"Nico","last_name":"Pampe","messages":{"2":"abc"}}

## PUT "message" endpoint
The main endpoint used to handle messaging.
A message has only two elements:
  Id: a unique Id that maps a user to the messages. Their should only ever be a 2:1 ratios (2 users : 1 message array)
  Messages: an array of string that contain all the messages sent between two users

Create a message
```
PUT "/message"
BODY: {
    "UserA": "1",
    "UserB": "2",
    "Message": "This is the way"
  }
```
For this project, since auth is not necessary, it's assumed that two users can freely send messages. I also decided to assume that anyone consuming the API would have all user Id's.
The body must have 3 elements:
  UserA: the ID of the first user
  UserB: the ID of the second user
  Message: the test of the message

If a user doesn't exist, an error will be returned:
`Error: User 3 does not exist`

## "message" endpoint
Get all messages from all users
```
"/message"
```
The response is a json blob. Example:
[
  {
    "Id":"abc",
    "messages":["Hello bilbo!","Hello Nico! Have you seen my precious","I'm afraid not"]
  }
]

## "/message/{id}" endpoint
Get all the messages sent/received for a user, provided a correct ID.
If the user ID does not exist, the response will be `null` (for security)
```
"/message/{ID}"
```
Example response:
["Hello bilbo!","Hello Nico! Have you seen my precious","I'm afraid not", "this is the way"]

## "/message/{id}/{message_id_req}" endpoint
Get all the messages sent/received for a user, provided a correct ID and Key in the Messages map.
If the user ID does not exist, the response will be `null` (for security)
```
"/message/{ID}/{message_id_req}"
```
Example: `"/message/1/abc"` - get the message of user "1" and the mapping "abc"
["Hello bilbo!","Hello Nico! Have you seen my precious","I'm afraid not"]


# Self documenting Docs
To generate swaggo's docs, run
```
swag init
```
and a directory "./docs" will be created. The file docs.go should run along side the rest of the go application. It is reliant on the additional go packages to be correctly installed.
