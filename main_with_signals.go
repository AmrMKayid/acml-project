package main

import (
  	"fmt"
  	"os"
    "os/signal"
    "syscall"
	"time" 
)

var (
	canTerminate = false
)

func waiting() {
    fmt.Println("\nWaiting for loop to finish\n")
}

func goodbye() {
	fmt.Println("\nGoodbye!\n")
}

func check() {
	c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        if(canTerminate) {
        	goodbye()
        	os.Exit(1)
        } else {
        	waiting()
        	canTerminate = true
        }
    }()
}

func main() {
	for i := 1; i <= 10; i++ {
		check()
		fmt.Println(i)
		time.Sleep(time.Second) 
	}
}