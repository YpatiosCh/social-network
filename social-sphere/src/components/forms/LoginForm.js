"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Eye, EyeOff } from "lucide-react";
import { loginClient } from "@/services/auth/login-client";
import { useAuth } from "@/contexts/AuthContext";
import { useFormValidation } from "@/hooks/useFormValidation";

export default function LoginForm() {
    const router = useRouter();
    const { fetchUserProfile } = useAuth();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const [showPassword, setShowPassword] = useState(false);

    // Real-time validation hook
    const { errors: fieldErrors, validateField } = useFormValidation();

    async function handleSubmit(event) {
        event.preventDefault();
        setIsLoading(true)
        setError("");

        const formData = new FormData(event.currentTarget);

        try {
            const result = await loginClient(formData);

            if (result.success) {
                console.log("Login successful", result);

                // Fetch user profile using the UserId from response
                if (result.user && result.user.UserId) {
                    await fetchUserProfile(result.user.UserId);
                }

                // Navigate to feed
                router.push("/feed/public");
            } else {
                setError(result.error || "Invalid credentials");
                setIsLoading(false);
            }
        } catch (err) {
            setError("An unexpected error occurred");
            setIsLoading(false);
        }
    }

    // Real-time validation handlers
    function handleFieldValidation(name, value) {
        switch (name) {
            case "identifier":
                validateField("identifier", value, (val) => {
                    if (!val.trim()) return "Email or Username is required.";
                    return null;
                });
                break;

            case "password":
                validateField("password", value, (val) => {
                    if (!val) return "Password is required.";
                    return null;
                });
                break;
        }
    }

    return (
        <form onSubmit={handleSubmit} className="w-full space-y-6">
            {/* Email/Username Field */}
            <div>
                <label htmlFor="identifier" className="form-label pl-4">
                    Email or Username
                </label>
                <input
                    id="identifier"
                    name="identifier"
                    type="text"
                    required
                    className="form-input"
                    placeholder="Enter your email or username"
                    onChange={(e) => handleFieldValidation("identifier", e.target.value)}
                />
                {fieldErrors.identifier && (
                    <div className="form-error">{fieldErrors.identifier}</div>
                )}
            </div>

            {/* Password Field */}
            <div>
                <label htmlFor="password" className="form-label pl-4">
                    Password
                </label>
                <div className="relative">
                    <input
                        id="password"
                        name="password"
                        type={showPassword ? "text" : "password"}
                        required
                        className="form-input pr-12"
                        placeholder="Enter your password"
                        onChange={(e) => handleFieldValidation("password", e.target.value)}
                    />
                    <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                        className="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-500 hover:text-neutral-700"
                    >
                        {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                    </button>
                </div>
                {fieldErrors.password && (
                    <div className="form-error">{fieldErrors.password}</div>
                )}
            </div>

            {/* Error Message */}
            {error && (
                <div className="form-error-box animate-fade-in">
                    {error}
                </div>
            )}

            {/* Submit Button */}
            <button
                type="submit"
                disabled={isLoading}
                className="w-1/2 mx-auto flex self-center justify-center btn btn-primary"
            >
                {isLoading ? "Signing in..." : "Sign In"}
            </button>
        </form>
    );
}
