"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function removeEventResponse({id}) {
    try {
        const url = `/events/${id}/response`;
        const apiResp = await serverApiRequest(url, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Remove Event Response Action Error:", error);
        return { success: false, error: error.message };
    }
}
