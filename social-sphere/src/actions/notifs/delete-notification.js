"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function deleteNotification(notificationId) {
    try {
        const url = `/notifications/${notificationId}`;
        await serverApiRequest(url, {
            method: "DELETE",
            forwardCookies: true
        });
        return { success: true };
    } catch (error) {
        console.error("Error deleting notification:", error);
        return { success: false, error: error.message };
    }
}
