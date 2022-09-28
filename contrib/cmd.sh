#!/bin/sh

# Run this in the background to start the server
go run /app/main.go &

/sbin/my_init
