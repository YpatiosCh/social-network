"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getImageUrl({fileId, variant}) {
    try {
        const url = `/files/images/${fileId}/${variant}`;
        const res = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true,
        });

        return {success: true, url: res.download_url};

    } catch (error) {
        console.error("Error fetching post:", error);
        return null;
    }
}
