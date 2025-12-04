import { NextResponse } from 'next/server';

export async function GET(request, { params }) {
    try {
        const { userId } = await params;
        const cookieHeader = request.headers.get('cookie');
        console.log(cookieHeader);

        const apiBase = process.env.API_BASE || "http://localhost:8081";

        const headers = {};
        if (cookieHeader) {
            headers['Cookie'] = cookieHeader;
        }

        const backendResponse = await fetch(`${apiBase}/profile/${userId}`, {
            method: "GET",
            headers: headers,
        });

        if (!backendResponse.ok) {
            const errorData = await backendResponse.json().catch(() => null);
            return NextResponse.json(
                errorData || { error: "Failed to fetch profile" },
                { status: backendResponse.status }
            );
        }

        const profileData = await backendResponse.json();

        return NextResponse.json(profileData, { status: 200 });
    } catch (error) {
        console.error("Profile API route error:", error);
        return NextResponse.json(
            { error: "Network error. Please try again later." },
            { status: 500 }
        );
    }
}
