const express = require("express");
const socketio = require("socket.io");
const http = require("http");
const axios = require("axios");

const PORT = process.env.PORT || 4004;

const app = express();
const server = http.createServer(app);
const io = socketio(server, {
    cors: {
        origin: "*", // Permite que todos los orígenes accedan al WebSocket
    },
});

app.use(express.json());

// API para recibir notificaciones desde la otra API (http://localhost:8080/sensors)
app.post("/sensors-notification", (req, res) => {
    const { message } = req.body;

    if (message) {
        // Emitir la notificación recibida a todos los clientes conectados a través de WebSocket
        io.emit("nuevaNotificacion", `Notificación: ${message}`);
        console.log("Notificación recibida y emitida:", message);
        res.status(200).json({ success: true, message: "Notificación enviada a WebSocket." });
    } else {
        res.status(400).json({ success: false, message: "El mensaje es requerido." });
    }
});

// Función para simular que otro servidor (en 8080) envía un POST
async function notifyFromOtherApi() {
    try {
        // Aquí haces un POST a la API en puerto 8080 para simular la recepción de un evento
        const response = await axios.post("http://localhost:8080/sensors", {
            message: "pago registrado", // El mensaje que será notificado
        });
        console.log("Notificación enviada a la API 8080:", response.data);
    } catch (error) {
        console.error("Error al notificar a la API 8080:", error);
    }
}

io.on("connection", (socket) => {
    console.log("Cliente conectado");

    socket.emit("nuevaNotificacion", "Conexión establecida correctamente");

    socket.on("disconnect", () => {
        console.log("Cliente desconectado");
    });
});

// Iniciar servidor
server.listen(PORT, () => {
    console.log("Servidor WebSocket corriendo en el puerto: " + PORT);
    notifyFromOtherApi(); // Simula que algo sucedió en la API 8080
});
