package usecase

import (
	"encoding/json"
	"fmt"
)

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

func (p *payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Body string `json:"body"`
	}{
		string(p.t),
		fmt.Sprintf("%v", p.body),
	})
}

func (p *payload) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func NewWillExecutePayload(arg interface{}) Payload {
	return &payload{
		t:    PayloadTypeWillExecute,
		body: arg,
	}
}

func NewDidExecutePayload(arg interface{}) Payload {
	return &payload{
		t:    PayloadTypeDidExecute,
		body: arg,
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
