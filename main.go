package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

// struct which contains filepath
type ConfigHandler struct {
	CounterFilePath string
}

func readCountFromFile(filename string) int {
	slog.Info("Reading count from file", "filename", filename)
	var count int

	// open c.CounterFilePath and read count
	file, err := os.Open(filename)

	if errors.Is(err, os.ErrNotExist) {
		slog.Warn("File does not exists, returning 0", "filename", filename)
	} else {
		check(err)
		_, err = fmt.Fscanf(file, "%d", &count)
		check(err)
		slog.Info("Read count from file", "count", count)
	}

	file.Close()

	return count
}

func writerCountToFile(filename string, count int) {
	slog.Info("Writing count to file", "filename", filename, "count", count)

	// open c.CounterFilePath and write count
	file, err := os.Create(filename)
	check(err)
	_, err = fmt.Fprintf(file, "%d", count)
	check(err)

	file.Close()
	slog.Info("Wrote count to file", "count", count)
}

func (c ConfigHandler) handleCount(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received request", "method", r.Method, "path", r.URL.Path)

	count := readCountFromFile(c.CounterFilePath)

	if (r.Method == "POST") || (r.Method == "PUT") {
		// increment count
		count++

		slog.Info("Incrementing count", "newCount", count)

		writerCountToFile(c.CounterFilePath, count)
	}

	// write count to response
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Count: %d", count)
	check(err)

	slog.Info("Response sent", "count", count)
}

func check(e error) {
	if e != nil {
		slog.Error("Error occurred", "error", e)
		panic(e)
	}
}

func initConfig() ConfigHandler {
	// read filepath from env variables
	counterFilePath := os.Getenv("COUNTER_FILE_PATH")
	if counterFilePath == "" {
		counterFilePath = "counter.txt"
	}

	slog.Info("Config using counter file path", "filePath", counterFilePath)

	var config ConfigHandler
	config.CounterFilePath = counterFilePath

	return config
}

func main() {
	// Initialize the logger
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	myLogger := slog.New(jsonHandler)
	slog.SetDefault(myLogger)

	configHandler := initConfig()

	var port = os.Getenv("COUNTER_PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/count", configHandler.handleCount)
	slog.Info("Server started on " + port)
	res := http.ListenAndServe("localhost:"+port, nil)

	if res != nil {
		slog.Error("Error starting server", "error", res.Error())
		return
	}
}
