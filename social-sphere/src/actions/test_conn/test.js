"use server";

export async function test(endpoint) {
    const result = await fetch(endpoint);

    return result.json();
}
