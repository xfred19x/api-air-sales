package controllers

import (
	"api-air-sales/configs"
	"api-air-sales/models"
	"api-air-sales/responses"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var reserveCollection *mongo.Collection = configs.GetCollection(configs.DB, "reserve")
var reserveValidate = validator.New()

// CreateReservation : Crea una reservacion
func CreateReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// valida el cuerpo del request
		var reserva models.Reservas
		if err := c.BindJSON(&reserva); err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Valida si ya hubo antes una reserva
		if existe, err := ValidateReservation(ctx, reserva.RutaId, reserva.NombreUsuario, reserva.Fecha); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		} else if existe {
			c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Ya se registro esta reserva"}})
			return
		}

		// Valida Campos Obligatorios
		if validationErr := reserveValidate.Struct(&reserva); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		// Actualiza el numero de asientos disponibles de la ruta
		if actualiza, err := UpdateAvailableSeat(ctx, reserva.RutaId, reserva.AsientosReserva, true); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		} else if !actualiza {
			c.JSON(http.StatusCreated, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "No se encontro la ruta para actualizar su asiento disponible"}})
			return
		}

		// Crea la Reserva
		reserva.Id = primitive.NewObjectID()
		if result, err := reserveCollection.InsertOne(ctx, reserva); err != nil {

			// Realiza rollbak en caso la Reserva falle, actualizando el numero de asientos disponibles como estaba inicialmente
			_, _ = UpdateAvailableSeat(ctx, reserva.RutaId, reserva.AsientosReserva, false)

			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			c.JSON(http.StatusCreated, responses.ErrorGenericoResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
		}

	}
}

// ValidateReservation : Valida si Existe una reservacion
func ValidateReservation(ctx context.Context, rutaID string, nombreUsuario string, fecha time.Time) (bool, error) {
	filter := bson.M{"ruta_id": rutaID, "nombre_usuario": nombreUsuario, "fecha": fecha}
	count, err := reserveCollection.CountDocuments(ctx, filter)
	return count != 0, err
}

// GetReserveId : Obtiene la reserva con el ID
func GetReserveId(ctx context.Context, reserveID string) (models.Reservas, error) {
	var reserve models.Reservas
	objId, _ := primitive.ObjectIDFromHex(reserveID)
	err := reserveCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&reserve)
	return reserve, err
}

// GetReservationsId : Obtiene una reserva con el ID
func GetReservationsId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		reserveId := c.Param("reserveId")
		var reserve models.Reservas
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(reserveId)
		err := reserveCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&reserve)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.ErrorGenericoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": reserve}})
	}
}

// CancelReservation : Cancela una reserva
func CancelReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		reserveId := c.Param("reserveId")
		objId, _ := primitive.ObjectIDFromHex(reserveId)
		update := bson.M{"estado": "Cancelado"}

		result, err := reserveCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.MatchedCount != 1 {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "No se puede cancelar la reservacion!"}})
			return
		}

		// Obtiene la reserva
		reserva, err := GetReserveId(ctx, reserveId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Actualiza el numero de asientos disponibles de la ruta
		_, _ = UpdateAvailableSeat(ctx, reserva.RutaId, reserva.AsientosReserva, false)

		c.JSON(http.StatusOK, responses.ErrorGenericoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": reserva}})

	}
}

// UpdateReservationStatus : Actualiza el estado de la reserva
func UpdateReservationStatus(ctx context.Context, reservaID string, estado string) (bool, error) {
	objId, _ := primitive.ObjectIDFromHex(reservaID)
	filter := bson.M{"_id": objId}
	update := bson.M{"estado": estado}
	result, err := reserveCollection.UpdateOne(ctx, filter, bson.M{"$set": update})
	return result.ModifiedCount != 0, err
}
