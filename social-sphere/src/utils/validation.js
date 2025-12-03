/**
 * Shared validation utilities for client-side and server-side validation
 */

// Email validation pattern
export const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

// Password strength pattern: at least 1 lowercase, 1 uppercase, 1 number, 1 symbol
export const STRONG_PASSWORD_PATTERN = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^\w\s]).+$/;

// Username/nickname pattern: letters, numbers, dots, underscores, dashes
export const USERNAME_PATTERN = /^[A-Za-z0-9_.-]+$/;

// Base64 image pattern for avatar validation
export const BASE64_IMAGE_PATTERN = /^data:image\/(jpeg|png|gif);base64,[A-Za-z0-9+/]+=*$/i;

/**
 * Calculate age from date of birth
 * @param {string} dateOfBirth - Date string in YYYY-MM-DD format
 * @returns {number} Age in years
 */
export function calculateAge(dateOfBirth) {
    const today = new Date();
    const birthDate = new Date(dateOfBirth);
    let age = today.getFullYear() - birthDate.getFullYear();
    const monthDiff = today.getMonth() - birthDate.getMonth();

    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
        age--;
    }

    return age;
}

/**
 * Validate email format
 * @param {string} email - Email to validate
 * @returns {boolean} True if valid
 */
export function isValidEmail(email) {
    return EMAIL_PATTERN.test(email);
}

/**
 * Validate password strength
 * @param {string} password - Password to validate
 * @returns {boolean} True if meets strength requirements
 */
export function isStrongPassword(password) {
    return password.length >= 8 && STRONG_PASSWORD_PATTERN.test(password);
}

/**
 * Validate username/nickname format
 * @param {string} username - Username to validate
 * @returns {boolean} True if valid
 */
export function isValidUsername(username) {
    return username.length >= 4 && USERNAME_PATTERN.test(username);
}

/**
 * Validate base64 image data
 * @param {string} dataUrl - Base64 data URL to validate
 * @returns {boolean} True if valid
 */
export function isValidBase64Image(dataUrl) {
    return BASE64_IMAGE_PATTERN.test(dataUrl);
}

/**
 * Validate registration form data (client-side)
 * @param {FormData} formData - Form data to validate
 * @param {string|null} avatarPreview - Base64 avatar preview
 * @returns {{valid: boolean, error: string}} Validation result
 */
export function validateRegistrationForm(formData, avatarPreview = null) {
    // First name validation
    const firstName = formData.get("firstName")?.trim() || "";
    if (!firstName) {
        return { valid: false, error: "First name is required." };
    }
    if (firstName.length < 2) {
        return { valid: false, error: "First name must be at least 2 characters." };
    }

    // Last name validation
    const lastName = formData.get("lastName")?.trim() || "";
    if (!lastName) {
        return { valid: false, error: "Last name is required." };
    }
    if (lastName.length < 2) {
        return { valid: false, error: "Last name must be at least 2 characters." };
    }

    // Email validation
    const email = formData.get("email")?.trim() || "";
    if (!isValidEmail(email)) {
        return { valid: false, error: "Please enter a valid email address." };
    }

    // Password validation
    const password = formData.get("password");
    const confirmPassword = formData.get("confirmPassword");
    if (!password || !confirmPassword) {
        return { valid: false, error: "Please enter both password and confirm password." };
    }
    if (password.length < 8) {
        return { valid: false, error: "Password must be at least 8 characters." };
    }
    if (!STRONG_PASSWORD_PATTERN.test(password)) {
        return { valid: false, error: "Password needs 1 lowercase, 1 uppercase, 1 number, and 1 symbol." };
    }
    if (password !== confirmPassword) {
        return { valid: false, error: "Passwords do not match" };
    }

    // Date of birth validation
    const dateOfBirth = formData.get("dateOfBirth")?.trim() || "";
    if (!dateOfBirth) {
        return { valid: false, error: "Date of birth is required." };
    }
    const age = calculateAge(dateOfBirth);
    if (age < 13 || age > 111) {
        return { valid: false, error: "You must be between 13 and 111 years old." };
    }

    // Nickname validation (optional)
    const username = formData.get("nickname")?.trim() || "";
    if (username) {
        if (username.length < 4) {
            return { valid: false, error: "Username must be at least 4 characters." };
        }
        if (!USERNAME_PATTERN.test(username)) {
            return { valid: false, error: "Username can only use letters, numbers, dots, underscores, or dashes." };
        }
    }

    // About me validation (optional)
    const aboutMe = formData.get("aboutMe")?.trim() || "";
    if (aboutMe && aboutMe.length > 800) {
        return { valid: false, error: "About me must be at most 800 characters." };
    }

    // Avatar validation (optional)
    if (avatarPreview) {
        if (!isValidBase64Image(avatarPreview)) {
            return { valid: false, error: "Avatar must be base64 JPEG, PNG, or GIF." };
        }
    }

    return { valid: true, error: "" };
}

/**
 * Validate login form data (client-side)
 * @param {FormData} formData - Form data to validate
 * @returns {{valid: boolean, error: string}} Validation result
 */
export function validateLoginForm(formData) {
    // Identifier validation
    const identifier = formData.get("identifier")?.trim() || "";
    if (!identifier) {
        return { valid: false, error: "Email or username is required." };
    }

    // Password validation
    const password = formData.get("password");
    if (!password) {
        return { valid: false, error: "Password is required." };
    }
    if (password.length < 8) {
        return { valid: false, error: "Password must be at least 8 characters." };
    }

    return { valid: true, error: "" };
}
