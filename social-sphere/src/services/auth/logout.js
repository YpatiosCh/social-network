"use client";

import { apiRequest } from "@/lib/api";

export async function logout() {
    try {
        // make api call 
        const apiResp = await apiRequest("/logout", {
            method: "POST",
        });
        
        // No need to manually clear cookies - backend already cleared the httpOnly cookie
        // Browser will automatically stop sending it on next request  
        return { success: true };

    } catch (error) {
        console.error("Logout error: ", error);
        return { success:false};
    }
}