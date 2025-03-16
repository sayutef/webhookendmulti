const admin = require("firebase-admin");

const serviceAccount = require("./firebase-key.json");

admin.initializeApp({
    credential: admin.credential.cert(serviceAccount),
});

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

module.exports = FCM;






