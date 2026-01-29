"use client";

import { motion } from "motion/react";
import { X, Bell } from "lucide-react";
import { constructLiveNotif } from "@/lib/notifications";
import Link from "next/link";

export default function Toast({ notification, onDismiss, onMouseEnter, onMouseLeave, onClick }) {

    const notif = constructLiveNotif(notification);

    return (
        <motion.div
            layout
            initial={{ opacity: 0, x: 100, scale: 0.9 }}
            animate={{ opacity: 1, x: 0, scale: 1 }}
            exit={{ opacity: 0, x: 100, scale: 0.9 }}
            transition={{ type: "spring", stiffness: 400, damping: 30 }}
            onMouseEnter={onMouseEnter}
            onMouseLeave={onMouseLeave}
            onClick={onClick}
            className="pointer-events-auto w-80 bg-background border border-(--border) border-l-4 border-l-(--accent) rounded-xl shadow-lg backdrop-blur-md overflow-hidden"
        >
            <div className="flex items-center gap-3 p-4">
                <div className="shrink-0 w-8 h-8 bg-(--accent)/10 rounded-full flex items-center justify-center">
                    <Bell className="w-4 h-4 text-(--accent)" />
                </div>
                <div className="flex-1 text-sm text-foreground leading-snug">
                    {notif?.who && (
                        <Link
                            href={`/profile/${notif.whoID}`}
                            prefetch={false}
                            onClick={(e) => e.stopPropagation()}
                            className="text-sm text-(--accent) hover:text-(--accent-hover) hover:underline cursor-pointer"
                        >
                            {notif.who}
                        </Link>
                    )}
                    {notif?.whoEvent && (
                        <Link
                            href={`/groups/${notif.whoID}?t=events`}
                            prefetch={false}
                            onClick={(e) => e.stopPropagation()}
                            className="text-sm text-(--accent) hover:text-(--accent-hover) hover:underline truncate cursor-pointer"
                        >
                            {notif.whoEvent}
                        </Link>
                    )}
                    <span className="text-sm text-foreground mt-0.5">
                        {notif.message}
                    </span>
                    {notif?.wherePost && (
                        <Link
                            href={`/posts/${notif.whereID}`}
                            prefetch={false}
                            onClick={(e) => e.stopPropagation()}
                            className="text-sm text-(--accent) hover:text-(--accent-hover) hover:underline cursor-pointer"
                        >
                            {notif.wherePost}
                        </Link>
                    )}
                    {notif?.whereGroup && (
                        <Link
                            href={`/groups/${notif.whereID}`}
                            prefetch={false}
                            onClick={(e) => e.stopPropagation()}
                            className="text-sm text-(--accent) hover:text-(--accent-hover) hover:underline truncate cursor-pointer"
                        >
                            {notif.whereGroup}
                        </Link>
                    )}
                    {notif?.whereEvent && (
                        <Link
                            href={`/groups/${notif.whereID}?t=events`}
                            prefetch={false}
                            onClick={(e) => e.stopPropagation()}
                            className="text-sm text-(--accent) hover:text-(--accent-hover) hover:underline truncate cursor-pointer"
                        >
                            {notif.whereEvent}
                        </Link>
                    )}
                    {notif?.extra && (
                        <p className="text-sm text-foreground mt-0.5">
                            {notif.extra}
                        </p>
                    )}
                    {/* <br></br>
                    {notif?.needs_action && (
                        <span>act</span>
                    )} */}
                </div>
                <button
                    onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        onDismiss();
                    }}
                    className="shrink-0 p-1 text-(--muted) hover:text-foreground hover:bg-(--muted)/10 rounded-full transition-colors cursor-pointer"
                >
                    <X className="w-4 h-4" />
                </button>
            </div>
        </motion.div>
    );
}
