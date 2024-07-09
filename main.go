package main

import (
	"log"
	"os"
)

func main() {
  log.SetFlags(log.Lshortfile)
	portString := os.Getenv("PORT")
  if portString == "" {
    portString = "8080"
    log.Println("Using Default Port: ",portString)
  } else {
    log.Println("Using Port: ",portString)
  }
}
