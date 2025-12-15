import { apiRequest } from "@/lib/api";

/**
 * Searches for groups.
 * @param {string} query - The search query.
 * @param {number} limit - The number of results to return.
 * @param {number} offset - The offset for pagination.
 * @returns {Promise<Object>} The response from the server.
 */
export const SearchGroups = async ({ query, limit, offset }) => {
    try {
        const response = await apiRequest("/search/group", {
            method: "POST",
            body: JSON.stringify({
                query: query,
                limit: limit,
                offset: offset,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error searching groups:", error);
        throw error;
    }
};