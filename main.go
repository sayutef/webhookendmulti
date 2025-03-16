package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

// Eventos válidos
var eventosValidos = map[string]bool{
    "temperatura_alta":  true,
    "temperatura_baja":  true,
    "temperatura_normal": true,
}

func main() {
    r := gin.Default()

    // Ruta para recibir eventos desde sensores
    r.POST("/sensors", func(c *gin.Context) {
        var sensorData struct {
            UserID string                 `json:"user_id"`
            Event  string                 `json:"event"`
            Data   map[string]interface{} `json:"data"`
        }

        if err := c.BindJSON(&sensorData); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
            return
        }

        if _, valido := eventosValidos[sensorData.Event]; !valido {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Evento no válido"})
            return
        }

        fmt.Printf("Sensor envió evento: Usuario %s - Evento %s\n", sensorData.UserID, sensorData.Event)

        // Llamar al webhook automáticamente
        go enviarWebhook(sensorData)

        c.JSON(http.StatusOK, gin.H{"message": "Evento recibido y procesado"})
    })

    // Webhook que escucha eventos
    r.POST("/webhook", func(c *gin.Context) {
        var event struct {
            UserID string                 `json:"user_id"`
            Event  string                 `json:"event"`
            Data   map[string]interface{} `json:"data"`
        }

        if err := c.BindJSON(&event); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
            return
        }

        fmt.Printf("Webhook recibió evento: Usuario %s - Evento %s\n", event.UserID, event.Event)

        // Obtener el token del dispositivo
        token := os.Getenv("DEVICE_TOKEN")
        if token == "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Token de dispositivo no configurado"})
            return
        }

        // Enviar notificación a Firebase
        err := sendNotification(token, "Nuevo Evento: "+event.Event, fmt.Sprintf("Detalles: %v", event.Data))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo enviar la notificación"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Webhook procesado y notificación enviada"})
    })

    r.Run(":8081") // API ejecutándose en el puerto 8080
}

// Función para enviar datos al webhook
func enviarWebhook(eventData interface{}) {
    webhookURL := "http://localhost:8080/sensors"

    jsonData, err := json.Marshal(eventData)
    if err != nil {
        fmt.Println("Error al convertir evento a JSON:", err)
        return
    }

    req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error al crear la solicitud:", err)
        return
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error al enviar webhook:", err)
        return
    }
    defer resp.Body.Close()

    fmt.Println("Webhook enviado con estado:", resp.Status)
}
