package dbrepo

import (
	"errors"
	"time"

	"github.com/shwethadia/HotelReservation/internal/models"
)

func (m *postgresDBTestRepo) InsertReservation(res models.Reservation) (int, error) {

	//if the room id is 2 , then fail otherwise pass

	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

//InsertRoomRestiction inserts a room restriction into the databases
func (m *postgresDBTestRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	if r.RoomID == 1000 {
		return errors.New("some error")
	}

	return nil
}

//SerachAvailabilityByDate returns True if availability exists
func (m *postgresDBTestRepo) SearchAvailabilityByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil

}

//SerachAvailability for all rooms returns a slice of available rooms if any given date range
func (m *postgresDBTestRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room
	return rooms, nil
}

//GetRoomByID gets a room by ID
func (m *postgresDBTestRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room
	if id > 2 {

		return room, errors.New("some error")
	}

	return room, nil
}

//GetUserByID gets a user by ID
func (m *postgresDBTestRepo) GetUserByID(id int) (models.User, error) {

	var user models.User

	return user, nil
}

//UpdateUser updates a user in the database
func (m *postgresDBTestRepo) UpdateUser(user models.User) error {

	return nil

}

//Authenticate authenticate user
func (m *postgresDBTestRepo) Authenticate(email, password string) (int, string, error) {

	var id int

	var hashedPassword string

	return id, hashedPassword, nil

}
