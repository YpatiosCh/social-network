"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getPendingRequests({ groupId, limit = 20, offset = 0 }) {
    try {
        const url = `/groups/${groupId}/pending-requests?limit=${limit}&offset=${offset}`;
        const response = await serverApiRequest(url, {
            method: "GET"
        });
        return response;
    } catch (error) {
        console.error("Error fetching pending requests:", error);
        return { users: [] };
    }
}
