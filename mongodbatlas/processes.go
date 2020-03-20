package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const processesPath = "groups/%s/processes"

type ProcessesService interface {
	List(context.Context, string, *ListOptions) ([]Process, *Response, error)
}

type ProcessesServiceOp struct {
	Client RequestDoer
}

var _ ProcessesService = &ProcessesServiceOp{}

type Process struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type processesResponse struct {
	Links      []*Link   `json:"links,omitempty"`
	Results    []Process `json:"results,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"`
}

func (s *ProcessesServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]Process, *Response, error) {
	path := fmt.Sprintf(processesPath, groupID)

	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(processesResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}
