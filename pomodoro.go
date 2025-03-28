package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"log"
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

// Updates the text in the bar by writing to ~/.local/share/pomodoro/output.txt
// and signaling waybar to read the change
func updateBar(str string) {
	programDir, err := getProgramDir()
	if err != nil {
		fmt.Println("Error getting program directory:", err)
		return
	}

	outputPath := filepath.Join(programDir, "output.txt")

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write output to file
	_, err = file.WriteString(str)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// Signal waybar to update
	cmd := exec.Command("pkill", "-SIGRTMIN+10", "waybar")
	if err := cmd.Run(); err != nil {
		log.Fatal("Error sending signal to waybar:", err)
	}
}

// Get the absolute path of ~/.local/share/pomodoro
func getProgramDir() (string, error) {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	programDir := filepath.Join(homeDir, ".local", "share", "pomodoro")

	// Create program directory if non existent
	err = os.MkdirAll(programDir, 0755)
	if err != nil {
		return "", err
	}

	return programDir, nil
}

// Send a notification with a message
func notifySend(msg string) {
	go playSound()

	cmd := exec.Command("notify-send", "Pomodoro", msg)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// Play embedded wav file
func playSound() {
	// Create reader from embedded wav data
	wavReader := bytes.NewReader(wavData)

	// Decode .wav file and get streamer
	streamer, format, err := wav.Decode(wavReader)
	if err != nil {
		log.Fatal(err)
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
		updateBar("{\"text\": \"󰓛 " + t + " " + strconv.Itoa(l-i) + " min\"}")
		time.Sleep(time.Minute)
	}
}

func main() {
	// Parse flags
	workLen := flag.Int("w", 25, "Work length in minutes")
	breakLen := flag.Int("b", 5, "Break length in minutes")
	flag.Parse()

	// Capture interupt/terminate signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// On received signal, update bar to reflect stop and exit
	go func() {
		sig := <-sigs
		fmt.Println("Received signal:", sig)
		updateBar("{\"text\": \"󰐊 Work " + strconv.Itoa(*workLen) + " min\"}")
		os.Exit(0)
	}()

	for {
		session(*workLen, "Work")
		notifySend("Time for a short break.")
		session(*breakLen, "Break")
		notifySend("Let's get back to work.")
	}
}
