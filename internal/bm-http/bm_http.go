package bmhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	bm "github.com/Tsapen/bm/internal/bm"
	bs "github.com/Tsapen/bm/internal/book-service"
)

type Server struct {
	*http.Server
}

type Config struct {
	Addr         string
	ConnMaxCount int
	Timeout      time.Duration
}

type serviceBundle struct {
	bookService *bs.Service
	router      *mux.Router
}

func NewServer(cfg Config, bookService *bs.Service) (*Server, error) {
	r := mux.NewRouter()
	b := &serviceBundle{
		router:      r,
		bookService: bookService,
	}

	r = r.PathPrefix("/api/v1").Subrouter()
	r.HandleFunc("/books", handleFunc(parseGetBooksReq, b.getBooks)).
		Queries(
			"author", "{author}",
			"genre", "{genre}",
			"start_date", "{start_date}",
			"finish_date", "{finish_date}",
			"order_by", "{order_by}",
			"desc", "{desc}",
			"collection_id", "{collection_id}",
			"page", "{page}",
			"page_size", "{page_size}",
		).
		Methods(http.MethodGet)
	r.HandleFunc("/book", handleFunc(parseJSONReq[createBookReq], b.createBook)).Methods(http.MethodPost)
	r.HandleFunc("/book", handleFunc(parseJSONReq[updateBookReq], b.updateBook)).Methods(http.MethodPatch)
	r.HandleFunc("/books", handleFunc(parseJSONReq[deleteBooksReq], b.deleteBooks)).Methods(http.MethodDelete)

	r.HandleFunc("/collections", handleFunc(parseGetCollectionsReq, b.getCollections)).
		Queries(
			"ids", "{ids}",
			"order_by", "{order_by}",
			"desc", "{desc}",
			"page", "{page}",
			"page_size", "{page_size}",
		).
		Methods(http.MethodGet)
	r.HandleFunc("/collection", handleFunc(parseJSONReq[createCollectionReq], b.createCollection)).Methods(http.MethodPost)
	r.HandleFunc("/collection", handleFunc(parseJSONReq[updateCollectionReq], b.updateCollection)).Methods(http.MethodPatch)
	r.HandleFunc("/collection", handleFunc(parseJSONReq[deleteCollectionReq], b.deleteCollection)).Methods(http.MethodDelete)

	r.HandleFunc("/collection/books", handleFunc(parseJSONReq[createBooksCollectionReq], b.createBooksCollection)).Methods(http.MethodPost)
	r.HandleFunc("/collection/books", handleFunc(parseJSONReq[deleteBooksCollectionReq], b.deleteBooksCollection)).Methods(http.MethodDelete)

	var s = &Server{
		&http.Server{
			Addr:         cfg.Addr,
			Handler:      r,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
		},
	}

	return s, nil
}

type (
	handler[Req, Resp any] struct {
		handle   func(context.Context, Req) (Resp, error)
		parseReq func(r *http.Request) (Req, error)
	}
)

func newHandler[Req, Resp any](
	parseReq func(r *http.Request) (Req, error),
	handle func(context.Context, Req) (Resp, error),
) *handler[Req, Resp] {
	return &handler[Req, Resp]{
		parseReq: parseReq,
		handle:   handle,
	}
}

func handleFunc[Req, Resp any](
	parseReq func(r *http.Request) (Req, error),
	handle func(context.Context, Req) (Resp, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := bm.WithReqID(r.Context(), uuid.NewString())

		req, err := parseReq(r)
		if err != nil {
			renderErr(ctx, fmt.Errorf("parse request: %w", err), w)

			return
		}

		resp, err := handle(ctx, req)
		if err != nil {
			renderErr(ctx, fmt.Errorf("handle request: %w", err), w)

			return
		}

		renderResponse(ctx, resp, w)
	}
}

// Start runs server.
func (s *Server) Start() error {
	return s.ListenAndServe()
}

func parseJSONReq[Req any](r *http.Request) (*Req, error) {
	req := new(Req)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}

	return req, nil
}

func renderErr(ctx context.Context, err error, w http.ResponseWriter) {

}

func renderResponse(ctx context.Context, resp any, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Info().Str("request_id", bm.ReqIDFromCtx(ctx)).Err(err).Msg("send message")
	}
}
