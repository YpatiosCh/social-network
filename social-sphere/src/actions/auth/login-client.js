"use client";

import { validateLoginForm } from "@/utils/validation";

/**
 * Client-side login function that calls the API route directly.
 * Must be called from client components so browser cookies are included.
 */
export async function loginClient(formData) {
    // Client-side validation
    const validation = validateLoginForm(formData);
    if (!validation.valid) {
        return { success: false, error: validation.error };
    }

    // Extract fields for backend request
    const identifier = formData.get("identifier")?.trim();
    const password = formData.get("password");

    // Prepare JSON payload
    const payload = {
        identifier,
        password,
    };

    try {
        // Call API route directly from client
        const response = await fetch("/api/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
            credentials: "include", // Important: include cookies
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            console.error("Login failed:", errorData);
            return { success: false, error: errorData.error || "Login failed. Please try again." };
        }

        const data = await response.json();

        // Return success with user data
        return { success: true, user: data };
    } catch (error) {
        console.error("Login error:", error);
        return { success: false, error: "Network error. Please try again later." };
    }
}
