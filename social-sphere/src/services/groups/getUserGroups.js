import { serverApiRequest } from "@/lib/server-api";

/**
 * Gets the user's groups.
 * @param {number} pagination.limit - The number of groups to return.
 * @param {number} pagination.offset - The offset for pagination.
 * @returns {Promise<Object>} The response from the server.
 */
export async function getUserGroups({ limit, offset }) {
    try {
        const groups = await serverApiRequest("/groups/user/", {
            method: "POST",
            body: JSON.stringify({
                limit: limit,
                offset: offset,
            })
        })

        return groups;

    } catch (error) {
        console.error("Error fetching user groups: ", error);
        return { success: false, error: error };
    }
}
