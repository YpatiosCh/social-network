"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getFollowers({ userId, limit = 100, offset = 0 } = {}) {
    try {
        if (!userId) {
            console.error("User ID is required to fetch followers");
            return [];
        }
        const url = `/users/${userId}/followers?limit=${limit}&offset=${offset}`;
        const followers = await serverApiRequest(url, {
            method: "GET"
        });

        return followers;

    } catch (error) {
        console.error("Error fetching followers:", error);
        return [];
    }
}
