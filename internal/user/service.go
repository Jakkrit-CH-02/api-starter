package user

type(
	Service interface {
	}
	service struct {
		repository Repository
	}
)

func NewService(repository Repository) service {
	return service{repository}
}