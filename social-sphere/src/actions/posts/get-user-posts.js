"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getUserPosts({ creatorId, limit = 10, offset = 0 } = {}) {
    try {
        const url = `/users/${creatorId}/posts?limit=${limit}&offset=${offset}`;
        const posts = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });

        return posts;

    } catch (error) {
        console.error("Error fetching user posts:", error);

        return [];
    }
}
