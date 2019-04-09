package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func parseReservationRequests(filepath string) (Theater, []Reservation, error) {
	var event string
	var layout []string
	var theater Theater
	var requests []Reservation

	file, err := os.Open(filepath)
	if err != nil {
		return theater, requests, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		rowLayout := readSectionLayout(line)

		if len(line) > 0 && rowLayout != nil {
			layout = append(layout, line)
		} else if len(line) > 0 {
			if strings.HasPrefix(strings.ToLower(line), "event") {
				event = strings.Join(strings.Fields(line)[1:], " ")
			} else {
				name, seats, err := readRequestLine(line)
				if err != nil {
					continue
				}

				requests = append(requests, Reservation{Name: name, Seats: seats})
			}

		}
	}

	if len(layout) > 0 {
		theater = newTheater(strings.Join(layout, "\n"))
		theater.Event = event
	}
	/*
		theaterLayout := `
		6 6
		3 5 5 3
		4 6 6 4
		2 8 8 2
		6 6`

		requests = []Reservation{
			Reservation{Name: "Smith", Seats: 2},
			Reservation{Name: "Jones", Seats: 5},
			Reservation{Name: "Davis", Seats: 6},
			Reservation{Name: "Wilson", Seats: 100},
			Reservation{Name: "Johnson", Seats: 3},
			Reservation{Name: "Williams", Seats: 4},
			Reservation{Name: "Brown", Seats: 8},
			Reservation{Name: "Miller", Seats: 12},
		}

		theater = newTheater(theaterLayout)

	*/

	return theater, requests, err
}

func main() {
	filepath := flag.String("file", "", "The path to the ticket request file")
	verbose := flag.Bool("verbose", false, "Print extra information")
	flag.Parse()

	theater, requests, err := parseReservationRequests(*filepath)
	if err != nil {
		fmt.Printf("Unable to process ticket requests from %v -> %s\n", *filepath, err)
		os.Exit(1)
	}

	err = theater.refreshFromFile()

	if err != nil {
		fmt.Printf("Unable to load reservations from theater file %s: %s\n", theater.filePath(), err)
	}

	fmt.Println("Ticket Request Batch File Processing Results")
	fmt.Println("--------------------------------------------")
	for i := 0; i < len(requests); i++ {
		requests[i].process(&theater)
		fmt.Println(requests[i].toString())
	}
	fmt.Println("--------------------------------------------")

	if *verbose {
		fmt.Println(theater.toString(true))
		fmt.Printf("There are %d seats available...\n", theater.availableSeats())
	}

	err = theater.save()
	if err != nil {
		fmt.Printf("Unable to save reservations to %s, %s\n", theater.filePath(), err)
	} else if *verbose {
		fmt.Printf("Reservations saved successfully to %s\n", theater.filePath())
	}

}
