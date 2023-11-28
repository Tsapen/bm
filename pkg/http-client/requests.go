package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/google/go-querystring/query"

	bm "github.com/Tsapen/bm/internal/bm"
	"github.com/Tsapen/bm/pkg/api"
)

const (
	booksPath            = "/api/v1/books"
	bookPath             = "/api/v1/book"
	updateBookPath       = "/api/v1/book/update"
	collectionsPath      = "/api/v1/collections"
	collectionPath       = "/api/v1/collection"
	collectionUpdatePath = "/api/v1/collection/update"
	booksCollectionPath  = "/api/v1/collection/books"
)

func (c *Client) GetBooks(ctx context.Context, req *api.GetBooksReq) (*api.GetBooksResp, error) {
	resp := new(api.GetBooksResp)
	err := c.doRequestWithURLParams(ctx, booksPath, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) CreateBook(ctx context.Context, req *api.CreateBookReq) (*api.CreateBookResp, error) {
	resp := new(api.CreateBookResp)
	err := c.doRequestWithJSON(ctx, bookPath, http.MethodPost, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) UpdateBook(ctx context.Context, req *api.UpdateBookReq) (*api.UpdateBookResp, error) {
	resp := new(api.UpdateBookResp)
	err := c.doRequestWithJSON(ctx, updateBookPath, http.MethodPost, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) DeleteBooks(ctx context.Context, req *api.DeleteBooksReq) (*api.DeleteBooksResp, error) {
	resp := new(api.DeleteBooksResp)
	err := c.doRequestWithJSON(ctx, booksPath, http.MethodDelete, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetCollections(ctx context.Context, req *api.GetCollectionsReq) (*api.GetCollectionsResp, error) {
	resp := new(api.GetCollectionsResp)
	err := c.doRequestWithURLParams(ctx, collectionsPath, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) CreateCollection(ctx context.Context, req *api.CreateCollectionReq) (*api.CreateCollectionResp, error) {
	resp := new(api.CreateCollectionResp)
	err := c.doRequestWithJSON(ctx, collectionPath, http.MethodPost, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) UpdateCollection(ctx context.Context, req *api.UpdateCollectionReq) (*api.UpdateCollectionResp, error) {
	resp := new(api.UpdateCollectionResp)
	err := c.doRequestWithJSON(ctx, collectionUpdatePath, http.MethodPost, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) DeleteCollection(ctx context.Context, req *api.DeleteCollectionReq) (*api.DeleteCollectionResp, error) {
	resp := new(api.DeleteCollectionResp)
	err := c.doRequestWithJSON(ctx, collectionPath, http.MethodDelete, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) CreateBooksCollection(ctx context.Context, req *api.CreateBooksCollectionReq) (*api.CreateBooksCollectionResp, error) {
	resp := new(api.CreateBooksCollectionResp)
	err := c.doRequestWithJSON(ctx, booksCollectionPath, http.MethodPost, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) DeleteBooksCollection(ctx context.Context, req *api.DeleteBooksCollectionReq) (*api.DeleteBooksCollectionResp, error) {
	resp := new(api.DeleteBooksCollectionResp)
	err := c.doRequestWithJSON(ctx, booksCollectionPath, http.MethodDelete, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) doRequestWithJSON(ctx context.Context, urlPath, method string, reqData, respData any) (err error) {
	body, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	u, err := url.Parse(c.cfg.Address)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	u.Path = path.Join(u.Path, urlPath)

	req, err := http.NewRequest(method, u.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("construct request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	defer func() {
		err = bm.HandleErrPair(resp.Body.Close(), err)
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get error http status: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(respData); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (c *Client) doRequestWithURLParams(ctx context.Context, urlPath string, reqData, respData any) (err error) {
	vals, err := query.Values(reqData)
	if err != nil {
		return fmt.Errorf("construct request: %w", err)
	}

	u, err := url.Parse(c.cfg.Address)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	u.Path = path.Join(u.Path, urlPath)

	u.RawQuery = vals.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("construct request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	defer func() {
		err = bm.HandleErrPair(resp.Body.Close(), err)
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get error http status: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(respData); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
