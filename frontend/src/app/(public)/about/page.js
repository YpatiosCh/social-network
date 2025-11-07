"use client";

import { motion } from "framer-motion";
import { Footer } from "../../../components/layout/Footer";
import Link from "next/link";

export default function About() {
    return (
        <main className="min-h-screen flex flex-col bg-[#F5F3EE] text-[#1F1B16] overflow-hidden">
            <section className="relative flex flex-col items-center text-center px-8 md:px-20 pt-20">
                <motion.h1
                    initial={{ opacity: 0, y: 25 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6 }}
                    className="text-5xl md:text-7xl font-extrabold leading-tight mb-6"
                >
                    We’re redefining <br />
                    <span className="text-[#67b0c4]">connection.</span>
                </motion.h1>
                <motion.p
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.2 }}
                    className="text-lg max-w-2xl text-[#4A453F] leading-relaxed"
                >
                    Our mission is to build a platform that amplifies genuine interaction —
                    a place where communities grow, conversations matter, and users own their digital presence.
                </motion.p>
            </section>

            {/* Divider */}
            <div className="w-full h-px bg-linear-to-r from-transparent via-[#DAD5C9] to-transparent my-12"></div>

            {/* Split: Communities */}
            <section className="flex flex-col md:flex-row items-center gap-16 px-8 md:px-20 py-16">
                <motion.div
                    initial={{ opacity: 0, x: -30 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.6 }}
                    className="md:w-1/2"
                >
                    <h2 className="text-4xl font-bold mb-4 text-[#1F1B16]">
                        Communities that evolve with you
                    </h2>
                    <p className="text-lg text-[#4A453F] leading-relaxed">
                        SocialSphere empowers users to create meaningful communities.
                        Groups form around shared passions, collaboration, and discovery — not follower counts or algorithmic trends.
                    </p>
                </motion.div>
            </section>

            {/* Divider */}
            <div className="w-full h-px bg-linear-to-r from-transparent via-[#DAD5C9] to-transparent my-12"></div>

            {/* Split: Profiles */}
            <section className="flex flex-col md:flex-row-reverse items-center gap-16 px-8 md:px-20 py-16">
                <motion.div
                    initial={{ opacity: 0, x: 30 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.6 }}
                    className="md:w-1/2"
                >
                    <h2 className="text-4xl font-bold mb-4 text-[#1F1B16]">
                        Profiles that tell your story
                    </h2>
                    <p className="text-lg text-[#4A453F] leading-relaxed">
                        Every member of SocialSphere has a space that reflects who they are — not just what they post.
                        Customize your profile, showcase your work, share your goals, and connect with others who value the same things.
                        You control the narrative — we just help you tell it beautifully.
                    </p>
                </motion.div>
            </section>

            {/* Divider */}
            <div className="w-full h-px bg-linear-to-r from-transparent via-[#DAD5C9] to-transparent my-12"></div>

            {/* Split: Privacy */}
            <section className="flex flex-col md:flex-row items-center gap-16 px-8 md:px-20 py-16">
                <motion.div
                    initial={{ opacity: 0, x: -30 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.6 }}
                    className="md:w-1/2"
                >
                    <h2 className="text-4xl font-bold mb-4 text-[#1F1B16]">
                        Privacy isn’t optional — it’s built-in
                    </h2>
                    <p className="text-lg text-[#4A453F] leading-relaxed">
                        We never trade engagement for ethics.
                        Users own their content and decide what to share.
                        Transparency, control, and respect for your data come first — always.
                    </p>
                </motion.div>
            </section>

            <div className="w-full h-px bg-linear-to-r from-transparent via-[#DAD5C9] to-transparent my-12"></div>

            {/* Closing CTA */}
            <section className="text-center pb-24 bg-linear-to-t from-[#EDE8E0] to-[#F5F3EE]">
                <motion.h2
                    initial={{ opacity: 0, y: 15 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6 }}
                    className="text-4xl md:text-5xl font-extrabold mb-8"
                >
                    Connection deserves better.
                </motion.h2>
                <p className="max-w-xl mx-auto text-lg text-[#4A453F] mb-10 leading-relaxed">
                    SocialSphere brings people together around purpose — not noise.
                    Build your profile, join a group, or start your own space.
                    Let’s build a better internet, one authentic interaction at a time.
                </p>
                <Link
                    href="/auth/register"
                    className="btn-primary"
                >
                    Join The Community
                </Link>
            </section>

            <Footer/>
        </main>
    );
}
