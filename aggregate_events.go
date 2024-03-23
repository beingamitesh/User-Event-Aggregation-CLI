package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"
)

type UserEvent struct {
	UserID    int    `json:"userId"`
	EventType string `json:"eventType"`
	Timestamp int64  `json:"timestamp"`
}

type DailySummary map[string]interface{}

type DailySummaries []DailySummary

var UserIdTimestampToEventTypeMap map[int]map[string]map[string]int

func main() {
	inputFile := flag.String("i", "", "Input JSON file path")
	outputFile := flag.String("o", "", "Output JSON file path")
	updateFlag := flag.Bool("update", false, "Update existing summary with new events")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Usage: ./main -i input.json -o output.json [--update]")
		os.Exit(1)
	}

	events, err := readEventsFromFile(*inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}

	summaries, err := aggregateEvents(events, *outputFile, *updateFlag)
	if err != nil {
		fmt.Println("Error aggregating events:", err)
		os.Exit(1)
	}

	err = writeSummariesToFile(summaries, *outputFile)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		os.Exit(1)
	}
}

func readEventsFromFile(filePath string) ([]UserEvent, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var events []UserEvent
	err = json.Unmarshal(file, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func readEventHashes() (map[string]interface{}, error) {
	file, err := os.ReadFile("hash.json")
	if err != nil {
		return nil, err
	}

	var eventHashes map[string]interface{}
	err = json.Unmarshal(file, &eventHashes)
	if err != nil {
		return nil, err
	}
	return eventHashes, nil
}

func checkIfNewEventHashExists(eventHashes map[string]interface{}, event UserEvent) (bool, error) {
	out, err := json.Marshal(event)
	if err != nil {
		return false, err
	}
	keyToCheck := fmt.Sprintf("%x", sha256.Sum256([]byte(string(out))))

	if _, exists := eventHashes[keyToCheck]; exists {
		return true, nil
	} else {
		return false, nil
	}
}

func aggregateEvents(events []UserEvent, outputFile string, updateFlag bool) (DailySummaries, error) {
	summaries := make(DailySummaries, 0)
	var err error
	UserIdTimestampToEventTypeMap := make(map[int]map[string]map[string]int)
	hashMap := make(map[string]interface{})
	if updateFlag {
		summaries, err = readSummariesFromFile(outputFile)
		if err != nil {
			return nil, err
		}
		hashMap, err = readEventHashes()
		if err != nil {
			return nil, err
		}
		for _, event := range events {
			eventHashExists, err := checkIfNewEventHashExists(hashMap, event)
			if err != nil {
				return nil, err
			}
			if !eventHashExists {
				updateSummary(&summaries, event, UserIdTimestampToEventTypeMap)
				err = hashFunc(event, hashMap)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		for _, event := range events {
			updateSummary(&summaries, event, UserIdTimestampToEventTypeMap)
			err := hashFunc(event, hashMap)
			if err != nil {
				return nil, err
			}
		}
	}
	err = writeEventHashToFile(hashMap)
	if err != nil {
		return nil, err
	}
	return summaries, nil
}

func readSummariesFromFile(filePath string) (DailySummaries, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var summaries DailySummaries
	err = json.Unmarshal(file, &summaries)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

func updateSummary(summaries *DailySummaries, event UserEvent, UserIdTimestampToEventTypeMap map[int]map[string]map[string]int) {
	date := time.Unix(event.Timestamp, 0).Format("2006-01-02")
	summary := getOrCreateSummary(summaries, event.UserID, date)
	if _, exists := UserIdTimestampToEventTypeMap[event.UserID]; !exists {
		UserIdTimestampToEventTypeMap[event.UserID] = make(map[string]map[string]int)
	}
	if _, exists := UserIdTimestampToEventTypeMap[event.UserID][date]; !exists {
		UserIdTimestampToEventTypeMap[event.UserID][date] = make(map[string]int)
	}
	UserIdTimestampToEventTypeMap[event.UserID][date][event.EventType]++
	updateSummaryFields(summary, event, UserIdTimestampToEventTypeMap)
}

func getOrCreateSummary(summaries *DailySummaries, userID int, date string) DailySummary {
	for i, summary := range *summaries {
		if v, ok := summary["userId"]; ok {
			if w, ok := summary["date"]; ok {
				if ((reflect.TypeOf(v).Kind() == reflect.Int && v.(int) == userID) || (reflect.TypeOf(v).Kind() == reflect.Float64 && int(v.(float64)) == userID)) && w.(string) == date {
					return (*summaries)[i]
				}
			}
		}
	}

	newSummary := make(DailySummary)
	newSummary["userId"] = userID
	newSummary["date"] = date

	*summaries = append(*summaries, newSummary)

	return (*summaries)[len(*summaries)-1]
}

func updateSummaryFields(summary DailySummary, event UserEvent, UserIdTimestampToEventTypeMap map[int]map[string]map[string]int) {
	if v, ok := summary["userId"]; ok {
		if w, ok := summary["date"]; ok {
			if reflect.TypeOf(v).Kind() == reflect.Int {
				summary[event.EventType] = UserIdTimestampToEventTypeMap[v.(int)][w.(string)][event.EventType]
			} else if reflect.TypeOf(v).Kind() == reflect.Float64 {
				summary[event.EventType] = UserIdTimestampToEventTypeMap[int(v.(float64))][w.(string)][event.EventType]
			}
		}
	}
}

func writeSummariesToFile(summaries DailySummaries, filePath string) error {
	file, err := json.MarshalIndent(summaries, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func hashFunc(event UserEvent, hashMap map[string]interface{}) error {
	out, err := json.Marshal(event)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%x", sha256.Sum256([]byte(string(out))))
	hashMap[key] = nil
	return nil
}

func writeEventHashToFile(hashMap map[string]interface{}) error {
	file, err := json.MarshalIndent(hashMap, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile("hash.json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}
