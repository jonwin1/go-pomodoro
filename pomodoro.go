package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

// Uppdates the text in the bar by writing to file
// and signaling waybar to read the change
func updateBar(str string) {
	// Create/truncate file
	file, err := os.Create("/home/jonwin/go-pomodoro/log")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close() // Close when done

	// Write meassage to file
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

// Send a notification with a message
func notifySend(msg string) {
	go playSound("./mixkit-correct-answer-tone-2870.wav")

	cmd := exec.Command("notify-send", "Pomodoro", msg)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// Play a .wav file given it's path
func playSound(path string) {
    // Open the file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

    // Decode .wav file and get streamer
	streamer, format, err := wav.Decode(file)
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
