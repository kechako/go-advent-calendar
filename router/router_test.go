package router

import (
	"net/http"
	"reflect"
	"testing"
)

func TestGetPathParams(t *testing.T) {
	r1, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	params1, err := GetPathParams(r1, ":year")
	if err != nil {
		t.Error(err)
	}

	expect1 := map[string]string{"year": ""}
	if !reflect.DeepEqual(params1, expect1) {
		t.Errorf("Want %v\ngot %v", expect1, params1)
	}

	r2, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/2016", nil)
	params2, err := GetPathParams(r2, ":year")
	if err != nil {
		t.Error(err)
	}

	expect2 := map[string]string{"year": "2016"}
	if !reflect.DeepEqual(params2, expect2) {
		t.Errorf("Want %v\ngot %v", expect2, params2)
	}

	r3, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/entries/2016/24", nil)
	params3, err := GetPathParams(r3, "entries", ":year", ":day")
	if err != nil {
		t.Error(err)
	}

	expect3 := map[string]string{"year": "2016", "day": "24"}
	if !reflect.DeepEqual(params3, expect3) {
		t.Errorf("Want %v\ngot %v", expect3, params3)
	}
}
