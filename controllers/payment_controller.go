package controllers

import (
	"api-air-sales/configs"
	"api-air-sales/models"
	"api-air-sales/responses"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var paymentCollection *mongo.Collection = configs.GetCollection(configs.DB, "payment")
var paymentValidate = validator.New()

// CreatePayment : Crea un Pago
func CreatePayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// valida el cuerpo del request
		var pago models.Pagos
		if err := c.BindJSON(&pago); err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Valida Campos Obligatorios
		if validationErr := paymentValidate.Struct(&pago); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		// Crea el Pago
		pago.Id = primitive.NewObjectID()
		if result, err := paymentCollection.InsertOne(ctx, pago); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			// Actualiza el estado de la reserva
			_, _ = UpdateReservationStatus(ctx, pago.ReservaId, "Completado")

			c.JSON(http.StatusCreated, responses.ErrorGenericoResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
		}
	}
}

// CancelPayment : Cancela un Pago
func CancelPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		paymentId := c.Param("paymentId")
		objId, _ := primitive.ObjectIDFromHex(paymentId)
		update := bson.M{"estado": "Cancelado"}

		result, err := paymentCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.MatchedCount != 1 {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "No se puede cancelar el pago!"}})
			return
		}

		// Obtiene el pago
		payment, err := GetPaymentId(ctx, paymentId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Actualiza el estado de la reserva
		_, _ = UpdateReservationStatus(ctx, payment.ReservaId, "Pendiente")

		c.JSON(http.StatusOK, responses.ErrorGenericoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": payment}})

	}
}

// GetPaymentId : Obtiene el pago con el ID
func GetPaymentId(ctx context.Context, paymentId string) (models.Pagos, error) {
	var payment models.Pagos
	objId, _ := primitive.ObjectIDFromHex(paymentId)
	err := paymentCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&payment)
	return payment, err
}

// GetPaymentsId : Obtiene el pago con el ID
func GetPaymentsId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		paymentId := c.Param("paymentId")
		var payment models.Pagos
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(paymentId)
		err := paymentCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&payment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.ErrorGenericoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": payment}})
	}
}
