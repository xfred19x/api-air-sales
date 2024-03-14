package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Ticket struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	ReservaId    string             `json:"reserva_id,omitempty" bson:"reserva_id" validate:"required"`
	PagoId       string             `json:"pago_id,omitempty" bson:"pago_id" validate:"required"`
	FechaEmision time.Time          `json:"fecha_emision,omitempty" bson:"fecha_emision" validate:"required"`
}
