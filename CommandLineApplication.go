package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Entry represents an individual record in the CSV file.
type Entry struct {
	SiteID                int
	FxiletID              int
	Name                  string
	Criticality           string
	RelevantComputerCount int
}

// ReadCSV reads the CSV file and returns a slice of Entry structs.
func ReadCSV(filename string) ([]Entry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []Entry
	reader := csv.NewReader(file)
	_, _ = reader.Read() // Skip header
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		siteID, _ := strconv.Atoi(record[0])
		fxiletID, _ := strconv.Atoi(record[1])
		relevantComputerCount, _ := strconv.Atoi(record[4])
		entries = append(entries, Entry{siteID, fxiletID, record[2], record[3], relevantComputerCount})
	}
	return entries, nil
}

// WriteCSV writes the list of entries to the CSV file.
func WriteCSV(filename string, entries []Entry) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Write([]string{"SiteID", "FxiletID", "Name", "Criticality", "RelevantComputerCount"})
	for _, e := range entries {
		writer.Write([]string{
			strconv.Itoa(e.SiteID), strconv.Itoa(e.FxiletID), e.Name, e.Criticality, strconv.Itoa(e.RelevantComputerCount),
		})
	}
	writer.Flush()
	return nil
}

// ListEntries displays all entries in the CSV file.
func ListEntries(entries []Entry) {
	if len(entries) == 0 {
		fmt.Println("No entries available.")
		return
	}
	for _, e := range entries {
		fmt.Printf("SiteID: %d, FxiletID: %d, Name: %s, Criticality: %s, Computers: %d\n", e.SiteID, e.FxiletID, e.Name, e.Criticality, e.RelevantComputerCount)
	}
}

// QueryEntry searches for entries by name or criticality.
func QueryEntry(entries []Entry, query string) {
	query = strings.ToLower(query)
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Name), query) || strings.Contains(strings.ToLower(e.Criticality), query) {
			fmt.Printf("SiteID: %d, FxiletID: %d, Name: %s, Criticality: %s, Computers: %d\n", e.SiteID, e.FxiletID, e.Name, e.Criticality, e.RelevantComputerCount)
			return
		}
	}
	fmt.Println("No entries found.")
}

// SortEntries sorts entries by the relevant computer count in ascending order.
func SortEntries(entries []Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].RelevantComputerCount < entries[j].RelevantComputerCount
	})
}

// AddEntry adds a new entry to the list.
func AddEntry(entries []Entry) ([]Entry, error) {
	var siteID, fxiletID, relevantComputerCount int
	var name, criticality string
	fmt.Println("Enter SiteID, FxiletID, Name, Criticality, RelevantComputerCount:")
	if _, err := fmt.Scanf("%d %d %s %s %d", &siteID, &fxiletID, &name, &criticality, &relevantComputerCount); err != nil {
		return entries, err
	}
	entries = append(entries, Entry{siteID, fxiletID, name, criticality, relevantComputerCount})
	return entries, nil
}

// DeleteEntry deletes an entry by FxiletID.
func DeleteEntry(entries []Entry, fxiletID int) ([]Entry, bool) {
	for i, e := range entries {
		if e.FxiletID == fxiletID {
			entries = append(entries[:i], entries[i+1:]...)
			return entries, true
		}
	}
	return entries, false
}

func main() {
	const filename = "fixlets.csv"
	// Read the existing CSV data
	entries, err := ReadCSV(filename)
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}
	// Command-line interactions
	for {
		var command string
		fmt.Println("\nChoose an operation: list, query, add, delete, sort, exit")
		fmt.Scanln(&command)

		switch command {
		case "list":
			ListEntries(entries)
		case "query":
			var query string
			fmt.Println("Enter name or criticality to query:")
			fmt.Scanln(&query)
			QueryEntry(entries, query)
		case "sort":
			SortEntries(entries)
			ListEntries(entries)
		case "add":
			entries, err = AddEntry(entries)
			if err != nil {
				fmt.Println("Error adding entry:", err)
			} else {
				WriteCSV(filename, entries)
				fmt.Println("Entry added.")
			}
		case "delete":
			var fxiletID int
			fmt.Println("Enter FxiletID to delete:")
			fmt.Scanln(&fxiletID)
			entries, found := DeleteEntry(entries, fxiletID)
			if found {
				WriteCSV(filename, entries)
				fmt.Println("Entry deleted.")
			} else {
				fmt.Println("Entry not found.")
			}
		case "exit":
			fmt.Println("Exiting program.")
			return
		default:
			fmt.Println("Invalid command.")
		}
	}
}
