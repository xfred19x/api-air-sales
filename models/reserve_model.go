package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Reservas struct {
	Id              primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	RutaId          string             `json:"ruta_id,omitempty" bson:"ruta_id" validate:"required"`
	NombreUsuario   string             `json:"nombre_usuario,omitempty" bson:"nombre_usuario" validate:"required"`
	AsientosReserva int                `json:"asientos_reserva,omitempty" bson:"asientos_reserva" validate:"required"`
	Equipaje        bool               `json:"equipaje,omitempty" bson:"equipaje" validate:"required"`
	Fecha           time.Time          `json:"fecha,omitempty" bson:"fecha" validate:"required"`
	Estado          string             `json:"estado,omitempty" bson:"estado" validate:"required"`
}
