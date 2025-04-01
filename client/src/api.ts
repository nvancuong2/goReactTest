import axios from "axios";

const API_BASE_URL = "http://127.0.0.1:5000/api/"; // Remplace par l'URL de ton API

export const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        "Content-Type": "application/json",
    },
});

// Exemple de fonction pour récupérer des données
export const fetchData = async () => {
    try {
        const response = await api.get("/todos"); // Adapter selon ton API
        return response.data;
    } catch (error) {
        console.error("Erreur lors de la récupération des données:", error);
        throw error;
    }
};
