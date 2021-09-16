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

//AllReservations returns a slice of all reservations
func (m *postgresDBTestRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil

}

//AllNewReservations returns a slice of all reservations
func (m *postgresDBTestRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation
	return reservations, nil

}

//getReservationByID returns the reservations by user ID
func (m *postgresDBTestRepo) GetReservationByID(id int) (models.Reservation, error) {

	var reservation models.Reservation
	return reservation, nil

}

//UpdateUser updates a reservation in the database
func (m *postgresDBTestRepo) UpdateReservation(reservation models.Reservation) error {

	return nil

}

//DeleteReservation deletes one reservation by ID
func (m *postgresDBTestRepo) DeleteReservation(id int) error {

	return nil

}

//UpdateProcessedForReservation updates processed for a reservation by id
func (m *postgresDBTestRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil

}

func (m *postgresDBTestRepo) AllRooms() ([]models.Room, error) {

	var rooms []models.Room
	return rooms, nil

}

//GetRestrictionsForRoomByDate returns restrictions for a room by date range
func (m *postgresDBTestRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	var restrictions []models.RoomRestriction
	return restrictions, nil

}

//InsertBlockForRoom inserts a room restrictions
func (m *postgresDBTestRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	return nil
}

//DeleteBlockForRoom delete a room restrictions
func (m *postgresDBTestRepo) DeleteBlockForRoom(id int) error {

	return nil
}
