package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

// Embed the notification sound file

//go:embed mixkit-correct-answer-tone-2870.wav
var wavData []byte

// Write output to ~/.local/share/pomodoro/output.txt
func writeOutput(output string) {
	dirPath, err := getOutputDir()
	if err != nil {
		fmt.Println("Error getting program directory:", err)
		return
	}

	filePath := filepath.Join(dirPath, "output.txt")

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

// Signal waybar to update
func updateWaybar() {
	cmd := exec.Command("pkill", "-SIGRTMIN+10", "waybar")
	if err := cmd.Run(); err != nil {
		fmt.Println("Error sending signal to waybar:", err)
		return
	}
}

// Get the absolute path of ~/.local/share/pomodoro
func getOutputDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	outDir := filepath.Join(homeDir, ".local", "share", "pomodoro")

	// Create program directory if non existent
	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return "", err
	}

	return outDir, nil
}

// Send a notification with a message
func notifySend(msg string) {
	go playSound()

	cmd := exec.Command("notify-send", "Pomodoro", msg)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error sending notification:", err)
		return
	}
}

// Play embedded wav file
func playSound() {
	// Create reader from embedded wav data
	wavReader := bytes.NewReader(wavData)

	// Decode .wav file and get streamer
	streamer, format, err := wav.Decode(wavReader)
	if err != nil {
		fmt.Println("Error decoding wav file:", err)
		return
	}
	defer streamer.Close()

	// NOTE: Should only be done once?
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Play sound and block until done
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

// Start a session with length l minutes and type t
// t can be any string, it only affects what is printed
func session(l int, t string) {
	for i := range l {
		writeOutput("{\"text\": \"󰓛 " + t + " " + strconv.Itoa(l-i) + " min\"}")
		updateWaybar()
		time.Sleep(time.Minute)
	}
}

// parseFlags defines and parses the command-line flags.
func parseFlags() (workLen *int, breakLen *int) {
	workLen = flag.Int("w", 25, "Work length in minutes")
	breakLen = flag.Int("b", 5, "Break length in minutes")
	flag.Parse()
	return workLen, breakLen
}

// cleanExit performs a clean exit updating Waybar to reflect the stop and exit
func cleanExit(workLen int) {
	fmt.Println("Exiting cleanly")
	writeOutput("{\"text\": \"󰐊 Work " + strconv.Itoa(workLen) + " min\"}")
	updateWaybar()
	os.Exit(0)
}

// monitorInterrupt listens for an interrupt signal and exits cleanly if received.
func monitorInterrupt(workLen int) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals
	fmt.Println("Received signal:", sig)
	cleanExit(workLen)
}

func main() {
	workLen, breakLen := parseFlags()

	go monitorInterrupt(*workLen)

	for {
		session(*workLen, "Work")
		notifySend("Time for a short break.")
		session(*breakLen, "Break")
		notifySend("Let's get back to work.")
	}
}
