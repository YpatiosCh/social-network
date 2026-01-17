"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function validateUpload(fileId) {
    try {
        const url = `/files/${fileId}/validate`;
        const res = await serverApiRequest(url, {
            method: "POST",
            body: JSON.stringify({
                return_url: true
            }),
            headers: {
                "Content-Type": "application/json"
            }
        });

        return { success: true, download_url: res.download_url };

    } catch (error) {
        console.error("Validate Upload Action Error:", error);
        return { success: false, error: error.message };
    }
}
