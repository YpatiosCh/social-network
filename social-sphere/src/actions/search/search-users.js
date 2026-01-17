"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function SearchUsers({ query, limit }) {
    try {
        const url = `/users/search?query=${query}&limit=${limit}`;
        const response = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });
        return response;
    } catch (error) {
        console.error("Error searching users:", error);
        return { success: false, error: error.message };
    }
}
