const admin = require("firebase-admin");
// Requiere la clave del servicio de Firebase, la cual se utiliza para autenticar la aplicación.
const serviceAccount = require("./firebase-key.json");

admin.initializeApp({
    credential: admin.credential.cert(serviceAccount),
});
// Clase FCM que maneja el envío de notificaciones de Firebase Cloud Messaging (FCM).
class FCM {
    async send(message) {
        try {
            await admin.messaging().send(message);
            console.log("✅ Notificación enviada con éxito.");
        } catch (error) {
            console.error("❌ Error al enviar notificación:", error);
            throw error;
        }
    }
}
// Exporta la clase FCM para que pueda ser utilizada en otros módulos.
module.exports = FCM;



// Método para enviar una notificación a un usuario específico utilizando su token.
// Recibe el 'token' del dispositivo del usuario y el 'message' que contiene la notificación.
// async sendToUser(token, message) {
//     try {
//         // Creamos un nuevo objeto de mensaje que incluye el token del usuario.
//         const messageWithToken = {
//             ...message, // Usamos el operador spread para copiar el contenido de 'message'.
//             token: token, // Agregamos el token del dispositivo del usuario al mensaje.
//         };

//         // Enviamos la notificación al usuario con el token proporcionado.
//         await admin.messaging().send(messageWithToken); // Utilizamos la API de Firebase para enviar el mensaje.
//         console.log(` Notificación enviada al usuario con token: ${token}`); // Registramos en consola que la notificación fue enviada exitosamente.
//     } catch (error) {
//         // Si ocurre un error durante el proceso, lo registramos en la consola.
//         console.error(" Error al enviar notificación al usuario:", error);
//         throw error; // Lanzamos el error para manejarlo fuera del método.
//     }
// }
