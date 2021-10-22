package main

import (
	"context"
	"fmt"
	"go.uber.org/multierr"
	"net/url"
)

// SearchResult represents a result of Search.
type SearchResult struct {
	Returned  int      `json:"results_returned"`
	Available int      `json:"results_available"`
	Start     int      `json:"results_start"`
	Events    []*Event `json:"events"`
}

type SearchService interface {
	Search(ctx context.Context, params url.Values) (*SearchResult, error)
}

func SearchParam(params ...Param) (url.Values, error) {
	vals := make(url.Values)
	var err error
	for _, p := range params {
		err = multierr.Append(err, p(vals))
	}
	if err != nil {
		return nil, err
	}

	return vals, nil
}

type searchService struct {
	cli *Client
}

func (s *searchService) Search(ctx context.Context, params url.Values) (*SearchResult, error) {
	var r SearchResult
	if err := s.cli.get(ctx, "event", params, &r); err != nil {
		return nil, fmt.Errorf("connpass.Search: %w", err)
	}
	return &r, nil
}
