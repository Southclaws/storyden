package authentication

type Service interface {
	//
}

type cookie struct{}

func New() Service {
	return &cookie{}
}
