"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function createPost(postData) {
    try {
        const url = `/posts`
        const apiResp = await serverApiRequest(url, {
            method: "POST",
            body: JSON.stringify(postData),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Create Post Action Error:", error);
        return { success: false, error: error.message };
    }
}
