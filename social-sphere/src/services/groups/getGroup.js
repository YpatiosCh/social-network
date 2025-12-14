import { serverApiRequest } from "@/lib/api";

export async function getGroup({ groupId }) {
    try {
        const group = await serverApiRequest(`/groups/get`, {
            method: "POST",
            body: JSON.stringify({
                group_id: groupId,
            }),
        });

        return group;

    } catch (error) {
        console.error("Error getting group:", error);
        return { success: false, error: error };
    }
}