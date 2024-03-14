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

var routeCollection *mongo.Collection = configs.GetCollection(configs.DB, "route")
var routeValidate = validator.New()

// CreateRoute : Crea una Ruta
func CreateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// valida el cuerpo del request
		var ruta models.Rutas
		if err := c.BindJSON(&ruta); err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Valida Campos Obligatorios
		if validationErr := routeValidate.Struct(&ruta); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		// Crea la Ruta
		ruta.Id = primitive.NewObjectID()
		if result, err := routeCollection.InsertOne(ctx, ruta); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			c.JSON(http.StatusCreated, responses.ErrorGenericoResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
		}
	}
}

// GetRouteFilters : Obtiene Rutas con Filtros
func GetRouteFilters() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		origen := c.Query("origen")
		destino := c.Query("destino")
		fechaIda := c.Query("fechaIda")
		fechaVuelta := c.Query("fechaVuelta")

		var rutas []models.Rutas
		var filter bson.M

		if origen == "" || destino == "" {
			c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Se requieren origen y destino"}})
			return
		}

		filter = bson.M{"origen": origen, "destino": destino}

		if fechaIda != "" && fechaVuelta != "" {
			start, err := time.Parse(time.RFC3339, fechaIda)
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Formato de fecha de ida inválido"}})
				return
			}

			end, err := time.Parse(time.RFC3339, fechaVuelta)
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Formato de fecha de vuelta inválido"}})
				return
			}

			filter["$and"] = []bson.M{
				{"fecha": bson.M{"$gte": start}},
				{"fecha": bson.M{"$lt": end}},
			}
		} else if fechaIda != "" {
			start, err := time.Parse(time.RFC3339, fechaIda)
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Formato de fecha de ida inválido"}})
				return
			}

			filter["fecha"] = start
		}

		results, err := routeCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		defer results.Close(ctx)

		for results.Next(ctx) {
			var singleRuta models.Rutas
			if err := results.Decode(&singleRuta); err != nil {
				c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			rutas = append(rutas, singleRuta)
		}
		if len(rutas) > 0 {
			c.JSON(http.StatusOK, responses.ErrorGenericoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": rutas}})
		} else {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "No se encontraron registros de Rutas para esos datos"}})
		}
	}
}

// GetRouteId : Obtiene una ruta con el ID
func GetRouteId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		routeId := c.Param("routeId")
		var ruta models.Rutas
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(routeId)
		err := routeCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&ruta)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorGenericoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.ErrorGenericoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": ruta}})
	}
}

// UpdateAvailableSeat : Actualiza el numero de asientos disponibles de la ruta
func UpdateAvailableSeat(ctx context.Context, rutaID string, asientosReserva int, estado bool) (bool, error) {
	objId, _ := primitive.ObjectIDFromHex(rutaID)
	filter := bson.M{"_id": objId}
	if estado {
		asientosReserva = -asientosReserva
	}
	update := bson.M{"$inc": bson.M{"asientos_disponibles": asientosReserva}}
	result, err := routeCollection.UpdateOne(ctx, filter, update)
	return result.ModifiedCount != 0, err
}
