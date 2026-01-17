"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function respondToEvent({id, going}) {
    try {
        const url = `/events/${id}/response`;
        const apiResp = await serverApiRequest(url, {
            method: "POST",
            body: JSON.stringify(going),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Respond to Event Action Error:", error);
        return { success: false, error: error.message };
    }
}
