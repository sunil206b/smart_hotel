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

	AllReservations = `select rs.id, rs.first_name, rs.last_name, rs.email, rs.phone, rs.check_in, rs.check_out,
						rs.created_at, rs.updated_at, rs.room_id, rs.processed, r.id, r.room_name from reservations rs
						inner join rooms r on rs.room_id = r.id order by rs.check_in desc`

	AllNewReservations = `select rs.id, rs.first_name, rs.last_name, rs.email, rs.phone, rs.check_in, rs.check_out,
						rs.created_at, rs.updated_at, rs.room_id, r.id, r.room_name from reservations rs
						inner join rooms r on rs.room_id = r.id where rs.processed = 0 order by rs.check_in desc`

	GetReservationByID = `select rs.id, rs.first_name, rs.last_name, rs.email, rs.phone, rs.check_in, rs.check_out,
							rs.created_at, rs.updated_at, rs.room_id, rs.processed, r.id, r.room_name from reservations rs
							inner join rooms r on rs.room_id = r.id where rs.id = $1`

	UpdateReservation = `update reservations set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
							where id = $6`

	DeleteRoomRestriction = `delete from room_restrictions where reservation_id = $1`

	DeleteReservation = `delete from reservations where id = $1`

	UpdateProcessedReservation = `update reservations set processed = $1, updated_at = $2 where id = $3`

	AllRooms = `select id, room_name, created_at, updated_at from rooms order by room_name`

	GetRoomRestrictionsByDate = `select id, start_date, end_date, room_id, coalesce(reservation_id, 0), restriction_id from room_restrictions
									where $1 < end_date and $2 >= start_date and room_id = $3`

	CreateBlockForRoom = `insert into room_restrictions(start_date, end_date, created_at, updated_at, room_id, restriction_id)
 							values ($1, $2, $3,$4, $5, $6)`
	DeleteBlockByID = `delete from room_restrictions where id = $1`
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
	defer rows.Close()
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

// AllReservations returns a slice of all reservations
func (pg *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(AllReservations)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in AllReservations() method while preparing query to get all reservations: %v\n", err))
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in AllReservations() method while executing query to get all reservations: %v\n", err))
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("error in AllReservations() method while scanning rows for reservations: %v\n", err))
	}
	var reservations []models.Reservation
	for rows.Next() {
		var rs models.Reservation
		err = rows.Scan(&rs.ID, &rs.FirstName, &rs.LastName, &rs.Email, &rs.Phone,
			&rs.CheckInDate, &rs.CheckOutDate, &rs.CreatedAt, &rs.UpdatedAt, &rs.RoomID,
			&rs.Processed, &rs.Room.ID, &rs.Room.RoomName)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error in AllReservations() method while scanning each row for reservation: %v\n", err))
		}
		reservations = append(reservations, rs)
	}
	return reservations, nil
}

// AllNewReservations returns a slice of all new reservations
func (pg *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(AllNewReservations)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in AllReservations() method while preparing query to get all reservations: %v\n", err))
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in AllReservations() method while executing query to get all reservations: %v\n", err))
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("error in AllReservations() method while scanning rows for reservations: %v\n", err))
	}
	var reservations []models.Reservation
	for rows.Next() {
		var rs models.Reservation
		err = rows.Scan(&rs.ID, &rs.FirstName, &rs.LastName, &rs.Email, &rs.Phone,
			&rs.CheckInDate, &rs.CheckOutDate, &rs.CreatedAt, &rs.UpdatedAt, &rs.RoomID,
			&rs.Room.ID, &rs.Room.RoomName)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error in AllReservations() method while scanning each row for reservation: %v\n", err))
		}
		reservations = append(reservations, rs)
	}
	return reservations, nil
}

// GetReservationByID returns one reservation by ID
func (pg *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rs models.Reservation
	stmt, err := pg.DB.Prepare(GetReservationByID)
	if err != nil {
		return rs, errors.New(fmt.Sprintf("error in GetReservationByID() method while preparing query to get a reservation: %v\n", err))
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, id).Scan(&rs.ID, &rs.FirstName, &rs.LastName, &rs.Email, &rs.Phone,
		&rs.CheckInDate, &rs.CheckOutDate, &rs.CreatedAt, &rs.UpdatedAt, &rs.RoomID,
		&rs.Processed, &rs.Room.ID, &rs.Room.RoomName)
	if err != nil {
		return rs, errors.New(fmt.Sprintf("error in GetReservationByID() method while executing query to get a reservation: %v\n", err))
	}
	return rs, nil
}

