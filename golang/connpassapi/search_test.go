package main

import (
	"context"
	"fmt"
	"testing"
)

func TestSearchParam(t *testing.T) {
	cli := NewClient()
	params, err := SearchParam(Keyword("rust"))
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	r, err := cli.Search(ctx, params)
	if err != nil {
		t.Fatal(err)
	}

	for _, e := range r.Events {
		t.Log(e.Title)
		fmt.Println(e.Title)
	}
}
