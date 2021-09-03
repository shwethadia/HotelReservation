package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi"
	"github.com/shwethadia/HotelReservation/internal/config"
)

func TestRoutes(t *testing.T) {

	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {

	case *chi.Mux:
		//Do Nothing /Test Passed

	default:
		t.Error(fmt.Sprintf("Type is not *chi.mux , type is %T", v))
	}
}