//UpdateReservation updates the reservation in the database
func (pg *postgresDBRepo) UpdateReservation(res *models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(UpdateReservation)
	if err != nil {
		return errors.New(fmt.Sprintf("error in UpdateReservation() method while preparing query to update a reservation: %v\n", err))
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, res.FirstName, res.LastName, res.Email, res.Phone, time.Now(), res.ID)
	if err != nil {
		return errors.New(fmt.Sprintf("error in UpdateReservation() method while executing the update reservation: %v\n", err))
	}
	return nil
}

// DeleteReservation deletes reservation by id
func (pg *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt1, err := pg.DB.Prepare(DeleteRoomRestriction)
	if err != nil {
		return errors.New(fmt.Sprintf("error in DeleteReservation() method while preparing query to delete a room restriction: %v\n", err))
	}
	defer stmt1.Close()
	_, err = stmt1.ExecContext(ctx, id)
	if err != nil {
		return errors.New(fmt.Sprintf("error in DeleteReservation() method while executing the delete room restriction: %v\n", err))
	}
	stmt2, err := pg.DB.Prepare(DeleteReservation)
	if err != nil {
		return errors.New(fmt.Sprintf("error in DeleteReservation() method while preparing query to delete a reservation: %v\n", err))
	}
	defer stmt2.Close()
	_, err = stmt2.ExecContext(ctx, id)
	if err != nil {
		return errors.New(fmt.Sprintf("error in DeleteReservation() method while executing the delete reservation: %v\n", err))
	}
	return nil
}

//UpdateProcessedReservation updates processed for a reservation
func (pg *postgresDBRepo) UpdateProcessedReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(UpdateProcessedReservation)
	if err != nil {
		return errors.New(fmt.Sprintf("error in UpdateProcessedReservation() method while preparing query to update a reservation: %v\n", err))
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, processed, time.Now(), id)
	if err != nil {
		return errors.New(fmt.Sprintf("error in UpdateProcessedReservation() method while executing the update reservation: %v\n", err))
	}
	return nil
}

// AllRooms returns all the rooms
func (pg *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(AllRooms)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in AllRooms() method while preparing query to get all rooms: %v\n", err))
	}
	defer stmt.Close()
	var rooms []models.Room
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in AllRooms() method while executing query to get all rooms: %v\n", err))
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("error in AllRooms() method while scanning rows to get all rooms: %v\n", err))
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
		if err = rows.Err(); err != nil {
			return nil, errors.New(fmt.Sprintf("error in AllRooms() method while scanning each row to get a rooms: %v\n", err))
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

//GetRestrictionsForRoomByDate returns all restrictions for a particular room by date range
func (pg *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(GetRoomRestrictionsByDate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in GetRestrictionsForRoomByDate() method while preparing query to get room restrictions: %v\n", err))
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, start, end, roomID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in GetRestrictionsForRoomByDate() method while executing query to get room restrictions: %v\n", err))
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("error in GetRestrictionsForRoomByDate() method while scanning all rows: %v\n", err))
	}
	var restrictions []models.RoomRestriction
	for rows.Next() {
		var res models.RoomRestriction
		err = rows.Scan(&res.ID, &res.StartDate, &res.EndDate, &res.RoomID, &res.ReservationID, &res.RestrictionID)
		if err = rows.Err(); err != nil {
			return nil, errors.New(fmt.Sprintf("error in GetRestrictionsForRoomByDate() method while scanning each row to get room restriction: %v\n", err))
		}
		restrictions = append(restrictions, res)
	}
	return restrictions, nil
}

//CreateBlockForRoom will create the restriction for a given room
func (pg postgresDBRepo) CreateBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(CreateBlockForRoom)
	if err != nil {
		return errors.New(fmt.Sprintf("error in CreateBlockForRoom() method while preparing query to create room restriction: %v\n", err))
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, startDate, startDate.AddDate(0, 0, 1), time.Now(), time.Now(), id, 2)
	if err != nil {
		return errors.New(fmt.Sprintf("error in CreateBlockForRoom() method while executing query to create room restriction: %v\n", err))
	}
	return nil
}

//DeleteBlockByID deletes a room restriction
func (pg postgresDBRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := pg.DB.Prepare(DeleteBlockByID)
	if err != nil {
		return errors.New(fmt.Sprintf("error in DeleteBlockByID() method while preparing query to delete room restriction: %v\n", err))
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return errors.New(fmt.Sprintf("error in DeleteBlockByID() method while executing query to delete room restriction: %v\n", err))
	}
	return nil
}
