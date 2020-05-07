package usecase

type UseCaseExecutor interface {
	UseCase() UseCase
	Execute(arg interface{}) error
	OnWillExecute() <-chan Payload
	OnDidExecute() <-chan Payload
}

type useCaseExecutor struct {
	useCase         UseCase
	willExecuteChan chan Payload
	didExecuteChan  chan Payload
}

func NewUseCaseExecutor(u UseCase) UseCaseExecutor {
	willExecuteChan := make(chan Payload)
	didExecuteChan := make(chan Payload)
	return &useCaseExecutor{u, willExecuteChan, didExecuteChan}
}

func (e *useCaseExecutor) UseCase() UseCase {
	return e.useCase
}

func (e *useCaseExecutor) Execute(arg interface{}) error {
	e.willExecuteChan <- NewWillExecutePayload()
	err := e.useCase.Execute(arg)
	if err != nil {
		return err
	}
	e.didExecuteChan <- NewDidExecutePayload()
	return nil
}

func (e *useCaseExecutor) OnWillExecute() <-chan Payload {
	return e.willExecuteChan
}

func (e *useCaseExecutor) OnDidExecute() <-chan Payload {
	return e.didExecuteChan
}
