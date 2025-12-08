"use client";

import { createContext, useContext, useState } from "react";
import { safeApiCall } from "@/lib/api-wrapper";


const AuthContext = createContext(null);

export function AuthProvider({ children }) {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(false);

    // Function to fetch user profile via proxy API
    const fetchUserProfile = async (userId) => {
        try {

            const url = `/api/auth/profile/${userId}`
            // Call Next.js proxy API route instead of directly calling backend
            const response = await safeApiCall(url, {
                method: "GET",
            });

            if (response.success) {
                const profileData = await response.data;
                setUser(profileData);
                return profileData;
            } else {
                console.error("Failed to fetch user profile");
                return null;
            }
        } catch (error) {
            console.error("Error fetching user profile:", error);
            return null;
        }
    };

    // Function to update user profile in context
    const updateUser = (userData) => {
        setUser(userData);
    };

    // Function to clear user on logout
    const clearUser = () => {
        setUser(null);
    };

    return (
        <AuthContext.Provider
            value={{
                user,
                loading,
                fetchUserProfile,
                updateUser,
                clearUser,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

// Custom hook to use auth context
export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}
