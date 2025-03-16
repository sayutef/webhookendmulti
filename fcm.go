package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

const fcmURL = "https://fcm.googleapis.com/fcm/send"

type Notification struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}

type FCMMessage struct {
    To           string       `json:"to"`
    Notification Notification `json:"notification"`
}

// Función para enviar notificaciones a Firebase
func sendNotification(token string, title string, body string) error {
    serverKey := os.Getenv("AIzaSyBMe-hgUxihJevy6xOk__ntpv5nndRaEcQ") // Clave de Firebase desde variable de entorno
    if serverKey == "" {
        return fmt.Errorf("clave del servidor no configurada")
    }

    message := FCMMessage{
        To: token,
        Notification: Notification{
            Title: title,
            Body:  body,
        },
    }

    jsonMessage, err := json.Marshal(message)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", fcmURL, bytes.NewBuffer(jsonMessage))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "key="+serverKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("error en la solicitud: %v", resp.Status)
    }

    return nil
}
