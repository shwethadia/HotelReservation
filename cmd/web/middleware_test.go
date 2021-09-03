package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	var my_handler myHandler
	h := NoSurf(&my_handler)

	switch v := h.(type) {

	case http.Handler:
		//Do Nothing

	default:
		t.Error(fmt.Sprintf("Type is not Http Handler, but got %T", v))
	}

}

func TestSession(t *testing.T) {

	var my_handler myHandler
	h := SessionLoad(&my_handler)

	switch v := h.(type) {

	case http.Handler:
		//Do Nothing

	default:
		t.Error(fmt.Sprintf("Type is not Http Handler, but got %T", v))
	}

}
