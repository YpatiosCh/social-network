"use client";

import { useState } from "react";
import Link from "next/link";
import { Menu, Bell, User } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import BlueLinkButton from "../ui/BlueLink";
import WhiteLinkButton from "../ui/WhiteLink";
import GrayBlueButton from "../ui/GreyBlueButton";

const BaseNav = ({ user }) => {
  const [menuOpen, setMenuOpen] = useState(false);
  const [notifOpen, setNotifOpen] = useState(false);

  const isAuthenticated = !!user; // you can later hook this to auth context

  return (
    <header className="sticky top-0 z-40 backdrop-blur-md bg-(--color-bg)/80 border-b border-(--color-border)">
      <nav className="flex items-center justify-between px-6 md:px-12 py-4">
        {/* Left: Logo */}
        <Link href="/" className="flex items-center gap-2 hover:opacity-80 transition-opacity">
          <img
            src="/logo.png"
            alt="SocialSphere"
            className="h-8 w-auto object-contain select-none pointer-events-none"
          />
          <span className="hidden md:inline text-lg font-semibold text-(--color-text)">
            SocialSphere
          </span>
        </Link>

        {/* Conditional render */}
        {!isAuthenticated ? (
          /* Public (hero) navigation */
          <div className="flex gap-4 justify-end">
            <WhiteLinkButton where="/login" what="Log In" />
            <BlueLinkButton where="/register" what="Join The Community" />
          </div>
        ) : (
          /* Private (user) navigation */
          <div className="flex items-center gap-4">
            {/* Center nav links */}
            <ul className="hidden md:flex items-center gap-8 text-(--color-text-muted) text-sm font-medium">
              <li><Link href="/feed" className="hover:text-=(--color-accent-hover) transition-colors">Feed</Link></li>
              <li><Link href="/communities" className="hover:text-=(--color-accent-hover) transition-colors">Communities</Link></li>
              <li><Link href="/events" className="hover:text-=(--color-accent-hover) transition-colors">Events</Link></li>
              <li><Link href="/messages" className="hover:text-=(--color-accent-hover) transition-colors">Messages</Link></li>
            </ul>

            {/* Notifications */}
            <button
              onClick={() => setNotifOpen(!notifOpen)}
              className="relative hover-glow p-2 rounded-full text-(--color-text-muted) hover:text-(--color-accent-hover) transition-colors"
            >
              <Bell size={20} />
              <span className="absolute top-1 right-1 h-2 w-2 bg-(--color-accent) rounded-full"></span>
            </button>

            {/* Profile menu */}
            <div className="relative">
              <GrayBlueButton onClick={() => setMenuOpen(!menuOpen)} something={<User size={18}/>} text={user}/>
              <AnimatePresence>
                {menuOpen && (
                  <motion.div
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -10 }}
                    className="absolute right-0 mt-3 w-44 rounded-xl bg-(--color-bg) border border-(--color-border) shadow-md overflow-hidden"
                  >
                    <Link href="/profile" className="block px-4 py-3 text-sm hover:bg-(--color-surface)">Profile</Link>
                    <Link href="/settings" className="block px-4 py-3 text-sm hover:bg-(--color-surface)">Settings</Link>
                    <hr className="border-(--color-border)" />
                    <Link href="/logout" className="block px-4 py-3 text-sm text-[#B23B3B] hover:bg-[#F7E4E4]">Log out</Link>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          </div>
        )}
      </nav>
    </header>
  );
};

export default BaseNav;
