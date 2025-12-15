import { apiRequest } from "@/lib/api";

/**
 * Leaves a group.
 * @param {string} userId - The ID of the user leaving the group.
 * @param {string} groupId - The ID of the group to leave.
 * @returns {Promise<Object>} The response from the server.
 */
export async function leaveGroup({ userId, groupId }) {
    try {
        const response = await apiRequest("/group/leave", {
            method: "POST",
            body: JSON.stringify({
                user_id: userId,
                group_id: groupId,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error leaving group:", error);
        throw error;
    }
}


