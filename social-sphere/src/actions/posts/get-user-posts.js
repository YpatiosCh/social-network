"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getUserPosts({ creatorId, limit = 10, offset = 0 } = {}) {
    try {
        const posts = await serverApiRequest("/user/posts", {
            method: "POST",
            body: JSON.stringify({
                creator_id: creatorId,
                limit: limit,
                offset: offset
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });

        return posts;

    } catch (error) {
        console.error("Error fetching user posts:", error);

        return [];
    }
}
