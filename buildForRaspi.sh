#!/bin/bash
env GOOS=linux GOARCH=arm GOARM=5 go build cocktail.go
