package bmtest

import (
	"context"
	"testing"

	httpclient "github.com/Tsapen/bm/pkg/http-client"
)

type tc struct {
	name     string
	testFunc func(ctx context.Context, t *testing.T, client *httpclient.Client)
}

// TestBM does integration testing.
func TestBM(t *testing.T, client *httpclient.Client) {
	ctx := context.Background()

	s := new(storage)
	testcases := []tc{
		{name: "test books CRUD", testFunc: s.testBooks},
		{name: "test collections CRUD", testFunc: s.testCollections},
		{name: "test books collection CRUD", testFunc: s.testBooksCollection},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.testFunc(ctx, t, client)
		})
	}
}
