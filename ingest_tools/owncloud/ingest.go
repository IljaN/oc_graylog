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

// toGELFTime converts owncloud log-timestamp to GELF time format
func toGELFTime(ocTime string) (float64, error) {
	var t time.Time
	var err error

	if t, err = time.Parse("2006-01-02 15:04:05.999999", ocTime); err != nil {
		return 0, fmt.Errorf("error parsing oc time: %v", err)
	}

	return float64(t.UnixNano()) / float64(time.Second), err
}

// convertToGELF converts a owncloud JSON log line to GELF format
func convertToGELF(logLine *[]byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(*logLine, &data); err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	ocTimestamp, ok := data["time"].(string)
	if !ok {
		return fmt.Errorf("error: 'time' field not found or not a string")
	}

	gelfTime, err := toGELFTime(ocTimestamp)
	if err != nil {
		return fmt.Errorf("error converting time to GELF format: %v", err)
	}

	// Required fields by GELF
	data["timestamp"] = gelfTime
	data["version"] = "1.1"
	data["host"] = "localhost"

	// Convert back to JSON
	if *logLine, err = json.Marshal(data); err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	return nil
}

func main() {

	var gelfUDPAddr string

	gelfUDPAddr = os.Getenv("OC_INGEST_GELF_UDP_ADDR")
	if gelfUDPAddr == "" {
		flag.StringVar(&gelfUDPAddr, "a", "127.0.0.1:12201", "GELF UDP address")
		flag.Parse()
	}

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
		line := scanner.Bytes()
		if err = convertToGELF(&line); err != nil {
			log.Printf("Error processing line: %v", err)
			continue
		}

		if _, err = conn.Write(line); err != nil {
			log.Printf("error sending data: %v", err)
		}

		msgCount++
	}

	elapsed := time.Since(startT)
	fmt.Printf("Processed %d messages in %f secs\n", msgCount, elapsed.Seconds())
}
