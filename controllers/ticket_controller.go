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

var ticketCollection *mongo.Collection = configs.GetCollection(configs.DB, "ticket")
var ticketValidate = validator.New()

// CreateTicket : Crea un Ticket
func CreateTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// valida el cuerpo del request
		var ticket models.Ticket
		if err := c.BindJSON(&ticket); err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Valida Campos Obligatorios
		if validationErr := ticketValidate.Struct(&ticket); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		// Crea el Pago
		ticket.Id = primitive.NewObjectID()
		if result, err := ticketCollection.InsertOne(ctx, ticket); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			c.JSON(http.StatusCreated, responses.ErrorGenericoResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
		}
	}
}

// GetTicketId : Obtiene el ticket con el ID
func GetTicketId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		ticketId := c.Param("ticketId")
		var ticket models.Ticket
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(ticketId)
		err := ticketCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&ticket)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.ErrorGenericoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": ticket}})
	}
}
