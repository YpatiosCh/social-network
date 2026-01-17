"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function inviteToGroup({ groupId, invitedIds }) {
    try {
        const url = `/groups/${groupId}/invite`;
        const response = await serverApiRequest(url, {
            method: "POST",
            body: JSON.stringify({
                invited_id: invitedIds
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error inviting to group:", error);
        return { success: false, error: error.message };
    }
}
