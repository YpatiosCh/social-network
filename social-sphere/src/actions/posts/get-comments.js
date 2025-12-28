"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getComments({ postId, limit = 10, offset = 0 }) {
    try {
        const comments = await serverApiRequest("/comments/", {
            method: "POST",
            body: JSON.stringify({
                entity_id: postId,
                limit,
                offset
            }),
            forwardCookies: true
        });

        return { success: true, comments };
    } catch (error) {
        console.error("Error fetching comments:", error);
        return { success: false, error: error.message, comments: [] };
    }
}
