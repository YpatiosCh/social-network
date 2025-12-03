"use server";

import {
    isValidEmail,
    isValidUsername,
    isValidBase64Image,
    calculateAge,
    STRONG_PASSWORD_PATTERN
} from "@/utils/validation";

export async function login(formData) {
    const identifier = formData.get("identifier")?.trim() || "";
    const password = formData.get("password");

    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 1000));

    // Server-side validation
    if (!identifier) {
        return { success: false, error: "Email or username is required." };
    }

    if (!password) {
        return { success: false, error: "Password is required." };
    }

    // Mock authentication logic
    // In a real app, this would make a server-to-server API call to the Golang backend
    if (identifier && password) {
        console.log("Server Action: Login successful", { identifier });
        return { success: true };
    } else {
        return { success: false, error: "Invalid credentials" };
    }
}

export async function register(formData) {
    // Extract and sanitize all fields
    const firstName = formData.get("firstName")?.trim() || "";
    const lastName = formData.get("lastName")?.trim() || "";
    const email = formData.get("email")?.trim() || "";
    const password = formData.get("password");
    const confirmPassword = formData.get("confirmPassword");
    const dateOfBirth = formData.get("dateOfBirth")?.trim() || "";
    const nickname = formData.get("nickname")?.trim() || "";
    const aboutMe = formData.get("aboutMe")?.trim() || "";
    const avatar = formData.get("avatar"); // Base64 string

    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 1500));

    // Server-side validation (defense in depth - never trust client)

    // First name validation
    if (!firstName) {
        return { success: false, error: "First name is required." };
    }
    if (firstName.length < 2) {
        return { success: false, error: "First name must be at least 2 characters." };
    }

    // Last name validation
    if (!lastName) {
        return { success: false, error: "Last name is required." };
    }
    if (lastName.length < 2) {
        return { success: false, error: "Last name must be at least 2 characters." };
    }

    // Email validation
    if (!isValidEmail(email)) {
        return { success: false, error: "Please enter a valid email address." };
    }

    // Password validation
    if (!password || !confirmPassword) {
        return { success: false, error: "Please enter both password and confirm password." };
    }
    if (password.length < 8) {
        return { success: false, error: "Password must be at least 8 characters." };
    }
    if (!STRONG_PASSWORD_PATTERN.test(password)) {
        return { success: false, error: "Password needs 1 lowercase, 1 uppercase, 1 number, and 1 symbol." };
    }
    if (password !== confirmPassword) {
        return { success: false, error: "Passwords do not match" };
    }

    // Date of birth validation
    if (!dateOfBirth) {
        return { success: false, error: "Date of birth is required." };
    }
    const age = calculateAge(dateOfBirth);
    if (age < 13 || age > 111) {
        return { success: false, error: "You must be between 13 and 111 years old." };
    }

    // Nickname validation (optional)
    if (nickname) {
        if (nickname.length < 4) {
            return { success: false, error: "Username must be at least 4 characters." };
        }
        if (!isValidUsername(nickname)) {
            return { success: false, error: "Username can only use letters, numbers, dots, underscores, or dashes." };
        }
    }

    // About me validation (optional)
    if (aboutMe && aboutMe.length > 800) {
        return { success: false, error: "About me must be at most 800 characters." };
    }

    // Avatar validation (optional)
    if (avatar) {
        if (!isValidBase64Image(avatar)) {
            return { success: false, error: "Avatar must be base64 JPEG, PNG, or GIF." };
        }
    }

    // Mock registration logic
    // In a real app, this would make a server-to-server API call to the Golang backend
    console.log("Server Action: Registration successful", {
        email,
        firstName,
        lastName,
        nickname,
        aboutMe,
        hasAvatar: !!avatar
    });
    return { success: true };
}

