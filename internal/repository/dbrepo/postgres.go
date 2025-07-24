package dbrepo

import (
	"context"
	"time"

	"github.com/bhehar/bed-and-breakfast/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		SELECT
			id,
			room_name
		FROM
			rooms
		WHERE
			id = $1`

	var room models.Room
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&room.ID, &room.RoomName)
	if err != nil {
		return room, err
	}
	return room, nil
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int

	stmt := `insert into reservations (
			first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO room_restrictions
			(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesAndRoomId checks the availability of a room
func (m *postgresDBRepo) SearchAvailabilityByDatesAndRoomId(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			select
				count(*)
			from
				room_restrictions
			where
				room_id = $1
				and $2 < end_date
				and $3 > start_date;`

	var count int
	err := m.DB.QueryRowContext(ctx, query, roomId, start, end).Scan(&count)
	if err != nil {
		return false, err
	}

	return (count == 0), nil
}

// SearchAvailabilityByDates returns a slice of available rooms for given date range
func (m *postgresDBRepo) SearchAvailabilityByDates(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select
			id,
			room_name
		from
			rooms r
		where 
			id not in (
			select
				rr.room_id
			from
				room_restrictions rr
			where
				$1 < end_date
				and $2 > start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return rooms, nil
}
