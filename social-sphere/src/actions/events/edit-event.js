"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function editEvent({id , data}) {
    try {
        const url = `/events/${id}`;
        const apiResp = await serverApiRequest(url, {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Edit Event Action Error:", error);
        return { success: false, error: error.message };
    }
}
