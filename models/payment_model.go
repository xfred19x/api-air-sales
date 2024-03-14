package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Pagos struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	ReservaId  string             `json:"reserva_id,omitempty" bson:"reserva_id" validate:"required"`
	MetodoPago string             `json:"metodo_pago,omitempty" bson:"metodo_pago" validate:"required"`
	Fecha      time.Time          `json:"fecha,omitempty" bson:"fecha" validate:"required"`
	Monto      float64            `json:"monto,omitempty" bson:"monto" validate:"required"`
	Estado     string             `json:"estado,omitempty" bson:"estado" validate:"required"`
}
