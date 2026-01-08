"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getMostPopular(groupId) {
    try {
        const response = await serverApiRequest(`/groups/popular`, {
            method: "POST",
            body: JSON.stringify({
                entity_id: groupId,
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, data: response };

    } catch (error) {
        console.error("Error getting group:", error);
        return { success: false, error: error.message };
    }
}
