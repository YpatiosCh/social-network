"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function deleteComment(commentId) {
    try {
        const url = `/comments/${commentId}`;
        const apiResp = await serverApiRequest(url, {
            method: "DELETE"
        });

        return { success: true, ...apiResp };

    } catch (error) {
        console.error("Delete Comment Action Error:", error);
        return { success: false, error: error.message };
    }
}
