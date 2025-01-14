module github.com/TheRangiCrew/WITS/services/nws/awips

go 1.23.2

require (
	github.com/TheRangiCrew/go-nws v0.0.0-20241219224257-055802b1425f
	github.com/joho/godotenv v1.5.1
	github.com/spf13/cobra v1.8.1
	github.com/surrealdb/surrealdb.go v0.3.2
	github.com/xmppo/go-xmpp v0.2.9
)

require (
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.33.0 // indirect
)

// replace github.com/surrealdb/surrealdb.go => ../../../../../surrealdb.go

// replace github.com/TheRangiCrew/go-nws => ../../../../go-nws
