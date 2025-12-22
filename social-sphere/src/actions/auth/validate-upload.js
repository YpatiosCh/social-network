"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function validateUpload(fileId) {
    try {
        const res = await serverApiRequest("/validate-file-upload", {
            method: "POST",
            body: JSON.stringify({
                file_id: fileId,
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
