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
A biggie that's missing with this code project is the connection to a real DB. User entries would have unique ID's. Quieres for messages in "last 30 day" would be significantly simpler.

# Testing
Another reason I picked Go was to start learning it's built in testing framework.
After getting stuck for a long time trying to simply wire up swagger, I ended up using Postman for most all of my manual testings. I did not get to writing automated tests (though I always highly prioritize having auto tests that strive for 100% code coverage)

# Review
After working on this project, I've learned a lot about Go. It was a totally new framework for me. Building it from the ground up was a blast! But... I probably should have chosen a langue and framework I already knew and expanded on topics I already knew/have used. Instead, most my time was spent learning the Go structure and basic concepts (working with array, var typs/struct, parsing request, go + postgres, ect). I also waisted a LONG time trying to get the self documenting swagger module. It helped me learn a lot about go's file/package structure and I feel I now have a solid grasp on go's module's (which seem to be a recently new concept for go)
I also have a lot of implementation missing for this service. It's very brittle and all the endpoints/feature's I'd like to expose are missing (such as requests for messages with time stamps). It would be more correct to have the service working with a database.
I also am frustrated with how much of the code and logic is structured. A lot of the coding feels hacky to me, and it is! It was a new langue. My variables names are a mess and very confusing to read. It feels like a rough POC, but a good starting point if this where to be "productized"
