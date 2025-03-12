package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	// Define flags with default values and descriptions
	workLen := flag.Int("w", 25, "Work length in minutes")
	breakLen := flag.Int("b", 5, "Break length in minutes")
	// longBreakLen := flag.Int("l", 15, "Long break length in minutes")
	// sessions := flag.Int("s", 4, "Number of sessions")

	// Parse the command-line arguments
	flag.Parse()

	// Print the values
	workMsg := "Work for " + strconv.Itoa(*workLen) + " minutes"
	breakMsg := "Break for " + strconv.Itoa(*breakLen) + " minutes"

	fmt.Println("work len:", *workLen, "m")
	fmt.Println("break len:", *breakLen, "m")
	// fmt.Println("long break len:", *longBreakLen, "m")
	// fmt.Println("sessions:", *sessions, "m")

	cmd := exec.Command("notify-send", "Pomodoro", workMsg)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(time.Duration(*workLen) * time.Minute)

		cmd = exec.Command("notify-send", "Pomodoro", "Work session ended\n"+breakMsg)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Duration(*breakLen) * time.Minute)

		cmd := exec.Command("notify-send", "Pomodoro", "Break is over\n"+workMsg)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

}
