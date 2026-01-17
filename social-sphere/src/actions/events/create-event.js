"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function createEvent({groupID, data}) {
    try {
        const url = `/groups/${groupID}/events`
        const apiResp = await serverApiRequest(url, {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Create Event Action Error:", error);
        return { success: false, error: error.message };
    }
}
