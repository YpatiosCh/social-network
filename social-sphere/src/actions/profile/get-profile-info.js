"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getProfileInfo(userId) {
    try {
        const url = `/users/${userId}/profile`;
        const user = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });

        return user;

    } catch (error) {
        console.error("Error fetching profile info:", error);
        return { success: false, error: error.message };
    }
}
