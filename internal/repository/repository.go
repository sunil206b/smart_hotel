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
}
