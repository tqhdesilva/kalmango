module tqhdesilva/kalmango

// +heroku goVersion 1.16
go 1.16

require (
	github.com/gorilla/websocket v1.4.2
	github.com/spf13/viper v1.8.0
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.3.6 // indirect
	gonum.org/v1/gonum v0.9.2
)

// +heroku install ./cmd/api
