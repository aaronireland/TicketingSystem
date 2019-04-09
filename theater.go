package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Section has a number of seats and an array of reserved seats
type Section struct {
	Size          int      `json:"seats"`
	ReservedSeats []string `json:"reserved"`
}

// Row is an array of seat Sections
type Row struct {
	Sections []Section `json:"sections"`
}

// Theater is an array of rows
type Theater struct {
	Rows  []Row  `json:"rows"`
	Event string `json:"event"`
}

// availableSeats returns the number of unreserved seats in a section
func (section Section) availableSeats() int {
	return section.Size - len(section.ReservedSeats)
}

// reserveSeats sets the ReservedSeats field of the Section
func (section *Section) reserveSeats(name string, seats int) error {
	if section.availableSeats() < seats {
		return errors.New("not enough seats in section")
	}

	for seat := 0; seat < seats; seat++ {
		updatedSection := append(section.ReservedSeats, name)
		section.ReservedSeats = updatedSection
	}

	return nil
}

// toString renders a human-readable of the current reservations for the theater
func (theater Theater) toString(includeReservations bool) string {
	var layoutString string
	for i, row := range theater.Rows {
		var delimiter string
		if len(row.Sections) > 1 {
			delimiter = "\t"
		}
		for j, section := range row.Sections {
			layoutString += fmt.Sprintf("[%d-%d(%d)]", i+1, j+1, section.Size)
			if includeReservations {
				layoutString += strings.Join(section.ReservedSeats, ", ")
			}
			layoutString += delimiter
		}
		layoutString += "\n"
	}
	return layoutString
}

func (theater Theater) dataDirectory() string {
	return filepath.Join(".", "data", theater.Event)
}

// fileName generates a JSON filename based on a hash of the string representation of the theater
func (theater Theater) fileName() string {
	fileNameHash := fnv.New32a()
	fileNameHash.Write([]byte(theater.toString(false)))
	fileName := fmt.Sprintf("%d.json", fileNameHash.Sum32())

	return fileName
}

// Build a full path to the JSON file so that each theater reservation file is unique by event
func (theater Theater) filePath() string {
	return filepath.Join(theater.dataDirectory(), theater.fileName())
}

// save persists a JSON representation of the theater to the filesystem
func (theater Theater) save() error {
	json, err := json.Marshal(theater)
	if err == nil {
		if _, err := os.Stat(theater.dataDirectory()); os.IsNotExist(err) {
			_ = os.MkdirAll(theater.dataDirectory(), os.ModePerm)
		}

		err = ioutil.WriteFile(theater.filePath(), []byte(json), 0644)
	}

	return err
}

// refreshFromFile reads reservations from JSON file
func (theater Theater) refreshFromFile() error {
	if _, err := os.Stat(theater.dataDirectory()); os.IsNotExist(err) {
		return err
	}

	file, err := ioutil.ReadFile(theater.filePath())
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &theater)

	return err
}

// availableSeats returns the number of unreserved seats in the theater
func (theater Theater) availableSeats() int {
	var available int
	for _, row := range theater.Rows {
		for _, section := range row.Sections {
			available += section.availableSeats()
		}
	}

	return available
}

// readSectionLayout parses an array of seat counts from a line of Text
func readSectionLayout(line string) []int {
	var sections []int

	sectionsAsString := strings.Fields(line)

	for _, section := range sectionsAsString {
		sectionAsInt, err := strconv.Atoi(section)
		if err != nil {
			return nil
		}
		sections = append(sections, sectionAsInt)
	}

	return sections
}

// newTheater creates an instance of a theater from a string.
// Each line is a row and sections are delimited with spaces
// E.x. --->
// 1 2 3
// 4 4 6
// This would create a theater with two rows and six total sections with a
// total of 20 seats
func newTheater(layout string) Theater {

	var theater Theater
	rowLayout := strings.Split(layout, "\n")

	for _, layout := range rowLayout {
		if len(layout) > 0 {
			sectionLayout := readSectionLayout(layout)

			if len(sectionLayout) > 0 {
				var row Row

				for _, seats := range sectionLayout {
					section := Section{Size: seats}
					row.Sections = append(row.Sections, section)
				}

				if len(row.Sections) > 0 {
					theater.Rows = append(theater.Rows, row)
				}
			}
		}
	}

	return theater

}
