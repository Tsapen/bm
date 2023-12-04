package api

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	GetBookReq struct {
		ID int64 `json:"-"`
	}

	GetBookResp struct {
		Book Book `json:"book"`
	}

	GetBooksReq struct {
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
		PublishedDate time.Time `json:"published_date" layout:"2006-01-02"`
		Edition       string    `json:"edition"`
		Description   string    `json:"description"`
		Genre         string    `json:"genre"`
	}

	CreateBookResp struct {
		ID int64 `json:"id"`
	}

	UpdateBookReq struct {
		ID            int64     `json:"-"`
		Title         string    `json:"title"`
		Author        string    `json:"author"`
		PublishedDate time.Time `json:"published_date"`
		Edition       string    `json:"edition"`
		Description   string    `json:"description"`
		Genre         string    `json:"genre"`
	}

	UpdateBookResp struct {
		Success bool `json:"-"`
	}

	DeleteBooksReq struct {
		IDs []int64 `json:"ids"`
	}

	DeleteBooksResp struct {
		Success bool `json:"-"`
	}

	GetCollectionReq struct {
		ID int64 `json:"-"`
	}

	GetCollectionResp struct {
		Collection Collection `json:"collection"`
	}

	GetCollectionsReq struct {
		OrderBy  string `url:"order_by,omitempty" json:"order_by"`
		Desc     bool   `url:"desc,omitempty" json:"desc"`
		Page     int64  `url:"page,omitempty" json:"page"`
		PageSize int64  `url:"page_size,omitempty" json:"page_size"`
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
		ID          int64  `json:"-"`
		Name        string `json:"name"`
		Description string `json:"decription"`
	}

	UpdateCollectionResp struct {
		Success bool `json:"-"`
	}

	DeleteCollectionReq struct {
		ID int64 `json:"-"`
	}

	DeleteCollectionResp struct {
		Success bool `json:"-"`
	}

	CreateBooksCollectionReq struct {
		CID     int64   `json:"-"`
		BookIDs []int64 `json:"book_ids"`
	}

	CreateBooksCollectionResp struct {
		Success bool `json:"-"`
	}

	DeleteBooksCollectionReq struct {
		CID     int64   `json:"-"`
		BookIDs []int64 `json:"books_ids"`
	}

	DeleteBooksCollectionResp struct {
		Success bool `json:"-"`
	}
)

func (c *CreateBookReq) UnmarshalJSON(data []byte) error {
	var aux struct {
		Title         string `json:"title"`
		Author        string `json:"author"`
		PublishedDate string `json:"published_date"`
		Edition       string `json:"edition"`
		Description   string `json:"description"`
		Genre         string `json:"genre"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.Title = aux.Title
	c.Author = aux.Author
	c.Edition = aux.Edition
	c.Description = aux.Description
	c.Genre = aux.Genre

	if aux.PublishedDate != "" {
		parsedDate, err := time.Parse("2006-01-02", aux.PublishedDate)
		if err != nil {
			return fmt.Errorf("parse published_date: %w", err)
		}

		c.PublishedDate = parsedDate
	}

	return nil
}

func (c *CreateBookReq) MarshalJSON() ([]byte, error) {
	aux := struct {
		Title         string `json:"title"`
		Author        string `json:"author"`
		PublishedDate string `json:"published_date,omitempty"`
		Edition       string `json:"edition"`
		Description   string `json:"description"`
		Genre         string `json:"genre"`
	}{
		Title:       c.Title,
		Author:      c.Author,
		Edition:     c.Edition,
		Description: c.Description,
		Genre:       c.Genre,
	}

	if !c.PublishedDate.IsZero() {
		aux.PublishedDate = c.PublishedDate.Format("2006-01-02")
	}

	return json.Marshal(&aux)
}

func (u *UpdateBookReq) UnmarshalJSON(data []byte) error {
	var aux struct {
		Title         string `json:"title"`
		Author        string `json:"author"`
		PublishedDate string `json:"published_date"`
		Edition       string `json:"edition"`
		Description   string `json:"description"`
		Genre         string `json:"genre"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	u.Title = aux.Title
	u.Author = aux.Author
	u.Edition = aux.Edition
	u.Description = aux.Description
	u.Genre = aux.Genre

	if aux.PublishedDate != "" {
		parsedDate, err := time.Parse("2006-01-02", aux.PublishedDate)
		if err != nil {
			return fmt.Errorf("parse published_date: %w", err)
		}

		u.PublishedDate = parsedDate
	}

	return nil
}

func (u *UpdateBookReq) MarshalJSON() ([]byte, error) {
	aux := struct {
		Title         string `json:"title"`
		Author        string `json:"author"`
		PublishedDate string `json:"published_date"`
		Edition       string `json:"edition"`
		Description   string `json:"description"`
		Genre         string `json:"genre"`
	}{
		Title:       u.Title,
		Author:      u.Author,
		Edition:     u.Edition,
		Description: u.Description,
		Genre:       u.Genre,
	}

	if !u.PublishedDate.IsZero() {
		aux.PublishedDate = u.PublishedDate.Format("2006-01-02")
	}

	return json.Marshal(&aux)
}
