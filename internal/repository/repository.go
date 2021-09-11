package repository

import (
	"time"

	"github.com/shwethadia/HotelReservation/internal/models"
)

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)

	
	GetUserByID(id int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, password string) (int, string, error)
}
