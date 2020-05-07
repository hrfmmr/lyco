package usecase

type UseCase interface {
	Execute(arg interface{}) error
}
