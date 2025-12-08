import { NextResponse } from 'next/server';
import { validateLoginForm } from '@/lib/validation';

export async function POST(request) {
    try {
        const payload = await request.json();

        // Validate payload
        const formData = new FormData();
        if (payload.identifier) formData.append("identifier", payload.identifier);
        if (payload.password) formData.append("password", payload.password);

        const validation = validateLoginForm(formData);
        if (!validation.valid) {
            return NextResponse.json(
                { error: validation.error },
                { status: 400 }
            );
        }

        const apiBase = process.env.API_BASE || "http://localhost:8081";
        const loginEndpoint = process.env.LOGIN || "/login";

        const backendResponse = await fetch(`${apiBase}${loginEndpoint}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });

        const responseData = await backendResponse.json().catch(() => null);
        const setCookieHeader = backendResponse.headers.get('set-cookie');

        console.log(responseData);

        const response = NextResponse.json(
            responseData || { error: "Login failed" },
            { status: backendResponse.status }
        );

        if (setCookieHeader) {
            const modifiedCookie = setCookieHeader.includes('Domain=')
                ? setCookieHeader
                : setCookieHeader + '; Domain=localhost';
            response.headers.set('Set-Cookie', modifiedCookie);
        }

        return response;
    } catch (error) {
        console.error("Login API route error:", error);
        return NextResponse.json(
            { error: "Network error. Please try again later." },
            { status: 500 }
        );
    }
}
