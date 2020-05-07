package usecase

type PayloadType string

const (
	PayloadTypeWillExecute PayloadType = "PAYLOAD_TYPE_WILL_EXECUTE"
	PayloadTypeDidExecute  PayloadType = "PAYLOAD_TYPE_DID_EXECUTE"
)

type Payload interface {
	Type() PayloadType
	Body() interface{}
}

type payload struct {
	t    PayloadType
	body interface{}
}

func (p *payload) Type() PayloadType {
	return p.t
}

func (p *payload) Body() interface{} {
	return p.body
}

func NewWillExecutePayload() Payload {
	return &payload{
		t: PayloadTypeWillExecute,
	}
}

func NewDidExecutePayload() Payload {
	return &payload{
		t: PayloadTypeDidExecute,
	}
}

type PayloadMeta interface {
	UseCase() UseCase
}

type payloadMeta struct {
	useCase UseCase
}

func NewPayLoadMeta(u UseCase) PayloadMeta {
	return &payloadMeta{u}
}

func (m *payloadMeta) UseCase() UseCase {
	return m.useCase
}
