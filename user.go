package user

// package model
//
// type Account struct {
//     ID   int    `json:"id" example:"1"`
//     Name string `json:"name" example:"account name"`
// }

type User struct {
    ID   int    `json:"id" example:"1"`
    Name string `json:"name" example:"account name"`
    Title string `json:"Title"`
    Desc string `json:"desc"`
    Content string `json:"content"`
}

var Users []User
