package main

import (
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/app/hash"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
	hashService := hash.New(portInt)

	ctrlC := make(chan os.Signal)
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctrlC
		fmt.Println("Interrupt caught, requesting shutdown")
		hashService.Stop()
	}()
	hashService.Start()
}