package api

import (
	"time"
)

type (
	GetBooksReq struct {
		ID           int64     `url:"id,omitempty" json:"id"`
		Author       string    `url:"author,omitempty" json:"author"`
		Genre        string    `url:"genre,omitempty" json:"genre"`
		CollectionID int64     `url:"collection_id,omitempty" json:"collection_id"`
		StartDate    time.Time `url:"start_date,omitempty" json:"start_date" layout:"2006-01-02"`
		FinishDate   time.Time `url:"finish_date,omitempty" json:"finish_date" layout:"2006-01-02"`
		OrderBy      string    `url:"order_by,omitempty" json:"order_by"`
		Desc         bool      `url:"desc,omitempty" json:"desc"`
		Page         int64     `url:"page,omitempty" json:"date"`
		PageSize     int64     `url:"page_size,omitempty" json:"page_size"`
	}

	GetBooksResp struct {
		Books []Book `json:"books"`
	}

	Book struct {
		ID            int64     `json:"id"`
		Title         string    `json:"title"`
		Author        string    `json:"author"`
		PublishedDate time.Time `json:"published_date"`
		Edition       string    `json:"edition"`
		Description   string    `json:"description"`
		Genre         string    `json:"genre"`
	}

	CreateBookReq struct {
		Title         string    `json:"title"`
		Author        string    `json:"author"`
		PublishedDate time.Time `json:"published_date"`
		Edition       string    `json:"edition"`
		Description   string    `json:"description"`
		Genre         string    `json:"genre"`
	}

	CreateBookResp struct {
		ID int64 `json:"id"`
	}

	UpdateBookReq struct {
		ID            int64     `json:"id"`
		Title         string    `json:"title"`
		Author        string    `json:"author"`
		PublishedDate time.Time `json:"published_date"`
		Edition       string    `json:"edition"`
		Description   string    `json:"description"`
		Genre         string    `json:"genre"`
	}

	UpdateBookResp struct {
		Success bool `json:"success"`
	}

	DeleteBooksReq struct {
		IDs []int64 `json:"ids"`
	}

	DeleteBooksResp struct {
		Success bool `json:"success"`
	}

	GetCollectionsReq struct {
		IDs      []int64 `url:"ids,omitempty"`
		OrderBy  string  `url:"order_by,omitempty"`
		Desc     bool    `url:"desc,omitempty"`
		Page     int64   `url:"page,omitempty"`
		PageSize int64   `url:"page_size,omitempty"`
	}

	Collection struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"decription"`
	}

	GetCollectionsResp struct {
		Collections []Collection `json:"collections"`
	}

	CreateCollectionReq struct {
		Name        string `json:"name"`
		Description string `json:"decription"`
	}

	CreateCollectionResp struct {
		ID int64 `json:"id"`
	}

	UpdateCollectionReq struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"decription"`
	}

	UpdateCollectionResp struct {
		Success bool `json:"success"`
	}

	DeleteCollectionReq struct {
		ID int64 `json:"id"`
	}

	DeleteCollectionResp struct {
		Success bool `json:"success"`
	}

	CreateBooksCollectionReq struct {
		CID     int64   `json:"collection_id"`
		BookIDs []int64 `json:"books_ids"`
	}

	CreateBooksCollectionResp struct {
		Success bool `json:"success"`
	}

	DeleteBooksCollectionReq struct {
		CID     int64   `json:"collection_id"`
		BookIDs []int64 `json:"books_ids"`
	}

	DeleteBooksCollectionResp struct {
		Success bool `json:"success"`
	}
)
