"use server";

import { serverApiRequest } from "@/lib/server-api";

export async function updateProfileInfo(data) {
    try {
        const url = `/my/profile`;
        const response = await serverApiRequest(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
            forwardCookies: true
        });

        return {
            success: true,
            UserId: response.UserId,
            FileId: response.FileId,
            UploadUrl: response.UploadUrl
        };
    } catch (error) {
        console.error("Error updating profile info:", error);
        return { success: false, error: error.message };
    }
}

export async function updateProfilePrivacy({ bool }) {
    try {
        const url = `/my/profile/privacy`;
        const response = await serverApiRequest(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                public: bool,
            }),
            forwardCookies: true
        });
        return response;
    } catch (error) {
        console.error("Error updating profile privacy:", error);
        return { success: false, error: error.message };
    }
}

export async function updateProfileEmail({ email }) {
    try {
        const url = `/my/profile/email`;
        const response = await serverApiRequest(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: email,
            }),
            forwardCookies: true
        });
        return response;
    } catch (error) {
        console.error("Error updating profile email:", error);
        return { success: false, error: error.message };
    }
}

export async function updateProfilePassword({ oldPassword, newPassword }) {
    try {
        const url = `/my/profile/password`;
        const response = await serverApiRequest(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                old_password: oldPassword,
                new_password: newPassword,
            }),
            forwardCookies: true
        });
        return response;
    } catch (error) {
        console.error("Error updating profile password:", error);
        return { success: false, error: error.message };
    }
}
