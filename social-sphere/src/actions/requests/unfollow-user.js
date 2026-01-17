"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function unfollowUser(userId) {
    try {
        const url = `/users/${userId}/unfollow`;
        const response = await serverApiRequest(url, {
            method: "POST",
            forwardCookies: true
        });
        return { success: true, data: response };
    } catch (error) {
        console.error("Error unfollowing user:", error);
        return { success: false, error: error.message };
    }
}
