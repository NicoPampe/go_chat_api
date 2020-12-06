#!/bin/bash

# This is a very very simple bash script with next to no error handlings.
# TODO: make the deployment script more robust.
docker build . -t go_chat_api
docker run -p 10000:10000 go_chat_api
