package bmhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	bm "github.com/Tsapen/bm/internal/bm"
	bs "github.com/Tsapen/bm/internal/book-service"
	"github.com/Tsapen/bm/pkg/api"
)

type Server struct {
	cfg Config

	tcpServer        *http.Server
	unixSocketServer *http.Server
}

type Config struct {
	Addr         string
	SocketPath   string
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
	r.HandleFunc("/books/{book_id}", handleFunc(parseGetBookReq, b.getBook)).Methods(http.MethodGet)
	r.HandleFunc("/books", handleFunc(parseGetBooksReq, b.getBooks)).Methods(http.MethodGet)
	r.HandleFunc("/books", handleFunc(parseJSONReq[api.CreateBookReq], b.createBook)).Methods(http.MethodPost)
	r.HandleFunc("/books/{book_id}", handleFunc(parseUpdateBookReq, b.updateBook)).Methods(http.MethodPut)
	r.HandleFunc("/books", handleFunc(parseJSONReq[api.DeleteBooksReq], b.deleteBooks)).Methods(http.MethodDelete)

	r.HandleFunc("/collections/{collection_id}", handleFunc(parseGetCollectionReq, b.getCollection)).Methods(http.MethodGet)
	r.HandleFunc("/collections", handleFunc(parseGetCollectionsReq, b.getCollections)).Methods(http.MethodGet)
	r.HandleFunc("/collections", handleFunc(parseJSONReq[api.CreateCollectionReq], b.createCollection)).Methods(http.MethodPost)
	r.HandleFunc("/collections/{collection_id}", handleFunc(parseUpdateCollectionReq, b.updateCollection)).Methods(http.MethodPut)
	r.HandleFunc("/collections/{collection_id}", handleFunc(parseDeleteCollectionReq, b.deleteCollection)).Methods(http.MethodDelete)

	r.HandleFunc("/collections/{collection_id}/books", handleFunc(parseCreateBooksCollectionReq, b.createBooksCollection)).Methods(http.MethodPost)
	r.HandleFunc("/collections/{collection_id}/books", handleFunc(parseDeleteBooksCollectionReq, b.deleteBooksCollection)).Methods(http.MethodDelete)

	var s = &Server{
		cfg: cfg,
		tcpServer: &http.Server{
			Addr:         cfg.Addr,
			Handler:      r,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
		},
		unixSocketServer: &http.Server{
			Handler:      r,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
		},
	}

	return s, nil
}

func handleFunc[Req any](
	parseReq func(r *http.Request) (Req, error),
	handle func(context.Context, Req) (any, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := bm.WithReqID(r.Context(), uuid.NewString())

		logger := log.With().Str("method", r.Method).Str("path", r.URL.String()).Str("request_id", bm.ReqIDFromCtx(ctx)).Logger()
		logger.Info().Msg("received request")

		req, err := parseReq(r)
		if err != nil {
			renderErr(ctx, logger, fmt.Errorf("parse request: %w", err), w)

			return
		}

		log.Info().Any("request", req).Msg("request is parsed")

		resp, err := handle(ctx, req)
		if err != nil {
			renderErr(ctx, logger, fmt.Errorf("handle request: %w", err), w)

			return
		}

		renderResponse(ctx, logger, resp, w)
	}
}

// Start runs server with tcp as transport.
func (s *Server) StartTCPServer() error {
	log.Info().Msgf("HTTP server (tcp) started to listen %s", s.cfg.Addr)

	return s.tcpServer.ListenAndServe()
}

// Start runs server with unix-socket as transport.
func (s *Server) StartUnixSocketServer() error {
	log.Info().Msgf("HTTP server (unix-socket) started to listen %s", s.cfg.Addr)

	unixListener, err := net.Listen("unix", s.cfg.SocketPath)
	if err != nil {
		return err
	}

	return s.unixSocketServer.Serve(unixListener)
}

func parseJSONReq[Req any](r *http.Request) (*Req, error) {
	req := new(Req)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, bm.NewValidationError("parse request: %w", err)
	}

	return req, nil
}

func renderErr(ctx context.Context, logger zerolog.Logger, err error, w http.ResponseWriter) {
	statusCode := httpStatus(err)
	w.WriteHeader(statusCode)

	logger.Info().Err(err).Int("status code", statusCode).Msg("failed to process message")
	renderResponse(ctx, logger, map[string]any{"error": err.Error()}, w)
}

func renderResponse(ctx context.Context, logger zerolog.Logger, resp any, w http.ResponseWriter) {
	if resp == nil {
		logger.Info().Msg("finish processing")

		return
	}

	logger.Info().Any("response", resp).Msg("finish processing")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Info().Err(err).Msg("send message")
	}
}
