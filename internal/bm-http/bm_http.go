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
	"github.com/Tsapen/bm/pkg/api"
)

type Server struct {
	cfg Config

	httpServer *http.Server
}

type Config struct {
	Addr         string
	ConnMaxCount int
	Timeout      time.Duration
}

type serviceBundle struct {
	bookService *bs.Service
}

func NewServer(cfg Config, bookService *bs.Service) (*Server, error) {
	b := &serviceBundle{
		bookService: bookService,
	}

	r := mux.NewRouter()
	r = r.PathPrefix("/api/v1").Subrouter()
	r.HandleFunc("/books", handleFunc(parseGetBooksReq, b.getBooks)).Methods(http.MethodGet)
	r.HandleFunc("/book", handleFunc(parseJSONReq[api.CreateBookReq], b.createBook)).Methods(http.MethodPost)
	r.HandleFunc("/book/update", handleFunc(parseJSONReq[api.UpdateBookReq], b.updateBook)).Methods(http.MethodPost)
	r.HandleFunc("/books", handleFunc(parseJSONReq[api.DeleteBooksReq], b.deleteBooks)).Methods(http.MethodDelete)

	r.HandleFunc("/collections", handleFunc(parseGetCollectionsReq, b.getCollections)).Methods(http.MethodGet)
	r.HandleFunc("/collection", handleFunc(parseJSONReq[api.CreateCollectionReq], b.createCollection)).Methods(http.MethodPost)
	r.HandleFunc("/collection/update", handleFunc(parseJSONReq[api.UpdateCollectionReq], b.updateCollection)).Methods(http.MethodPost)
	r.HandleFunc("/collection", handleFunc(parseJSONReq[api.DeleteCollectionReq], b.deleteCollection)).Methods(http.MethodDelete)

	r.HandleFunc("/collection/books", handleFunc(parseJSONReq[api.CreateBooksCollectionReq], b.createBooksCollection)).Methods(http.MethodPost)
	r.HandleFunc("/collection/books", handleFunc(parseJSONReq[api.DeleteBooksCollectionReq], b.deleteBooksCollection)).Methods(http.MethodDelete)

	var s = &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:         cfg.Addr,
			Handler:      r,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
		},
	}

	return s, nil
}

func handleFunc[Req, Resp any](
	parseReq func(r *http.Request) (Req, error),
	handle func(context.Context, Req) (Resp, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := bm.WithReqID(r.Context(), uuid.NewString())
		log.Info().Str("request_id", bm.ReqIDFromCtx(ctx)).Str("method", r.Method).Any("path", r.URL.String()).Msg("request")

		req, err := parseReq(r)
		if err != nil {
			renderErr(ctx, fmt.Errorf("parse request: %w", err), w)

			return
		}

		log.Info().Str("request_id", bm.ReqIDFromCtx(ctx)).Str("method", r.Method).Any("path", r.URL.String()).Any("request", req).Msg("start processing")

		resp, err := handle(ctx, req)
		if err != nil {
			renderErr(ctx, fmt.Errorf("handle request: %w", err), w)

			return
		}

		log.Info().Str("request_id", bm.ReqIDFromCtx(ctx)).Str("method", r.Method).Any("path", r.URL.String()).Any("response", resp).Msg("finish processing")

		renderResponse(ctx, resp, w)
	}
}

// Start runs server.
func (s *Server) Start() error {
	log.Info().Msgf("HTTP server started to listen %s", s.cfg.Addr)

	return s.httpServer.ListenAndServe()
}

func parseJSONReq[Req any](r *http.Request) (*Req, error) {
	req := new(Req)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
	}

	return req, nil
}

func renderErr(ctx context.Context, err error, w http.ResponseWriter) {
	log.Info().Str("request_id", bm.ReqIDFromCtx(ctx)).Err(err).Msg("process message")
	w.WriteHeader(httpStatus(err))
	renderResponse(ctx, err, w)
}

func renderResponse(ctx context.Context, resp any, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Info().Str("request_id", bm.ReqIDFromCtx(ctx)).Err(err).Msg("send message")
	}
}
