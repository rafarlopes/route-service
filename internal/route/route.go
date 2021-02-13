package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/rafarlopes/route-service/internal/osrm"
)

type (
	// Response struct represents the API output for the routes request
	Response struct {
		Source string        `json:"source"`
		Routes []*osrm.Route `json:"routes"`
	}

	// coordinates is a internal struct used to store the parsed destinations
	coordinates struct {
		Longitude float64
		Latitude  float64
	}

	// ByDurationAndDistance is used to sort osrm.Route first by duration and then by distance
	ByDurationAndDistance []*osrm.Route
)

func (r ByDurationAndDistance) Len() int {
	return len(r)
}

func (r ByDurationAndDistance) Less(i, j int) bool {
	return r[i].Duration < r[j].Duration && r[i].Distance < r[j].Distance
}

func (r ByDurationAndDistance) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// RoutesHandler handles the API requests for the /routes
func RoutesHandler(w http.ResponseWriter, r *http.Request) {
	source, ok := r.URL.Query()["src"]
	if !ok || len(source) != 1 {
		handleError(w, http.StatusBadRequest, "one src parameter must be specified")
		return
	}

	srcLong, srcLat, err := parseAndValidateLongitudeAndLatitude(source[0])
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("invalid src parameter: %q - %q", source[0], err.Error()))
		return
	}

	destinations, ok := r.URL.Query()["dst"]
	if !ok {
		handleError(w, http.StatusBadRequest, "at least one dst parameter must be specified")
		return
	}

	parsedDestinations := make([]coordinates, len(destinations))
	for idx, dst := range destinations {
		dstLong, dstLat, err := parseAndValidateLongitudeAndLatitude(dst)
		if err != nil {
			handleError(w, http.StatusBadRequest, fmt.Sprintf("invalid dst parameter: %q - %q", dst, err.Error()))
			return
		}

		parsedDestinations[idx] = coordinates{
			Longitude: dstLong,
			Latitude:  dstLat,
		}
	}

	response := &Response{
		Source: source[0],
	}

	for _, dst := range parsedDestinations {
		route, err := osrm.GetRoute(r.Context(), srcLong, srcLat, dst.Longitude, dst.Latitude)
		if err != nil {
			if errors.Is(err, osrm.ErrInvalidInput) {
				handleError(w, http.StatusBadRequest, err.Error())
				return
			}
			handleError(w, http.StatusInternalServerError, "unable to retrieve route for the given coordinates")
			return
		}
		response.Routes = append(response.Routes, route)
	}

	sort.Sort(ByDurationAndDistance(response.Routes))

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// Parses the coordinates string spliting by comman and validates if the lat and long are valid
func parseAndValidateLongitudeAndLatitude(input string) (float64, float64, error) {
	items := strings.Split(input, ",")

	if len(items) != 2 {
		return 0, 0, errors.New("invalid longitude and latitude - must be separeted by comma - long,lat")
	}

	long, err := strconv.ParseFloat(items[0], 64)
	if err != nil || long < -180 || long > 180 {
		return 0, 0, errors.New("invalid longitude - must be a float number between -180 and 180")
	}

	lat, err := strconv.ParseFloat(items[1], 64)
	if err != nil || lat < -90 || lat > 90 {
		return 0, 0, errors.New("invalid latitude - must be a float number between -90 and 90")
	}

	return long, lat, nil
}

// prepare the error response on the json format
func handleError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(message)
}
