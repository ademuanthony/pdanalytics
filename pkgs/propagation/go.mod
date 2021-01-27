module github.com/planetdecred/pdanalytics/pkgs/propagation

go 1.13

require (
	github.com/decred/dcrd/blockchain/stake v1.1.0
	github.com/decred/dcrd/chaincfg/chainhash v1.0.2
	github.com/decred/dcrd/chaincfg/v2 v2.3.0
	github.com/decred/dcrd/dcrutil v1.3.0
	github.com/decred/dcrd/rpc/jsonrpc/types/v2 v2.3.0
	github.com/decred/dcrd/rpcclient/v5 v5.0.1
	github.com/decred/dcrd/wire v1.3.0
	github.com/decred/slog v1.1.0
	github.com/friendsofgo/errors v0.9.2
	github.com/kat-co/vala v0.0.0-20170210184112-42e1d8b61f12
	github.com/lib/pq v1.9.0
	github.com/planetdecred/dcrextdata v0.0.0-20201006010145-cddd04eb454b
	github.com/planetdecred/pdanalytics/dbhelpers v0.0.0-00010101000000-000000000000
	github.com/planetdecred/pdanalytics/pkgs/cache v0.0.0-00010101000000-000000000000
	github.com/planetdecred/pdanalytics/web v0.0.0-20210125191324-0735b483e313
	github.com/spf13/viper v1.3.2
	github.com/volatiletech/null v8.0.0+incompatible
	github.com/volatiletech/sqlboiler v3.7.1+incompatible
	golang.org/x/text v0.3.3 // indirect
)

replace (
	github.com/planetdecred/pdanalytics/dbhelpers => ../../dbhelpers
	github.com/planetdecred/pdanalytics/pkgs/cache => ../cache
	github.com/planetdecred/pdanalytics/web => ../../web
)
