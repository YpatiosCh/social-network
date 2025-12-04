"use client";

import { validateRegistrationForm } from "@/utils/validation";

/**
 * Client-side registration function that calls the API route directly.
 * Must be called from client components so browser cookies are included.
 */
export async function registerClient(formData) {
    // Extract avatar file
    const avatar = formData.get("avatar");

    // Client-side validation
    const validation = validateRegistrationForm(formData, avatar);
    if (!validation.valid) {
        return { success: false, error: validation.error };
    }

    // Extract fields for backend request
    const firstName = formData.get("firstName")?.trim();
    const lastName = formData.get("lastName")?.trim();
    const email = formData.get("email")?.trim();
    const password = formData.get("password");
    const dateOfBirth = formData.get("dateOfBirth")?.trim();
    const nickname = formData.get("nickname")?.trim();
    const aboutMe = formData.get("aboutMe")?.trim();

    // Prepare data for backend
    const backendFormData = new FormData();
    backendFormData.append("first_name", firstName);
    backendFormData.append("last_name", lastName);
    backendFormData.append("email", email);
    backendFormData.append("password", password);
    backendFormData.append("date_of_birth", dateOfBirth);
    backendFormData.append("public", "true"); // Default to public

    if (nickname) backendFormData.append("username", nickname);
    if (aboutMe) backendFormData.append("about", aboutMe);
    if (avatar && avatar.size > 0) backendFormData.append("avatar", avatar);

    try {
        // Call API route directly from client
        const response = await fetch("/api/auth/register", {
            method: "POST",
            body: backendFormData,
            credentials: "include", // Important: include cookies
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            console.error("Registration failed, Backend response:", errorData);
            return { success: false, error: errorData.error || "Registration failed. Please try again." };
        }

        return { success: true };
    } catch (error) {
        console.error("Registration error:", error);
        return { success: false, error: "Network error. Please try again later." };
    }
}
