"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function deleteEvent(eventId) {
    try {
        const url = `/events/${eventId}`;
        const apiResp = await serverApiRequest(url, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Delete Event Action Error:", error);
        return { success: false, error: error.message };
    }
}
