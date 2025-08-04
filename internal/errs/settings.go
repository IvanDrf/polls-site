package errs

// Can't load config
func ErrCFGLoad() error {
	return Error{Code: -1, Msg: "can't load cfg file .env"}
}

// Can't connect to database
func ErrDBConnection() error {
	return Error{Code: -1, Msg: "can't connect to database"}
}

// Invalid logger level in database
func ErrLoggerLevel() error {
	return Error{Code: -1, Msg: "can't set up logger's level"}
}

// Error of starting server, cant start new server
func ErrCantStartServer() error {
	return Error{Code: -1, Msg: "can't start server"}
}
