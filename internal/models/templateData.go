package models

import "github.com/shwethadia/HotelReservation/internal/forms"

//TemplateData holds data sent from handlres to templates
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FlatMap         map[string]float32
	Data            map[string]interface{}
	CSRFtoken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
}
