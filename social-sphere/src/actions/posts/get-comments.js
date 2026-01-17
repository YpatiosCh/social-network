"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getComments({ postId, limit = 10, offset = 0 }) {
    try {
        const url = `/comments?entity_id=${postId}&limit=${limit}&offset=${offset}`;
        const comments = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });

        return { success: true, comments };
    } catch (error) {
        console.error("Error fetching comments:", error);
        return { success: false, error: error.message, comments: [] };
    }
}
