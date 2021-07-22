package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/sunil206b/smart_booking/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	InsertReservation = `insert into reservations(first_name, last_name, email, phone, check_in, check_out, created_at, updated_at, room_id) 
						values($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	InsertRoomRestriction = `insert into room_restrictions(start_date, end_date, created_at, updated_at, room_id, reservation_id, restriction_id)
							values($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	SearchAvailableRoomByDate = `select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date`

	SearchAllAvailableRooms = `select r.id, r.room_name from rooms r
								where r.id not in (select rr.room_id from room_restrictions rr
								where $1 < rr.end_date and $2 > rr.start_date)`

	SearchRoomByID = `select id, room_name, created_at, updated_at from rooms where id = $1`

	GetUserByID = `select id, first_name, last_name, email, password, access_level, created_at, 
       				updated_at from users where id = $1`

	UpdateUser = `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, 
                 	updated_at = $5 where id = $6`

	GetUserByEmail = `select id, password from users where email = $1`
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
func (pg *postgresDBRepo) SearchAvailabilityByDatesByRoom(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(SearchAvailableRoomByDate)
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

// SearchAllAvailableRooms returns all available rooms if any, with the given date range
func (pg *postgresDBRepo) SearchAllAvailableRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(SearchAllAvailableRooms)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in SearchAllAvailableRooms() method while preparing query to search all available rooms: %v\n", err))
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, start, end)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in SearchAllAvailableRooms() method while executing query to search all available rooms: %v\n", err))
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("error in SearchAllAvailableRooms() method while checking for errors in the query rows: %v\n", err))
	}

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error in SearchAllAvailableRooms() method while scanning results into rooms model: %v\n", err))
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (pg *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	room := models.Room{}
	stmt, err := pg.DB.Prepare(SearchRoomByID)
	if err != nil {
		return room, errors.New(fmt.Sprintf("error in GetRoomByID() method while preparing query to search for a room: %v\n", err))
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, id).Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, errors.New(fmt.Sprintf("error in GetRoomByID() method while executing the query: %v\n", err))
	}
	return room, nil
}

// GetUserByID returns user with the given id
func (pg *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user := models.User{}
	stmt, err := pg.DB.Prepare(GetUserByID)
	if err != nil {
		return user, errors.New(fmt.Sprintf("error in GetUserByID() method while preparing query to search for a user: %v\n", err))
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, errors.New(fmt.Sprintf("error in GetUserByID() method while executing the query: %v\n", err))
	}
	return user, nil
}

//UpdateUser updates the user in the database
func (pg *postgresDBRepo) UpdateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(UpdateUser)
	if err != nil {
		return errors.New(fmt.Sprintf("error in UpdateUser() method while preparing query to update a user: %v\n", err))
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, user.FirstName, user.LastName, user.Email, user.AccessLevel, user.UpdatedAt, user.ID)
	if err != nil {
		return errors.New(fmt.Sprintf("error in UpdateUser() method while executing the query: %v\n", err))
	}
	return nil
}

//Authenticate authenticates the user
func (pg *postgresDBRepo) Authenticate(email, testPass string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	var hashedPass string
	stmt, err := pg.DB.Prepare(GetUserByEmail)
	if err != nil {
		return id, "", errors.New(fmt.Sprintf("error in Authenticate() method while preparing query to get a user: %v\n", err))
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, email).Scan(&id, &hashedPass)
	if err != nil {
		return id, "", errors.New(fmt.Sprintf("error in Authenticate() method while executing the query: %v\n", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(testPass))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPass, nil
}
