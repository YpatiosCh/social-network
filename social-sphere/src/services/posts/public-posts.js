import { serverApiRequest } from "@/lib/server-api";

export async function getPublicPosts() {
    try {
        const posts = await serverApiRequest("/public-feed", {
            method: "POST",
        })

        return posts;

    } catch (error) {
        console.error("Error fetching public posts: ", error);
        return { success: false, error: error };
    }
}