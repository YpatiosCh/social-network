"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getFriendsPosts({ limit = 10, offset = 0 } = {}) {
    try {
        const posts = await serverApiRequest("/friends-feed", {
            method: "POST",
            body: JSON.stringify({
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
        console.error("Error fetching friends posts:", error);
        return [];
    }
}
