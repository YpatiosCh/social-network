import { apiRequest } from "@/lib/api";

// accept or decline member as OWNER of the group
export async function handleJoinGroupRequest({ ownerId, groupId, requesterId, accepted }) {
    try {
        const response = await apiRequest("/group/handle-request", {
            method: "POST",
            body: JSON.stringify({
                owner_id: ownerId,
                group_id: groupId,
                requester_id: requesterId,
                accepted: accepted,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error handling join group request:", error);
        throw error;
    }
}

// request to join group as MEMBER of the group
export async function requestJoinGroup({ userId, groupId }) {
    try {
        const response = await apiRequest("/group/join", {
            method: "POST",
            body: JSON.stringify({
                user_id: userId,
                group_id: groupId,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error requesting to join group:", error);
        throw error;
    }
}

// respond to invitation to join group
export async function respondToGroupInvite({ groupId, accepted }) {
    try {
        const response = await apiRequest("/group/invite/response", {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId,
                accepted: accepted,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error responding to group invite:", error);
        throw error;
    }
}

/**
 * Invites a user to a group.
 * @param {string} inviterId - The ID of the user inviting the other user.
 * @param {string} groupId - The ID of the group to invite the user to.
 * @param {string} invitedId - The ID of the user being invited.
 * @returns {Promise<Object>} The response from the server.
 */
export async function inviteToGroup({ inviterId, groupId, invitedId }) {
    try {
        const response = await apiRequest("/group/invite/user", {
            method: "POST",
            body: JSON.stringify({
                inviter_id: inviterId,
                group_id: groupId,
                invited_id: invitedId,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error inviting to group:", error);
        throw error;
    }
}