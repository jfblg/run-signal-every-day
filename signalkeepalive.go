package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func checkConnectivity() error {
	_, err := http.Get("https://opennic.org/") // use a public website
	return err
}

func waitForConnectivity(timeout time.Duration) error {
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(1 * time.Second) {
		if err := checkConnectivity(); err != nil {
			log.Printf("connectivity check failed: %v", err)
			continue
		}
		return nil // connectivity check succeeded
	}
	return fmt.Errorf("no connectivity established within %v", timeout)
}

func keepalive() error {
	// Limit to one attempt per day by exclusively creating a logfile.
	home := os.Getenv("HOME")
	if home == "" {
		home = "/Users/frantisekjanus"
	}

	logDir := filepath.Join(home, "signal-keepalive", "_logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}
	logFn := filepath.Join(logDir, time.Now().Format("2006-01-02")+".txt")
	f, err := os.OpenFile(logFn, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil // nothing to do, already ran today
		}
		return err
	}
	// Intentionally not closing this file so that even the log.Fatal()
	// call in the calling function will end up in the log file.

	log.SetOutput(f) // redirect standard library logging into this file
	log.Printf("signal-keepalive, waiting for internet connectivity")

	// Wait for network connectivity
	if err := waitForConnectivity(10 * time.Minute); err != nil {
		return err
	}

	// Start signal
	log.Printf("connectivity verified, starting signal")
	signal := exec.Command("/Applications/Signal.app/Contents/MacOS/Signal", "--start-in-tray")
	signal.Stdout = f
	signal.Stderr = f
	if err := signal.Start(); err != nil {
		return err
	}

	// Wait for some time to give Signal a chance to synchronize messages.
	const signalWaitTime = 5 * time.Minute
	// const signalWaitTime = 30 * time.Second
	log.Printf("giving signal %v to sync messages", signalWaitTime)
	time.Sleep(signalWaitTime)

	// Stop signal
	log.Printf("killing signal")
	if err := signal.Process.Kill(); err != nil {
		return err
	}
	log.Printf("waiting for signal")
	log.Printf("signal returned: %v", signal.Wait())
	log.Printf("all done")

	return f.Sync()
}

func main() {
	if err := keepalive(); err != nil {
		log.Fatal(err)
	}
}
