import { apiRequest } from "@/lib/api";

/**
 * Searches for users.
 * @param {string} query - The search query.
 * @param {number} limit - The number of results to return.
 * @returns {Promise<Object>} The response from the server.
 */
export const SearchUsers = async ({ query, limit }) => {
    try {
        const response = await apiRequest("/users/search", {
            method: "POST",
            body: JSON.stringify({
                query: query,
                limit: limit,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error searching users:", error);
        throw error;
    }
};