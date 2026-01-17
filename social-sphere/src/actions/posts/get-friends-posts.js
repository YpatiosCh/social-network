"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getFriendsPosts({ limit = 10, offset = 0 } = {}) {
    try {
        const url = `/posts/friends?limit=${limit}&offset=${offset}`;
        const posts = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });

        return posts;

    } catch (error) {
        console.error("Error fetching friends posts:", error);
        return [];
    }
}
