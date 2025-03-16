require("dotenv").config();
const express = require("express");
const socketio = require("socket.io");
const http = require("http");
const axios = require("axios");
const FCM = require("./fcm");

const PORT = process.env.PORT || 4004;
const app = express();
const server = http.createServer(app);
const io = socketio(server, {
    cors: {
        origin: "*",
    },
});

const fcm = new FCM(); 

app.use(express.json());

// 📌 Lista de eventos válidos
const EVENTOS_VALIDOS = ["temperatura_alta", "temperatura_baja", "temperatura_normal"];

// 📌 Webhook para recibir eventos de la API secundaria
debugger;
app.post("/webhook", async (req, res) => {
    const { event, data } = req.body;

    console.log("Evento recibido:", event);
    console.log("Datos:", data);

    // Verificación de eventos válidos
    if (!event || !data || !EVENTOS_VALIDOS.includes(event)) {
        return res.status(400).json({ error: "Evento inválido o datos faltantes." });
    }

    let message = {
        notification: {
            title: `Nuevo Evento: ${event.replace("_", " ")}`,
            body: `Detalles: ${JSON.stringify(data)}`,
        },
        topic: "notificaciones", // Se enviará a todos los suscritos
    };

    try {
        // Enviar mensaje a Firebase Cloud Messaging
        await fcm.send(message);
        console.log(`✅ Notificación enviada para evento: ${event}`);

        // Emitir el mensaje por WebSocket a los clientes conectados
        io.emit("nuevaNotificacion", message.notification);
        
        res.status(200).json({ message: "Notificación enviada con éxito." });
    } catch (error) {
        console.error("❌ Error al enviar notificación:", error);
        res.status(500).json({ error: "Error enviando notificación" });
    }
});

// 📌 Simulación de eventos desde otro servidor en 8080
async function notifyFromOtherApi() {
    try {
        const response = await axios.post("http://localhost:8080/sensors", {
            event: "temperatura_alta",
            data: { ciudad: "México", temperatura: 35 },
        });
        console.log("📩 Notificación enviada a la API 8080:", response.data);
    } catch (error) {
        console.error("❌ Error al notificar a la API 8080:", error);
    }
}

// 📌 Manejo de conexiones WebSocket
io.on("connection", (socket) => {
    console.log("🟢 Cliente conectado");

    socket.emit("nuevaNotificacion", { title: "Bienvenido", body: "Conexión establecida correctamente" });

    socket.on("disconnect", () => {
        console.log("🔴 Cliente desconectado");
    });
});

// 📌 Iniciar el servidor
server.listen(PORT, () => {
    console.log(`🚀 Servidor WebSocket corriendo en el puerto: ${PORT}`);
    notifyFromOtherApi(); // Simular evento desde la API secundaria
});


