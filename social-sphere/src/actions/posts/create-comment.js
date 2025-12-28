"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function createComment({ postId, commentBody, imageName, imageSize, imageType }) {
    try {
        const body = {
            parent_id: postId,
            comment_body: commentBody
        };

        // Add image info if provided
        if (imageName && imageSize && imageType) {
            body.image_name = imageName;
            body.image_size = imageSize;
            body.image_type = imageType;
        }

        const response = await serverApiRequest("/comments/create", {
            method: "POST",
            body: JSON.stringify(body),
            forwardCookies: true
        });

        return { success: true, ...response };
    } catch (error) {
        console.error("Error creating comment:", error);
        return { success: false, error: error.message };
    }
}
