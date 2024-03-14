package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rutas struct {
	Id                  primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Origen              string             `json:"origen,omitempty" bson:"origen" validate:"required"`
	Destino             string             `json:"destino,omitempty" bson:"destino" validate:"required"`
	AsientosDisponibles int                `json:"asientos_disponibles,omitempty" bson:"asientos_disponibles" validate:"required"`
	Kilometos           string             `json:"kilometros,omitempty" bson:"kilometros" validate:"required"`
	Fecha               time.Time          `json:"fecha,omitempty" bson:"fecha" validate:"required"`
	Precio              float64            `json:"precio,omitempty" bson:"precio" validate:"required"`
}
