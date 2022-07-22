package authentication

type Service interface {
	//
}

func New() Service {
	return struct{}{}
}
