package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-audio/wav"
	"github.com/hajimehoshi/oto"
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
	data, err := os.ReadFile("./mixkit-correct-answer-tone-2870.wav")
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}

	dec := wav.NewDecoder(bytes.NewReader(data))
	buf, err := dec.FullPCMBuffer()
	if err != nil {
		log.Fatal(err)
	}

	byteBuffer := make([]byte, len(buf.Data)*2)
	for i, sample := range buf.Data {
		binary.LittleEndian.PutUint16(byteBuffer[i*2:], uint16(sample))
	}

	ctx, err := oto.NewContext(buf.Format.SampleRate, buf.Format.NumChannels, 2, 4096)
	if err != nil {
		log.Fatal(err)
	}

	player := ctx.NewPlayer()
	_, err = player.Write(byteBuffer)
	if err != nil {
		log.Fatal(err)
	}
	player.Close()
	ctx.Close()

	cmd := exec.Command("notify-send", "Pomodoro", msg)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// Start a session with length l and type t
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
