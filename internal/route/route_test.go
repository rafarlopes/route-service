package route

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHappyPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(RoutesHandler)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("status code does not match: expect %v got %v", http.StatusOK, status)
	}
}

func TestSrcParameterIsRequired(t *testing.T) {
	req, err := http.NewRequest("GET", "/routes?dst=13.397634,52.529407&dst=13.428555,52.523219", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(RoutesHandler)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("status code does not match: expect %v got %v", http.StatusBadRequest, status)
	}
}

func TestDstParameterIsRequired(t *testing.T) {
	req, err := http.NewRequest("GET", "/routes?src=13.388860,52.517037", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(RoutesHandler)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("status code does not match: expect %v got %v", http.StatusBadRequest, status)
	}
}

func TestCanParseAndValidateLongAndLatString(t *testing.T) {
	long, lat, err := parseAndValidateLongitudeAndLatitude("13.397634,52.529407")

	if err != nil {
		t.Error("parseAndValidateLongitudeAndLatitude should not return an error")
	}

	if long != 13.397634 {
		t.Errorf("Longitude does not match: expected %f got %f", 13.397634, long)
	}

	if lat != 52.529407 {
		t.Errorf("Latitude does not match: expected %f got %f", 52.529407, lat)
	}
}

func TestReturnsErrorWhenStringIsNotValidLongAndLat(t *testing.T) {
	_, _, err := parseAndValidateLongitudeAndLatitude("invalid")

	if err == nil {
		t.Error("parseAndValidateLongitudeAndLatitude should return an error when the string is invalid")
	}
}

func TestReturnsErrorWhenLongitudeIsOutOfRange(t *testing.T) {
	_, _, err := parseAndValidateLongitudeAndLatitude("1300.397634,52.529407")

	if err == nil {
		t.Error("parseAndValidateLongitudeAndLatitude should return an error when longitude is out of range")
	}
}

func TestReturnsErrorWhenLatitudeIsOutOfRange(t *testing.T) {
	_, _, err := parseAndValidateLongitudeAndLatitude("13.397634,520.529407")

	if err == nil {
		t.Error("parseAndValidateLongitudeAndLatitude should return an error when latitude is out of range")
	}
}
