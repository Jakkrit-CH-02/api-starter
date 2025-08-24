package user

type handler struct {
	service Service
}

func NewHandler(service Service) handler {
	return handler{service}
}