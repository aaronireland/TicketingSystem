package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Reservation represents a request to book seats in a theater
type Reservation struct {
	Name         string `json:"name"`
	Seats        int    `json:"seats"`
	Error        string `json:"error"`
	Confirmation string `json:"confirmation"`
}

func readRequestLine(line string) (string, int, error) {
	var name string
	var seats int
	var err error

	lineSlice := strings.Fields(line)

	seats, err = strconv.Atoi(lineSlice[len(lineSlice)-1])

	if err == nil {
		name = strings.Join(lineSlice[:len(lineSlice)-1], " ")

		if name == "" {
			err = errors.New("invalid name for ticket request")
		}
	}

	return name, seats, err

}

// toString builds a human-readable representation of the reservation request
func (r Reservation) toString() string {
	if r.Error != "" {
		return fmt.Sprintf("%s %s", r.Name, r.Error)
	} else if r.Confirmation != "" {
		return fmt.Sprintf("%s %s", r.Name, r.Confirmation)
	}

	// The request has not yet been processed
	return fmt.Sprintf("%s Pending Confirmation", r.Name)
}

// process attempts to reserve the requested number of seats as far front as possible.
// Updates the reservation with either an error or confirmation of success. If able,
// the ReservedSeats field of the selected Section will be updated with requester Name
// for each reserved seat.
func (r *Reservation) process(theater *Theater) {

	if r.Seats > theater.availableSeats() {
		r.Error = "Sorry, we can't handle your party."
		return
	}

	for row := 0; row < len(theater.Rows); row++ {
		for sec := 0; sec < len(theater.Rows[row].Sections); sec++ {
			section := &theater.Rows[row].Sections[sec]
			err := section.reserveSeats(r.Name, r.Seats)
			if err == nil {
				r.Confirmation = fmt.Sprintf("Row %d Section %d", row+1, sec+1)
				return
			}
		}
	}

	r.Error = "Call to split party."

}
