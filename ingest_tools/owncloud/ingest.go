package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func processLine(line string, conn net.Conn) error {
	// Parse the JSON string
	var data map[string]interface{}
	err := json.Unmarshal([]byte(line), &data)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)

	}

	// Get the value of the "time" field
	timeStr, ok := data["time"].(string)
	if !ok {
		return fmt.Errorf("error: 'time' field not found or not a string")
	}

	t, err := time.Parse("2006-01-02 15:04:05.999999", timeStr)
	if err != nil {
		return fmt.Errorf("error parsing time: %v", err)
	}

	data["timestamp"] = t.Unix()
	data["version"] = "1.1"
	data["host"] = "localhost"

	// Convert back to JSON
	modifiedLine, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Send modified line to the server
	_, err = conn.Write(modifiedLine)
	if err != nil {
		return fmt.Errorf("error sending data: %v", err)
	}

	return nil
}

func main() {

	var gelfUDPAddr string
	flag.StringVar(&gelfUDPAddr, "a", "127.0.0.1:12201", "GELF UDP address")
	flag.Parse()

	// Connect to server
	conn, err := net.Dial("udp", gelfUDPAddr)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	startT := time.Now()
	var msgCount uint64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		if err = processLine(line, conn); err != nil {
			log.Printf("Error processing line: %v", err)
			continue
		}
		msgCount++
	}

	elapsed := time.Since(startT)
	fmt.Printf("Processed %d messages in %f secs\n", msgCount, elapsed.Seconds())
}
