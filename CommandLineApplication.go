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
	_, err = reader.Read() // Skip header
	if err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		siteID, _ := strconv.Atoi(record[0])
		fxiletID, _ := strconv.Atoi(record[1])
		relevantComputerCount, _ := strconv.Atoi(record[4])
		entries = append(entries, Entry{
			SiteID:                siteID,
			FxiletID:              fxiletID,
			Name:                  record[2],
			Criticality:           record[3],
			RelevantComputerCount: relevantComputerCount,
		})
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
	for _, entry := range entries {
		writer.Write([]string{
			strconv.Itoa(entry.SiteID),
			strconv.Itoa(entry.FxiletID),
			entry.Name,
			entry.Criticality,
			strconv.Itoa(entry.RelevantComputerCount),
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
	for _, entry := range entries {
		fmt.Printf("SiteID: %d, FxiletID: %d, Name: %s, Criticality: %s, Computers: %d\n",
			entry.SiteID, entry.FxiletID, entry.Name, entry.Criticality, entry.RelevantComputerCount)
	}
}

// QueryEntry searches for entries by name or criticality.
func QueryEntry(entries []Entry, query string) {
	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(entry.Criticality), strings.ToLower(query)) {
			fmt.Printf("SiteID: %d, FxiletID: %d, Name: %s, Criticality: %s, Computers: %d\n",
				entry.SiteID, entry.FxiletID, entry.Name, entry.Criticality, entry.RelevantComputerCount)
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
func AddEntry(entries []Entry, siteID, fxiletID int, name, criticality string, relevantComputerCount int) []Entry {
	return append(entries, Entry{
		SiteID:                siteID,
		FxiletID:              fxiletID,
		Name:                  name,
		Criticality:           criticality,
		RelevantComputerCount: relevantComputerCount,
	})
}

// DeleteEntry deletes an entry by FxiletID.
func DeleteEntry(entries []Entry, fxiletID int) []Entry {
	for i, entry := range entries {
		if entry.FxiletID == fxiletID {
			return append(entries[:i], entries[i+1:]...)
		}
	}
	fmt.Println("Entry not found.")
	return entries
}

func main() {
	const filename = "fixlets.csv" // CSV file to be used

	// Read the existing CSV data
	entries, err := ReadCSV(filename)
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	// Command-line interactions
	var command string
	for {
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
			var siteID, fxiletID, relevantComputerCount int
			var name, criticality string
			fmt.Println("Enter SiteID:")
			fmt.Scanln(&siteID)
			fmt.Println("Enter FxiletID:")
			fmt.Scanln(&fxiletID)
			fmt.Println("Enter name:")
			fmt.Scanln(&name)
			fmt.Println("Enter criticality:")
			fmt.Scanln(&criticality)
			fmt.Println("Enter relevant computer count:")
			fmt.Scanln(&relevantComputerCount)
			entries = AddEntry(entries, siteID, fxiletID, name, criticality, relevantComputerCount)
			if err := WriteCSV(filename, entries); err != nil {
				fmt.Println("Error saving CSV file:", err)
			} else {
				fmt.Println("Entry added successfully.")
			}
		case "delete":
			var fxiletID int
			fmt.Println("Enter FxiletID of entry to delete:")
			fmt.Scanln(&fxiletID)
			entries = DeleteEntry(entries, fxiletID)
			if err := WriteCSV(filename, entries); err != nil {
				fmt.Println("Error saving CSV file:", err)
			} else {
				fmt.Println("Entry deleted successfully.")
			}
		case "exit":
			fmt.Println("Exiting program.")
			return
		default:
			fmt.Println("Invalid command.")
		}
	}
}
