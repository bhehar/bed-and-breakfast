package repository

import (
	"time"

	"github.com/bhehar/bed-and-breakfast/internal/models"
)

type DatabaseRepo interface {
	GetRoomById(id int) (models.Room, error)
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesAndRoomId(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityByDates(start, end time.Time) ([]models.Room, error)
}
