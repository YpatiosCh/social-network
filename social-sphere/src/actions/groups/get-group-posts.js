"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getGroupPosts({ groupId, limit, offset }) {
    try {
        const url = `/groups/${groupId}/posts?limit=${limit}&offset=${off}`;
        const response = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });

        // Return success wrapper (use 'data' for consistency)
        return { success: true, data: response };

    } catch (error) {
        console.error("Error fetching groups:", error);
        return { success: false, error: error.message };
    }
}
