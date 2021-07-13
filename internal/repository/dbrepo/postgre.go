package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/sunil206b/smart_booking/internal/models"
	"time"
)

const (
	InsertReservation = `insert into reservations(first_name, last_name, email, phone, check_in, check_out, created_at, updated_at, room_id) 
						values($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	InsertRoomRestriction = `insert into room_restrictions(start_date, end_date, created_at, updated_at, room_id, reservation_id, restriction_id)
							values($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	SearchAvailableRooms = `select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date`
)

func (pg *postgresDBRepo) AllUsers() bool {
	return false
}

//CreateReservation creates reservation record in the database
func (pg *postgresDBRepo) CreateReservation(res *models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(InsertReservation)
	if err != nil {
		return errors.New(fmt.Sprintf("error in CreateReservation() method while preparing create reservations query: %v\n", err))
	}
	defer stmt.Close()

	reservationID := 0
	err = stmt.QueryRowContext(ctx, res.FirstName, res.LastName, res.Email, res.Phone, res.CheckInDate,
		res.CheckOutDate, res.CreatedAt, res.UpdatedAt, res.RoomID).Scan(&reservationID)
	if err != nil {
		return errors.New(fmt.Sprintf("error in CreateReservation() method while creating reservations: %v\n", err))
	}
	res.ID = reservationID
	return nil
}

// CreateRoomRestriction create room restrictions in the database
func (pg *postgresDBRepo) CreateRoomRestriction(r *models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(InsertRoomRestriction)
	if err != nil {
		return errors.New(fmt.Sprintf("error in CreateRoomRestriction() method while preparing create room restrictions query: %v\n", err))
	}
	defer stmt.Close()

	restrictionID := 0
	err = stmt.QueryRowContext(ctx, r.StartDate, r.EndDate, r.CreatedAt, r.UpdatedAt, r.RoomID, r.ReservationID, r.RestrictionID).Scan(&restrictionID)
	if err != nil {
		return errors.New(fmt.Sprintf("error in CreateRoomRestriction() method while creating room restriction: %v\n", err))
	}
	r.ID = restrictionID
	return nil
}

// SearchAvailabilityByDatesByRoom returns true if room available, and false if room not available in the given roomID
func (pg postgresDBRepo) SearchAvailabilityByDatesByRoom(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(SearchAvailableRooms)
	if err != nil {
		return false, errors.New(fmt.Sprintf("error in SearchAvailabilityByDatesByRoom() method while preparing search available rooms query: %v\n", err))
	}
	defer stmt.Close()

	numRows := 0
	err = stmt.QueryRowContext(ctx, roomID, start, end).Scan(&numRows)
	if err != nil {
		return false, errors.New(fmt.Sprintf("error in SearchAvailabilityByDatesByRoom() method while executing search available rooms query: %v\n", err))
	}
	return numRows == 0, nil
}
