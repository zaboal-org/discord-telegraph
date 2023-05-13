package main

import (
	"flag"
	"log"
)

// Token — an authorization token to Discord API for bots
var (
	Token string
)

// Retrieve the passed arguments on run
func init() {
	flag.StringVar(&Token, "token", "", "Authorization token to Discord API for bots")
	flag.Parse()
}

func main() {
	log.Println("Got the token:", Token)
}
