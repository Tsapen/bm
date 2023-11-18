package bmunixsocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/Tsapen/bm/internal/bm"
	bs "github.com/Tsapen/bm/internal/book-service"
)

type (
	command struct {
		Action string          `json:"action"`
		Data   json.RawMessage `json:"data"`
	}

	serviceBundle struct {
		bookService *bs.Service
	}

	Config struct {
		SocketPath   string
		ConnMaxCount int
		Timeout      time.Duration
	}

	Server struct {
		cfg      Config
		listener net.Listener
		bundle   *serviceBundle

		sem chan struct{}
	}
)

// NewServer creates new unix-socket server.
func NewServer(cfg Config, bookService *bs.Service) (*Server, error) {
	listener, err := net.Listen("unix", cfg.SocketPath)
	if err != nil {
		return nil, fmt.Errorf("make tcp listener: %w", err)
	}

	sem := make(chan struct{}, cfg.ConnMaxCount)
	for i := 0; i < cfg.ConnMaxCount; i++ {
		sem <- struct{}{}
	}

	return &Server{
		listener: listener,
		bundle: &serviceBundle{
			bookService: bookService,
		},

		sem: sem,
	}, nil
}

func (s *Server) Start() {
	defer s.listener.Close()

	log.Info().Msgf("Unix socket server started to listen %s", s.cfg.SocketPath)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Warn().Err(err).Msgf("failed to accept connection")

			return
		}

		go func() {
			<-s.sem

			conn.SetReadDeadline(time.Now().Add(s.cfg.Timeout))
			conn.SetWriteDeadline(time.Now().Add(s.cfg.Timeout))

			ctx := bm.WithReqID(context.Background(), uuid.NewString())

			s.handle(ctx, conn)

			s.sem <- struct{}{}
		}()
	}
}
func (s *Server) handle(ctx context.Context, conn net.Conn) {
	logger := log.With().Str("request_id", bm.ReqIDFromCtx(ctx)).Logger()
	cmd := new(command)
	err := json.NewEncoder(conn).Encode(cmd)
	if err != nil {
		logger.Info().Err(err).Any("request", cmd).Msgf("parse command")

		return
	}

	var resp any
	data := cmd.Data
	switch cmd.Action {
	case "get_books":
		resp, err = handleFunc[getBooksReq](s.bundle.getBooks)(ctx, data)

	case "create_book":
		resp, err = handleFunc[createBookReq](s.bundle.createBook)(ctx, data)

	case "update_book":
		resp, err = handleFunc[updateBookReq](s.bundle.updateBook)(ctx, data)

	case "delete_books":
		resp, err = handleFunc[deleteBooksReq](s.bundle.deleteBooks)(ctx, data)

	case "get_collections":
		resp, err = handleFunc[getCollectionsReq](s.bundle.getCollections)(ctx, data)

	case "create_collection":
		resp, err = handleFunc[createCollectionReq](s.bundle.createCollection)(ctx, data)

	case "update_collection":
		resp, err = handleFunc[updateCollectionReq](s.bundle.updateCollection)(ctx, data)

	case "delete_collection":
		resp, err = handleFunc[deleteCollectionReq](s.bundle.deleteCollection)(ctx, data)

	case "create_books_collection":
		resp, err = handleFunc[createBooksCollectionReq](s.bundle.createBooksCollection)(ctx, data)

	case "delete_books_collection":
		resp, err = handleFunc[deleteBooksCollectionReq](s.bundle.deleteBooksCollection)(ctx, data)

	default:
		logger.Info().Err(err).Msgf("command not found")

		return
	}

	if err != nil {
		logger.Info().Err(err).Any("request", cmd).Msgf("handle request")

		return
	}

	if err = json.NewEncoder(conn).Encode(resp); err != nil {
		logger.Info().Err(err).Msgf("encode request")

		return
	}
}

func handleFunc[Req any](
	handle func(context.Context, *Req) (any, error),
) func(context.Context, json.RawMessage) (any, error) {
	return func(ctx context.Context, rawReq json.RawMessage) (any, error) {
		req := new(Req)
		if err := json.Unmarshal(rawReq, req); err != nil {
			return nil, fmt.Errorf("parse request: %w", err)
		}

		resp, err := handle(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("handle request: %w", err)
		}

		return resp, nil
	}
}
