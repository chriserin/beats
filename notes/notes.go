package notes

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse notes string
func Parse(notes string) []int {
	switch {
	case strings.Contains(notes, "-"):
		return parseWithDash(notes)
	case strings.Contains(notes, ","):
		return parseWithCommas(notes)
	case notes == "DRUM":
		return []int{46, 48, 41, 58, 40, 49, 51, 42, 44, 39}
	}

	return parseWithDash("60-80")
}

func parseWithDash(notes string) []int {
	minmax := strings.Split(notes, "-")
	min := minmax[0]
	max := minmax[1]

	imin, errmin := strconv.Atoi(min)
	check(errmin, notes)
	imax, errmax := strconv.Atoi(max)
	check(errmax, notes)

	notesArray := []int{}

	for i := imin; i <= imax; i++ {
		notesArray = append(notesArray, i)
	}

	return notesArray
}

func parseWithCommas(notes string) []int {
	values := strings.Split(notes, ",")

	notesArray := []int{}

	for _, value := range values {
		if num, err := strconv.Atoi(value); err == nil {
			notesArray = append(notesArray, num)
		}
	}

	return notesArray
}

func check(err error, notes string) {
	if err != nil {
		panic(fmt.Sprintf("Unable to parse notes value '%s'", notes))
	}
}
