import { apiRequest } from "@/lib/api";

export async function getGroupMembers({ group_id, limit, offset }) {
    try {
        const members = await apiRequest("/group/members", {
            method: "POST",
            body: JSON.stringify({
                group_id: group_id,
                limit: limit,
                offset: offset
            })
        })

        return members;

    } catch (error) {
        console.error("Error fetching group members: ", error);
        return { success: false, error: error };
    }
}
