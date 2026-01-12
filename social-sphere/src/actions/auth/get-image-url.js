//    "/get-image"    "thumb"

"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getImageUrl({fileId, variant}) {
    try {
        const res = await serverApiRequest("/get-image", {
            method: "POST",
            body: JSON.stringify({
                image_id: fileId,
                variant: variant
            }),
            forwardCookies: true,
            headers: {
                "Content-Type": "application/json"
            }
        });

        return {success: true, url: res.download_url};

    } catch (error) {
        console.error("Error fetching post:", error);
        return null;
    }
}
