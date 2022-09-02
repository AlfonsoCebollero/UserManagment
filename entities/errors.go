package entities

import "google.golang.org/grpc/status"

var (
	AlreadyRegisteredEmailError = status.Error(6, "email already registered")
)
