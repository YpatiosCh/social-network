"use client";

import { Footer } from "../components/layout/Footer";
import "./globals.css";
import Link from "next/link";

export default function Home() {
  return (
    <main className="min-h-screen flex flex-col bg-(--color-bg) text-(--color-text)">
      {/* Hero */}
      <section className="section flex flex-col md:flex-row flex-1 items-center justify-between gap-16">

        <div className="max-w-xl fade-in">

          <h1 className="text-5xl md:text-6xl font-extrabold leading-tight tracking-tight">
            The next chapter of <br />
            <span className="text-accent">social connection.</span>
          </h1>

          <p className="mt-6 text-lg text-muted leading-relaxed">
            Connect, share, and grow in a space built for genuine people.
            SocialSphere makes staying close effortless — no noise, no clutter, just connection.
          </p>

          <div className="mt-10 flex flex-wrap gap-4">
            <Link href="/register" className="link-primary">
              Get Started
            </Link>
            <Link href="/about" className="link-secondary">
              Learn More
            </Link>
          </div>
        </div>

        <div className="relative w-full max-w-md flex">
          <img
            src="/logo.png"
            alt="App preview"
            className="w-full h-auto object-contain select-none pointer-events-none"
            style={{
              backgroundColor: "transparent",
              mixBlendMode: "multiply",
            }}
          />
        </div>
      </section>

      <Footer />

    </main>
  );
}
