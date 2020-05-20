package main

import (
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/app"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("A port must provided on the command line")
		os.Exit(1)
	}

	port := os.Args[1]
	portInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to parse `%v` as a port number", port))
		os.Exit(1)
	}
	hashService := app.NewHashService(portInt)
	hashService.Start()
}