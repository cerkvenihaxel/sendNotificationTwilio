package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"log"
	"net/http"
	"os"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}
}

func sendWhatsAppMessage(c *gin.Context) {
	var request struct {
		PhoneNumber      string `json:"phone"`        // Número de teléfono
		MessageVariables string `json:"message_vars"` // Variables del mensaje en formato JSON
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	contentSid := os.Getenv("TWILIO_CONTENT_SID_MESSAGE")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &api.CreateMessageParams{}
	params.SetTo("whatsapp:" + request.PhoneNumber)
	params.SetFrom(fromNumber)
	params.SetContentSid(contentSid)
	params.SetContentVariables(request.MessageVariables)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Error al enviar el mensaje: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo enviar el mensaje"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mensaje enviado correctamente"})
}

func sendNotificationApproved(c *gin.Context) {
	var request struct {
		PhoneNumber      string `json:"phone"`        // Número de teléfono
		MessageVariables string `json:"message_vars"` // Variables del mensaje en formato JSON
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	contentSid := os.Getenv("TWILIO_CONTENT_SID_APPROVED")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &api.CreateMessageParams{}
	params.SetTo("whatsapp:" + request.PhoneNumber)
	params.SetFrom(fromNumber)
	params.SetContentSid(contentSid)
	params.SetContentVariables(request.MessageVariables)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Error al enviar el mensaje: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo enviar el mensaje"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mensaje enviado correctamente"})
}

func main() {
	loadEnv()

	r := gin.Default()
	r.POST("/notification-medication", sendWhatsAppMessage)
	r.POST("/notification-approved", sendNotificationApproved)
	r.Run(":8080")
}
