"use client";

import { safeApiCall } from "@/lib/api-wrapper";

/**
 * Client-side login function that calls the API route directly.
 * Must be called from client components so browser cookies are included.
 */
export async function loginClient(formData) {
    // Extract fields for backend request
    const identifier = formData.get("identifier")?.trim();
    const password = formData.get("password");

    // Prepare JSON payload
    const payload = {
        identifier,
        password,
    };

    const result = await safeApiCall("/api/auth/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
    });

    if (result.success) {
        return { success: true, user: result.data };
    }

    return result;
}
