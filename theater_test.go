package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

// Section Struct Tests...
func TestAvailableSeats(t *testing.T) {

	var tests = []struct {
		sectionSize            int
		reservedSeats          []string
		expectedAvailableSeats int
	}{
		{0, []string{}, 0},
		{2, []string{}, 2},
		{2, []string{"Test"}, 1},
		{2, []string{"Test", "Test"}, 0},
	}

	for _, test := range tests {
		section := Section{Size: test.sectionSize, ReservedSeats: test.reservedSeats}
		assert.Equal(t, section.availableSeats(), test.expectedAvailableSeats)
	}
}

func TestReserveSeats(t *testing.T) {
	var tests = []struct {
		sectionSize            int
		reservedSeats          []string
		requestedSeats         int
		expectedAvailableSeats int
		expectedError          bool
		description            string
	}{
		{0, []string{}, 1, 0, true, "Empty section of size 0, requesting 1 seat"},
		{2, []string{}, 2, 0, false, "Empty section of size 2, requesting 2 seats"},
		{2, []string{}, 3, 2, true, "Empty section of size 2, requesting 3 seats"},
		{2, []string{"Test"}, 1, 0, false, "1 available in section of size 2, requesting 1 seat"},
		{2, []string{"Test"}, 100, 1, true, "1 available in section of size 2, requesting 100 seats"},
		{100, []string{"Test"}, 1, 98, false, "99 available in section of size 100, requesting 1 seat"},
	}

	for _, test := range tests {
		section := Section{Size: test.sectionSize, ReservedSeats: test.reservedSeats}
		err := section.reserveSeats("Test", test.requestedSeats)
		if test.expectedError {
			assert.NotNil(t, err, test.description)
		} else {
			assert.Nil(t, err, test.description)
		}
		assert.Equal(t, section.availableSeats(), test.expectedAvailableSeats, test.description)
	}
}

// Theater Struct Tests...
func TestNewTheater(t *testing.T) {
	var tests = []struct {
		layout          string
		expectedTheater Theater
		expectedSeats   int
		description     string
	}{
		{
			"", Theater{}, 0, "Empty theater",
		},
		{
			"1", Theater{Rows: []Row{Row{[]Section{Section{Size: 1}}}}}, 1,
			"Theater with 1 row, 1 section, and 1 seat",
		},
		{
			"\n1 2 3",
			Theater{Rows: []Row{Row{[]Section{Section{Size: 1}, Section{Size: 2}, Section{Size: 3}}}}},
			6,
			"Theater with 1 row, 3 sections, and 6 seats",
		},
		{
			"1 1 3\n12 2 1\n 2 2 invalid row",
			Theater{Rows: []Row{
				Row{[]Section{Section{Size: 1}, Section{Size: 1}, Section{Size: 3}}},
				Row{[]Section{Section{Size: 12}, Section{Size: 2}, Section{Size: 1}}},
			}},
			20,
			"Theater with 2 rows, 6 sections, and 20 seats",
		},
		{
			"1 2 3\n   3 2 1\n\t10 1 11 100\n",
			Theater{Rows: []Row{
				Row{[]Section{Section{Size: 1}, Section{Size: 2}, Section{Size: 3}}},
				Row{[]Section{Section{Size: 3}, Section{Size: 2}, Section{Size: 1}}},
				Row{[]Section{Section{Size: 10}, Section{Size: 1}, Section{Size: 11}, Section{Size: 100}}},
			}},
			134,
			"Theater with 3 rows, 10 sections, and 134 seats",
		},
	}

	for _, test := range tests {
		theater := newTheater(test.layout)
		assert.Equal(t, test.expectedSeats, theater.availableSeats(), test.description)
		assert.True(
			t, cmp.Equal(theater, test.expectedTheater),
			fmt.Sprintf(
				"Testing %s\nExpected Theater:\n%s\n\nActual Theater:\n%s\n\n",
				test.description, test.expectedTheater.toString(false), theater.toString(false),
			),
		)
	}
}
