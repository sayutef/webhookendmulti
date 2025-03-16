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
    // Convierte el mensaje a formato JSON utilizando json.Marshal
    jsonMessage, err := json.Marshal(message)
    if err != nil { // Si ocurre un error al convertir el mensaje a JSON, lo devuelve
        return err
    }
    // Crea una nueva solicitud HTTP POST hacia la URL de FCM
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


// Función comentada para obtener información sobre un dispositivo registrado
/*
func getDeviceInfo(deviceToken string) (map[string]interface{}, error) {
    // URL de la API de Firebase para obtener información del dispositivo
    deviceInfoURL := "https://fcm.googleapis.com/v1/projects/my-project-id/devices/" + deviceToken

    // Crear la solicitud GET para obtener la información del dispositivo
    req, err := http.NewRequest("GET", deviceInfoURL, nil)
    if err != nil {
        return nil, fmt.Errorf("error al crear la solicitud: %v", err)
=======



// Función comentada para registrar un dispositivo en Firebase
/*
func registerDevice(deviceToken string) error {
    // Definir el URL de registro del dispositivo en Firebase
    registerURL := "https://fcm.googleapis.com/v1/projects/my-project-id/registrationTokens/" + deviceToken

    // Crear el cuerpo de la solicitud de registro
    requestBody := map[string]string{
        "token": deviceToken,
    }
    jsonBody, err := json.Marshal(requestBody)
    if err != nil {
        return fmt.Errorf("error al convertir el cuerpo de la solicitud a JSON: %v", err)
    }

    // Realizar la solicitud HTTP para registrar el dispositivo
    req, err := http.NewRequest("POST", registerURL, bytes.NewBuffer(jsonBody))
    if err != nil {
        return fmt.Errorf("error al crear la solicitud: %v", err)

    }

    // Configurar los encabezados de la solicitud
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "key="+os.Getenv("SERVER_KEY"))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {

        return nil, fmt.Errorf("error al enviar la solicitud: %v", err)
    }
    defer resp.Body.Close()

    // Comprobar si la respuesta fue exitosa
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("error en la solicitud: %v", resp.Status)
    }

    // Decodificar la respuesta JSON en un mapa de datos
    var deviceInfo map[string]interface{}
    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(&deviceInfo)
    if err != nil {
        return nil, fmt.Errorf("error al decodificar la respuesta JSON: %v", err)
    }

    return deviceInfo, nil
}
*/

