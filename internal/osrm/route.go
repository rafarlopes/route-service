package osrm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidInput is used for the cases where the resquest was done successfully but the parameters are not good.
	ErrInvalidInput      error = errors.New("the input provided to the OSRM API is invalid - check longitude and latitude values")
	errOSRMRequestFailed error = errors.New("the request to OSRM API did not succeed")
)

type (
	// Route struct is used to hold the output from OSRM API containing the distance and the duration
	Route struct {
		Destination string  `json:"destination"`
		Distance    float64 `json:"distance"`
		Duration    float64 `json:"duration"`
	}

	// Internal response struct used to parse the json from the OSRM API response
	response struct {
		Code    string
		Message string
		Routes  []*Route
	}
)

// GetRoute requests the route to the OSRM API and return it using Route struct
func GetRoute(ctx context.Context, srcLong float64, srcLat float64, dstLong float64, dstLat float64) (*Route, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f", srcLong, srcLat, dstLong, dstLat), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create OSRM API route request")
	}

	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("overview", "false")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request route from OSRM API")
	}

	defer res.Body.Close()

	decodedResp := &response{}
	if err = json.NewDecoder(res.Body).Decode(decodedResp); err != nil {
		return nil, errors.Wrap(err, "failed to read route response from OSRM API")
	}

	switch decodedResp.Code {
	case "Ok":
		route := decodedResp.Routes[0]
		route.Destination = fmt.Sprintf("%f,%f", dstLong, dstLat)
		return route, nil
	case "InvalidValue":
		return nil, ErrInvalidInput
	default:
		return nil, errors.Wrapf(errOSRMRequestFailed, "code %q message %q", decodedResp.Code, decodedResp.Message)
	}

}
