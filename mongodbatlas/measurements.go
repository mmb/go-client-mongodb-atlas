package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const measurementsPath = "groups/%s/processes/%s:%d/measurements?granularity=PT24H&period=PT24H"

type MeasurementsService interface {
	List(context.Context, string, string, int, *ListOptions) ([]Measurement, *Response, error)
}

type MeasurementsServiceOp struct {
	Client RequestDoer
}

var _ MeasurementsService = &MeasurementsServiceOp{}

type Datapoint struct {
	Value     float64 `json:"value"`
	Timestamp string  `json:"timestamp"`
}

type Measurement struct {
	Datapoints []Datapoint `json:"dataPoints"`
	Name       string      `json:"name"`
}

type measurementsResponse struct {
	Links        []*Link       `json:"links,omitempty"`
	Measurements []Measurement `json:"measurements,omitempty"`
}

func (s *MeasurementsServiceOp) List(ctx context.Context, groupID string, host string, port int, listOptions *ListOptions) ([]Measurement, *Response, error) {
	path := fmt.Sprintf(measurementsPath, groupID, host, port)

	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(measurementsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Measurements, resp, nil
}
