"use client";

/**
 * Client-side logout function that calls the API route directly.
 * Must be called from client components so browser cookies are included.
 * @param {Function} clearUser - Optional function to clear user from AuthContext
 */
export async function logoutClient(clearUser) {
    try {
        // Call API route directly from client
        const response = await fetch("/api/auth/logout", {
            method: "POST",
            credentials: "include", // Important: include cookies
        });

        if (!response.ok) {
            console.error("Logout failed");
            return { success: false };
        }

        // Clear user from context if provided
        if (clearUser) {
            clearUser();
        }

        // Force a full page reload to clear all state and redirect to home
        window.location.href = "/";

        return { success: true };
    } catch (error) {
        console.error("Logout error:", error);
        return { success: false };
    }
}
