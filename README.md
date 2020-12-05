# go_chat_api

A messenger API built in Golang.
The API allows users to connect to a server where


# Why Golang
I wanted to practice using a framework that I don't often get to use.
A few months ago, a new project I started helping with used Go for the backend. I wanted to work with Go from the ground up.

I've also done this same project (messenger API) multiple times in node for Guild and other companies in the past. I feel that trying a new langue will be a good practice and show I'm excited to dive into new technologies. I also hope gain some new skills! Now let's Go!

# Documentation
https://github.com/swaggo/gin-swagger

# What I Found
Go is super cool! But wow... there's a lot to learn.
I am struggling with `go.mod` and go's package structure. As an idea, I understand the structure and want to dig in farther. But I've been stuck for several hours and need to change gears.


# What I started but didn't get to
Hooking up Go to a postgres db. I started wiring up the connection, successfully create a user DB, but struggled to insert data from a POST "createNewUser". I decided to keep the data in a local struct in the go application.
This would also be a more "correct" way to append messages to a data structure.
All of the meta data that would come with useing postgres could be used for the messaging API. You could get more accurate time stamps, order or insert, and much more!

# Testing
Another reason I picked Go was to start learning it's built in testing framework.
After getting stuck for a long time trying to simply wire up swagger, I ended up using Postman for most all of my manual testings. I did not get to writing automated tests (though I always highly prioritize having auto tests that strive for 100% code coverage)
