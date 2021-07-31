package repository

import (
	"github.com/sunil206b/smart_booking/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool

	CreateReservation(res *models.Reservation) error
	CreateRoomRestriction(r *models.RoomRestriction) error
	SearchAvailabilityByDatesByRoom(start, end time.Time, roomID int) (bool, error)
	SearchAllAvailableRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(user *models.User) error
	Authenticate(email, testPass string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(res *models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedReservation(id, processed int) error
	AllRooms() ([]models.Room, error)
	GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error)
}
