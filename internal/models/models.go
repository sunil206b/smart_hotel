package models

import (
	"time"
)

//User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Room is the rooms model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Restriction is the restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//Reservation is the reservations model
type Reservation struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
	RoomID       int       `json:"room_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Processed    int       `json:"processed"`
	Room         Room      `json:"-"`
}

//RoomRestriction is the room_restrictions model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

// MailData holds an email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
