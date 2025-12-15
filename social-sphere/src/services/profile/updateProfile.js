import { apiRequest } from "@/lib/api";

/**
 * Updates the user's profile privacy.
 * @param { boolean } bool - Whether the profile is public or not.
 * @returns {Promise<Object>} The response from the server.
 */
export async function updateProfilePrivacy({ bool }) {
    try {
        const response = await apiRequest("/account/update/public", {
            method: "POST",
            body: JSON.stringify({
                public: bool,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error updating profile:", error);
        throw error;
    }
};

/**
 * Updates the user's profile email.
 * @param { string } email - The new email address.
 * @returns {Promise<Object>} The response from the server.
 */
export async function updateProfileEmail({ email }) {
    try {
        // TODO: Add email validation
        const response = await apiRequest("/account/update/email", {
            method: "POST",
            body: JSON.stringify({
                email: email,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error updating profile:", error);
        throw error;
    }
};

/**
 * Updates the user's profile password.
 * @param { string } password - The new password.
 * @returns {Promise<Object>} The response from the server.
 */
export async function updateProfilePassword({ password }) {
    try {
        // TODO: Add password validation
        const response = await apiRequest("/account/update/password", {
            method: "POST",
            body: JSON.stringify({
                password: password,
            }),
        });
        return response;
    } catch (error) {
        console.error("Error updating profile:", error);
        throw error;
    }
};

/**
 * Updates the user's profile information.
 * @param {string} first_name - The new first name.
 * @param {string} last_name - The new last name.
 * @param {Date} date_of_birth - The new date of birth.
 * @param {string} avatar_id - The new avatar ID.
 * @param {string} about - The new about text.
 * @returns {Promise<Object>} The response from the server.
 */
export async function updateProfileInfo({ first_name, last_name, date_of_birth, avatar_id, about }) {
    try {
        const response = await apiRequest("/profile/update", {
            method: "POST",
            body: JSON.stringify({ first_name, last_name, date_of_birth, avatar_id, about }),
        });
        return response;
    } catch (error) {
        console.error("Error updating profile:", error);
        throw error;
    }
};
