package telegram

import "reflect"

const PlaceIdEmpty PlaceId = "empty"

type PlaceStorage interface {
	Set(id UserId, place *Place)
	Get(id UserId) *Place
	Delete(id UserId)
}

type PlaceId string

func NewPlace(id PlaceId, context ...interface{}) *Place {
	return &Place{
		Id:      id,
		Context: context,
	}
}

type Place struct {
	Id      PlaceId
	Context interface{}
}

func (p *Place) Eq(place *Place) bool {
	if p.Id != place.Id {
		return false
	}

	return reflect.DeepEqual(p.Context, place.Context)
}
