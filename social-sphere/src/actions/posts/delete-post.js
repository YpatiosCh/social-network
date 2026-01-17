"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function deletePost(postId) {
    try {
        const url = `/posts/${postId}`;
        const apiResp = await serverApiRequest(url, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Delete Post Action Error:", error);
        return { success: false, error: error.message };
    }
}
