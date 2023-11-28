package bmtest

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Tsapen/bm/pkg/api"
	httpclient "github.com/Tsapen/bm/pkg/http-client"
)

type storage struct {
	books       []*api.Book
	collections []*api.Collection
}

func (s *storage) testBooks(ctx context.Context, t *testing.T, client *httpclient.Client) {
	// 1. Create books.
	s.books = []*api.Book{
		{
			Title:         "The White Guard",
			Author:        "Mikhail Bulgakov",
			PublishedDate: time.Date(1925, time.January, 1, 0, 0, 0, 0, time.UTC),
			Edition:       "First Edition",
			Description:   "The White Guard is a novel by Mikhail Bulgakov, set in Kiev, Ukraine, during the Russian Civil War.",
			Genre:         "Historical Fiction",
		},
		{
			Title:         "A Hero of Our Time",
			Author:        "Mikhail Lermontov",
			PublishedDate: time.Date(1840, time.January, 1, 0, 0, 0, 0, time.UTC),
			Edition:       "First Edition",
			Description:   "A Hero of Our Time is a novel by Mikhail Lermontov, portraying the tragic and absurd life of a young officer in the Russian Caucasus.",
			Genre:         "Classic Fiction",
		},
		{
			Title:         "Chapayev and Void",
			Author:        "Victor Pelevin",
			PublishedDate: time.Date(1996, time.January, 1, 0, 0, 0, 0, time.UTC),
			Edition:       "First Edition",
			Description:   "Chapayev and Void is a satirical novel by Victor Pelevin, exploring the post-Soviet Russian society through a surreal and philosophical narrative.",
			Genre:         "Satire",
		},
		{
			Title:         "Journey to the End of the Night",
			Author:        "Louis-Ferdinand Céline",
			PublishedDate: time.Date(1932, time.January, 1, 0, 0, 0, 0, time.UTC),
			Edition:       "First Edition",
			Description:   "Journey to the End of the Night is a novel by Louis-Ferdinand Céline, providing a raw and intense depiction of the author's experiences as a young doctor in colonial Africa.",
			Genre:         "Modernist Literature",
		},
		{
			Title:         "The Possibility of an Island",
			Author:        "Michel Houellebecq",
			PublishedDate: time.Date(1994, time.January, 1, 0, 0, 0, 0, time.UTC),
			Edition:       "First Edition",
			Description:   "The Possibility of an Island is a novel by Michel Houellebecq, exploring themes of love, sexuality, and the impact of technology on human relationships.",
			Genre:         "Contemporary Fiction",
		},
	}
	for _, b := range s.books {
		req := &api.CreateBookReq{
			Title:         b.Title,
			Author:        b.Author,
			PublishedDate: b.PublishedDate,
			Edition:       b.Edition,
			Description:   b.Description,
			Genre:         b.Genre,
		}

		resp, err := client.CreateBook(ctx, req)
		assert.NoError(t, err)

		b.ID = resp.ID
	}

	// 2. Read testing.
	filterTests := []struct {
		give      *api.GetBooksReq
		wantBooks []api.Book
	}{
		{
			give: &api.GetBooksReq{
				ID: s.books[0].ID,
			},
			wantBooks: []api.Book{*s.books[0]},
		},
		{
			give: &api.GetBooksReq{
				Author: s.books[1].Author,
			},
			wantBooks: []api.Book{*s.books[1]},
		},
		{
			give: &api.GetBooksReq{
				Genre: s.books[2].Genre,
			},
			wantBooks: []api.Book{*s.books[2]},
		},
		{
			give: &api.GetBooksReq{
				StartDate:  time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC),
				FinishDate: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
				OrderBy:    "author",
				Desc:       true,
				Page:       2,
				PageSize:   2,
			},
			wantBooks: []api.Book{*s.books[4], *s.books[3]},
		},
	}
	for _, tt := range filterTests {
		got := getBooks(ctx, t, client, tt.give)
		assert.Equal(t, len(tt.wantBooks), len(got.Books))
		for i := 0; i < len(tt.wantBooks); i++ {
			assert.Equal(t, tt.wantBooks[i].ID, got.Books[i].ID)
			assert.Equal(t, tt.wantBooks[i].Title, got.Books[i].Title)
		}
	}

	// 3. Update testing.
	book := s.books[4]
	updateReq := api.UpdateBookReq{
		ID:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		PublishedDate: book.PublishedDate,
		Edition:       book.Edition,
		Description:   book.Description,
		Genre:         book.Genre,
	}
	updateTests := []struct {
		modify func(*api.UpdateBookReq)
	}{
		{
			modify: func(req *api.UpdateBookReq) {
				req.Title = uuid.NewString()
			},
		},
		{
			modify: func(req *api.UpdateBookReq) {
				req.Author = uuid.NewString()
			},
		},
		{
			modify: func(req *api.UpdateBookReq) {
				req.PublishedDate = time.Now()
			},
		},
		{
			modify: func(req *api.UpdateBookReq) {
				req.Edition = uuid.NewString()
			},
		},
		{
			modify: func(req *api.UpdateBookReq) {
				req.Description = uuid.NewString()
			},
		},
		{
			modify: func(req *api.UpdateBookReq) {
				req.Genre = uuid.NewString()
			},
		},
	}
	for _, tt := range updateTests {
		req := updateReq

		tt.modify(&req)

		resp, err := client.UpdateBook(ctx, &req)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		want := req
		got := getBooks(ctx, t, client, &api.GetBooksReq{ID: book.ID})
		assert.Equal(t, 1, len(got.Books))
		gotBook := got.Books[0]
		assert.Equal(t, want.Title, gotBook.Title)
		assert.Equal(t, want.Author, gotBook.Author)
		assert.Equal(t, want.PublishedDate.Truncate(24*time.Hour).String(), gotBook.PublishedDate.String())
		assert.Equal(t, want.Edition, gotBook.Edition)
		assert.Equal(t, want.Description, gotBook.Description)
		assert.Equal(t, want.Genre, gotBook.Genre)
	}

	// 4. Delete testing.
	deleteResp, err := client.DeleteBooks(ctx, &api.DeleteBooksReq{IDs: []int64{s.books[4].ID}})
	assert.NoError(t, err)
	assert.True(t, deleteResp.Success)

	s.books = s.books[:4]
}

