package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define flags with default values and descriptions
	workLen := flag.Int("w", 25, "Work length in minutes")
	breakLen := flag.Int("b", 5, "Break length in minutes")
	longBreakLen := flag.Int("l", 15, "Long break length in minutes")
	sessions := flag.Int("s", 4, "Number of sessions")

	// Parse the command-line arguments
	flag.Parse()

	// Print the values
	fmt.Println("work len:", *workLen)
	fmt.Println("break len:", *breakLen)
	fmt.Println("long break len:", *longBreakLen)
	fmt.Println("sessions:", *sessions)
}
