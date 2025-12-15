import { apiRequest } from "@/lib/api";
import { serverApiRequest } from "@/lib/server-api";

/**
 * Updates the user's profile privacy.
 * @param { boolean } bool - Whether the profile is public or not.
 * @returns {Promise<Object>} The response from the server.
 */
export async function updateProfilePrivacy({ bool }) {
    const isServer = typeof window === 'undefined';
    const apiFn = isServer ? serverApiRequest : apiRequest;

    try {
        const response = await apiFn("/account/update/public", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
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
    const isServer = typeof window === 'undefined';
    const apiFn = isServer ? serverApiRequest : apiRequest;

    try {
        const response = await apiFn("/account/update/email", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
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
 * @param {string} oldPassword - The old password.
 * @param {string} newPassword - The new password.
 * @returns {Promise<Object>} The response from the server.
 */
export async function updateProfilePassword({ oldPassword, newPassword }) {
    const isServer = typeof window === 'undefined';
    const apiFn = isServer ? serverApiRequest : apiRequest;

    try {
        const response = await apiFn("/account/update/password", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                old_password: oldPassword,
                new_password: newPassword,
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
 * @param {string} username - The new username.
 * @param {string} first_name - The new first name.
 * @param {string} last_name - The new last name.
 * @param {Date} date_of_birth - The new date of birth.
 * @param {string} avatar_id - The new avatar ID.
 * @param {string} about - The new about text.
 * @returns {Promise<Object>} The response from the server.
 */
export async function updateProfileInfo({ username, first_name, last_name, date_of_birth, avatar_id, about }) {
    const isServer = typeof window === 'undefined';
    const apiFn = isServer ? serverApiRequest : apiRequest;

    try {
        const response = await apiFn("/profile/update", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ username, first_name, last_name, date_of_birth, avatar_id, about }),
        });
        return response;
    } catch (error) {
        console.error("Error updating profile:", error);
        throw error;
    }
};
