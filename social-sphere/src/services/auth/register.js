"use client";

import { apiRequest } from "@/lib/api";

export async function register(formData) {
    try {
        // register with a public profile
        formData.append('public', 'true');

        // make API call 
        const apiResp = await apiRequest("/register", {
            method: "POST",
            body: formData,
        });

        console.log(apiResp)

        // check if user id is provided
        if (!apiResp.UserId) {
            return { 
                success: false, 
                error: "Registration failed - no user ID returned" 
            };
        }

        // all good
        return { success: true, user_id: apiResp.UserId};

    } catch (error) {
        console.error("Error: ", error);
        return { success: false, error: error};
    }
}