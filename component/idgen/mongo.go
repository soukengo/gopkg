package idgen

import "go.mongodb.org/mongo-driver/bson/primitive"

type mongoGenerator struct {
}

func NewMongoGenerator() Generator {
	return &mongoGenerator{}
}

func (g *mongoGenerator) GenID() (id string, err error) {
	id = primitive.NewObjectID().Hex()
	return
}
