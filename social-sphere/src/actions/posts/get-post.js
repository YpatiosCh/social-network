"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function getPost(postId) {
    try {
        const url = `/posts/${postId}`;
        const post = await serverApiRequest(url, {
            method: "GET",
            forwardCookies: true
        });

        return {success: true, error: null, post: post};

    } catch (error) {
        console.error("Error fetching post:", error);
        return {success:false, error: error.message, post: null};
    }
}
