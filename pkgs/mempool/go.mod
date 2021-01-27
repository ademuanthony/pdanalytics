module github.com/planetdecred/pdanalytics/pkgs/mempool

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/DATA-DOG/go-sqlmock v1.3.3 // indirect
	github.com/decred/dcrd/chaincfg/chainhash v1.0.2
	github.com/decred/dcrd/chaincfg/v2 v2.3.0
	github.com/decred/dcrd/rpc/jsonrpc/types/v2 v2.3.0
	github.com/decred/dcrd/rpcclient/v5 v5.0.1
	github.com/decred/dcrd/wire v1.4.0
	github.com/decred/slog v1.1.0
	github.com/friendsofgo/errors v0.9.2
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/kat-co/vala v0.0.0-20170210184112-42e1d8b61f12
	github.com/kr/pretty v0.2.0 // indirect
	github.com/lib/pq v1.9.0
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/planetdecred/pdanalytics/dbhelpers v0.0.0-00010101000000-000000000000
	github.com/planetdecred/pdanalytics/pkgs/cache v0.0.0-00010101000000-000000000000
	github.com/planetdecred/pdanalytics/web v0.0.0-20210125191324-0735b483e313
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/volatiletech/null v8.0.0+incompatible
	github.com/volatiletech/sqlboiler v3.7.1+incompatible
	golang.org/x/net v0.0.0-20190923162816-aa69164e4478 // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace (
	github.com/planetdecred/pdanalytics/dbhelpers => ../../dbhelpers
	github.com/planetdecred/pdanalytics/pkgs/cache => ../cache
	github.com/planetdecred/pdanalytics/web => ../../web
)
