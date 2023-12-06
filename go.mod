module github.com/Tsapen/bm

go 1.21.4

require (
	github.com/Tsapen/bm/pkg/api v0.0.0-00010101000000-000000000000
	github.com/Tsapen/bm/pkg/http-client v0.0.0-00010101000000-000000000000
	github.com/caarlos0/env/v9 v9.0.0
	github.com/golang-migrate/migrate/v4 v4.16.2
	github.com/google/uuid v1.3.1
	github.com/gorilla/mux v1.8.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	github.com/rs/zerolog v1.29.0
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/cobra v1.8.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/Tsapen/bm/pkg/api => ./pkg/api
	github.com/Tsapen/bm/pkg/http-client => ./pkg/http-client
)
