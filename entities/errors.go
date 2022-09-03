package entities

import "google.golang.org/grpc/status"

var (
	InvalidEmailError           = status.Error(3, "entered email is not valid")
	AlreadyRegisteredEmailError = status.Error(6, "email already registered")
	NotFoundUser                = status.Error(5, "could not found user")
)
