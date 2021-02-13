package osrm

import (
	"context"
	"testing"
)

func TestSuccessfulSimpleRequest(t *testing.T) {
	ctx := context.Background()

	_, err := GetRoute(ctx, 13.388860, 52.517037, 13.397634, 52.52940)

	if err != nil {
		t.Error("GetRoute should not return an error")
	}
}

func TestInvalidInputShouldReturnErrInvalidInput(t *testing.T) {
	ctx := context.Background()

	_, err := GetRoute(ctx, 200, 200, 200, 200)

	if err != ErrInvalidInput {
		t.Errorf("GetRoute return %v error, got %v", ErrInvalidInput, err)
	}
}

func TestGetRoutesReturnsValidRouteDistanceAndDuration(t *testing.T) {
	ctx := context.Background()

	route, _ := GetRoute(ctx, 13.388860, 52.517037, 13.397634, 52.52940)

	if route.Distance == 0 && route.Duration == 0 {
		t.Error("GetRoute did not return proper distance and duration")
	}
}
