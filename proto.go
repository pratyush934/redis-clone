package main

import "fmt"

type Command struct {
}

func parseMessage(msg []byte) (Command, error) {
	fmt.Print(string(msg))
	return nil, nil
}
