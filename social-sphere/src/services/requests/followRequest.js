import { apiRequest } from "@/lib/api";

/**
 * Handles a follow request.
 * @param {string} requesterId - The ID of the user who sent the follow request.
 * @param {boolean} accept - Whether to accept or reject the follow request.
 * @returns {Promise<Object>} The response from the server.
 */
export const handleFollowRequest = async ({ requesterId, accept }) => {
    try {
        const response = await apiRequest("/follow-request", {
            method: "POST",
            body: JSON.stringify({
                requester_id: requesterId,
                accept: accept,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error handling follow request:", error);
        throw error;
    }
};


/**
 * Unfollows a user.
 * @param {string} userId - The ID of the user to unfollow.
 * @returns {Promise<Object>} The response from the server.
 */
export const unfollowUser = async ({ userId }) => {
    try {
        const response = await apiRequest("/user/unfollow", {
            method: "POST",
            body: JSON.stringify({
                user_id: userId,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error unfollowing user:", error);
        throw error;
    }
};