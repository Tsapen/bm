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
		{name: "test create book validation", testFunc: s.testCreateBookValidation},
		{name: "test get books validation", testFunc: s.testGetBookValidation},
		{name: "test update book validation", testFunc: s.testUpdateBookValidation},
		{name: "test delete book validation", testFunc: s.testDeleteBookValidation},

		{name: "test collections CRUD", testFunc: s.testCollections},
		{name: "test create collection validation", testFunc: s.testCreateCollectionValidation},
		{name: "test get colleciton validation", testFunc: s.testGetCollectionValidation},
		{name: "test update collection validation", testFunc: s.testUpdateCollectionValidation},
		{name: "test delete collection validation", testFunc: s.testDeleteCollectionValidation},

		{name: "test books collection CRUD", testFunc: s.testBooksCollection},
		{name: "test create books collection validation", testFunc: s.testCreateBooksCollectionValidation},
		{name: "test remove books collection validation", testFunc: s.testDeleteBooksCollectionValidation},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.testFunc(ctx, t, client)
		})
	}
}
