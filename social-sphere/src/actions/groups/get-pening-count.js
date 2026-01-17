"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getPendingRequestsCount({ groupId }) {
    try {
        const url = `/groups/${groupId}/pending-count`
        const response = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });

        return { success: true, data: response };

    } catch (error) {
        console.error("Error fetching user groups: ", error);
        return { success: false, error: error.message };
    }
}