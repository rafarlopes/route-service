package route

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/rafarlopes/route-service/internal/osrm"
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

func TestSortReturnsSortestDurationAndDistance(t *testing.T) {
	routes := []*osrm.Route{
		{
			Duration: 30.0,
			Distance: 25.0,
		},
		{
			Duration: 15.0,
			Distance: 10.0,
		},
		{
			Duration: 15.0,
			Distance: 12.0,
		},
		{
			Duration: 20.0,
			Distance: 20.0,
		},
	}

	sort.Sort(ByDurationAndDistance(routes))
	assertSorting(t, routes[0], 15.0, 10.0)
	assertSorting(t, routes[1], 15.0, 12.0)
	assertSorting(t, routes[2], 20.0, 20.0)
	assertSorting(t, routes[3], 30.0, 25.0)
}

func assertSorting(t *testing.T, element *osrm.Route, expectedDuration, expectedDistance float64) {
	if element.Duration != expectedDuration && element.Distance != expectedDistance {
		t.Errorf("sort did not match - expect duration %f got %f and expected distance %f and got %f", expectedDuration, element.Duration, expectedDistance, element.Distance)
	}
}
