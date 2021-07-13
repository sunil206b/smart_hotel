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
}
