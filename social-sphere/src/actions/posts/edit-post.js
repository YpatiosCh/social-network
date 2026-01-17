"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function editPost(postData) {
    try {
        const url = `/posts/${postData.post_id}`;
        const apiResp = await serverApiRequest(url, {
            method: "POST",
            body: JSON.stringify(postData),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Edit Post Action Error:", error);
        return { success: false, error: error.message };
    }
}
