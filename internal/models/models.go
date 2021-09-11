package models

import "time"

//Users is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt    time.Time
}

//Room
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt  time.Time
}

//Restriction is the Restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt        time.Time
}

//Reservations is the Reservation Model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt  time.Time
	Room      Room
}

//RoomRestrictions is the RoomRestriction Model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	CreatedAt     time.Time
	UpdatedAt      time.Time
	Room          Room
	ReservationID int
	Reservation   Reservation
	RestrictionID int
	Restriction   Restriction
}