func (s *storage) testCollections(ctx context.Context, t *testing.T, client *httpclient.Client) {
	// 1. Create collections.
	s.collections = []*api.Collection{
		{
			Name:        "Classic Novels",
			Description: "A collection of classic novels from various authors.",
		},
		{
			Name:        "Russian Literature",
			Description: "A collection of novels and literature from Russian authors.",
		},
		{
			Name:        "Contemporary Fiction",
			Description: "A collection of modern and contemporary fiction books.",
		},
	}
	for _, c := range s.collections {
		req := &api.CreateCollectionReq{
			Name:        c.Name,
			Description: c.Description,
		}
		resp, err := client.CreateCollection(ctx, req)
		assert.NoError(t, err)
		c.ID = resp.ID
	}

	// 2. Read testing.
	filterTests := []struct {
		give            *api.GetCollectionsReq
		wantCollections []api.Collection
	}{
		{
			give: &api.GetCollectionsReq{
				IDs: []int64{s.collections[0].ID},
			},
			wantCollections: []api.Collection{*s.collections[0]},
		},
		{
			give: &api.GetCollectionsReq{
				OrderBy:  "name",
				Desc:     true,
				Page:     2,
				PageSize: 1,
			},
			wantCollections: []api.Collection{*s.collections[2]},
		},
	}
	for _, tt := range filterTests {
		got := getCollections(ctx, t, client, tt.give)
		assert.Equal(t, len(tt.wantCollections), len(got.Collections))
		for i := 0; i < len(tt.wantCollections); i++ {
			assert.Equal(t, tt.wantCollections[i].ID, got.Collections[i].ID)
			assert.Equal(t, tt.wantCollections[i].Name, got.Collections[i].Name)
		}
	}

	// 3. Update testing.
	collection := s.collections[2]
	updateReq := api.UpdateCollectionReq{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
	}
	updateTests := []struct {
		modify func(*api.UpdateCollectionReq)
	}{
		{
			modify: func(req *api.UpdateCollectionReq) {
				req.Name = uuid.NewString()
			},
		},
		{
			modify: func(req *api.UpdateCollectionReq) {
				req.Description = uuid.NewString()
			},
		},
	}
	for _, tt := range updateTests {
		req := updateReq

		tt.modify(&req)

		resp, err := client.UpdateCollection(ctx, &req)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		want := req
		got := getCollections(ctx, t, client, &api.GetCollectionsReq{IDs: []int64{collection.ID}})
		assert.Equal(t, 1, len(got.Collections))
		gotCollection := got.Collections[0]
		assert.Equal(t, want.Name, gotCollection.Name)
		assert.Equal(t, want.Description, gotCollection.Description)
	}

	// 4. Delete testing.
	deleteResp, err := client.DeleteCollection(ctx, &api.DeleteCollectionReq{ID: s.collections[2].ID})
	assert.NoError(t, err)
	assert.True(t, deleteResp.Success)

	s.collections = s.collections[:2]
}

func (s *storage) testBooksCollection(ctx context.Context, t *testing.T, client *httpclient.Client) {
	// 1. Create books collection.
	collectionID := s.collections[0].ID
	createCollectionResp, err := client.CreateBooksCollection(ctx, &api.CreateBooksCollectionReq{
		CID:     s.collections[0].ID,
		BookIDs: []int64{s.books[0].ID, s.books[1].ID, s.books[2].ID, s.books[3].ID},
	})
	assert.NoError(t, err)
	assert.True(t, createCollectionResp.Success)

	got, err := client.GetBooks(ctx, &api.GetBooksReq{CollectionID: collectionID})
	assert.NoError(t, err)
	assert.Len(t, got.Books, 4)
	for i, b := range s.books[:4] {
		assert.Equal(t, b.ID, got.Books[i].ID)
		assert.Equal(t, b.Title, got.Books[i].Title)
	}

	// 2. Remove books collection.
	deleteCollectionResp, err := client.DeleteBooksCollection(ctx, &api.DeleteBooksCollectionReq{
		CID:     s.collections[0].ID,
		BookIDs: []int64{s.books[2].ID, s.books[3].ID},
	})
	assert.NoError(t, err)
	assert.True(t, deleteCollectionResp.Success)

	got, err = client.GetBooks(ctx, &api.GetBooksReq{CollectionID: collectionID})
	assert.NoError(t, err)
	assert.Len(t, got.Books, 2)
	for i, b := range s.books[:2] {
		assert.Equal(t, b.ID, got.Books[i].ID)
		assert.Equal(t, b.Title, got.Books[i].Title)
	}
}

func getBooks(ctx context.Context, t *testing.T, client *httpclient.Client, req *api.GetBooksReq) *api.GetBooksResp {
	resp, err := client.GetBooks(ctx, req)
	assert.NoError(t, err)

	return resp
}

func getCollections(ctx context.Context, t *testing.T, client *httpclient.Client, req *api.GetCollectionsReq) *api.GetCollectionsResp {
	resp, err := client.GetCollections(ctx, req)
	assert.NoError(t, err)

	return resp
}
