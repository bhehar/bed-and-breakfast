package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/bhehar/bed-and-breakfast/internal/models"
)

func (m *testingDbRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}

// InsertReservation inserts a reservation into the database
func (m *testingDbRepo) InsertReservation(res models.Reservation) (int, error) {
	// setup for testing
	if res.RoomID == -1 {
		return -1, errors.New("failing for test purposes")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testingDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == -2 {
		return errors.New("failing for test purposes")
	}
	return nil
}

// SearchAvailabilityByDatesAndRoomId checks the availability of a room
func (m *testingDbRepo) SearchAvailabilityByDatesAndRoomId(start, end time.Time, roomId int) (bool, error) {
	// set up a test time
	str := "2049-12-31"
	t, err := time.Parse(time.DateOnly, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2060-01-01 as start
	testDateToFail, err := time.Parse(time.DateOnly, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start.Equal(testDateToFail) {
		return false, errors.New("some error")
	}

	// if the start date is after 2049-12-31, then return false,
	// indicating no availability;
	if start.After(t) {
		return false, nil
	}

	// otherwise, we have availability
	return true, nil
}

// SearchAvailabilityByDates returns a slice of available rooms for given date range
func (m *testingDbRepo) SearchAvailabilityByDates(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	// if the start date is after 2049-12-31, then return empty slice,
	// indicating no rooms are available;
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return rooms, errors.New("some error")
	}

	if start.After(t) {
		return rooms, nil
	}

	// otherwise, put an entry into the slice, indicating that some room is
	// available for search dates
	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}
