import { serverApiRequest } from "@/lib/server-api";


/**
 * Creates a new group.
 * @param {string} group_title - The title of the group.
 * @param {string} group_description - The description of the group.
 * @param {string} group_image - The image of the group.
 * @returns {Promise<Object>} The response from the server.
 */
export async function createGroup({ group_title, group_description, group_image }) {
    try {
        const group = await serverApiRequest(`/groups/create`, {
            method: "POST",
            body: JSON.stringify({
                group_title,
                group_description,
                ...(group_image && { group_image }),
            }),
        });

        // return group_id
        return group;

    } catch (error) {
        console.error("Error creating group:", error);
        return { success: false, error: error };
    }
}