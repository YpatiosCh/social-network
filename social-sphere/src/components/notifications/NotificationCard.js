"use client";

import { useState } from "react";
import Link from "next/link";
import { Bell, Trash2, Check, X, User, Users, Calendar, Heart, MessageCircle } from "lucide-react";
import { markNotificationAsRead } from "@/actions/notifs/mark-as-read";
import { deleteNotification } from "@/actions/notifs/delete-notification";
import { handleFollowRequest } from "@/actions/requests/handle-request";
import { respondToGroupInvite } from "@/actions/groups/respond-to-invite";
import { handleJoinRequest } from "@/actions/groups/handle-join-request";
import { getRelativeTime } from "@/lib/time";

export default function NotificationCard({ notification, onDelete, onUpdate }) {
    const [isActing, setIsActing] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);
    const [acted, setActed] = useState(notification.acted);

    const needsAction = notification.needs_action && !acted;

    const getNotificationIcon = () => {
        switch (notification.type) {
            case "follow_request":
            case "new_follower":
            case "follow_request_accepted":
                return <User className="w-4 h-4" />;
            case "group_invite":
            case "group_join_request":
            case "group_join_request_accepted":
            case "group_invite_accepted":
                return <Users className="w-4 h-4" />;
            case "new_event":
                return <Calendar className="w-4 h-4" />;
            case "post_reply":
                return <MessageCircle className="w-4 h-4" />;
            case "like":
                return <Heart className="w-4 h-4" />;
            default:
                return <Bell className="w-4 h-4" />;
        }
    };

    // Construct notification message - placeholder, will be refined later
    const constructMessage = () => {
        const { type, payload, count } = notification;

        switch (type) {
            case "new_follower":
                return {
                    who: payload?.follower_name,
                    whoId: payload?.follower_id,
                    message: " started following you"
                };
            case "follow_request":
                return {
                    who: payload?.requester_name,
                    whoId: payload?.requester_id,
                    message: " wants to follow you"
                };
            case "follow_request_accepted":
                return {
                    who: payload?.target_name,
                    whoId: payload?.target_id,
                    message: " accepted your follow request"
                };
            case "post_reply":
                return {
                    who: payload?.commenter_name,
                    whoId: payload?.commenter_id,
                    message: " commented on your post",
                    link: `/posts/${payload?.post_id}`
                };
            case "like":
                return {
                    who: payload?.liker_name,
                    whoId: payload?.liker_id,
                    message: " liked your post",
                    link: `/posts/${payload?.post_id}`
                };
            case "group_invite":
                return {
                    who: payload?.inviter_name,
                    whoId: payload?.inviter_id,
                    message: " invited you to join ",
                    groupName: payload?.group_name,
                    groupId: payload?.group_id
                };
            case "group_join_request":
                return {
                    who: payload?.requester_name,
                    whoId: payload?.requester_id,
                    message: " wants to join ",
                    groupName: payload?.group_name,
                    groupId: payload?.group_id
                };
            case "group_join_request_accepted":
                return {
                    message: "You were accepted to ",
                    groupName: payload?.group_name,
                    groupId: payload?.group_id
                };
            case "group_invite_accepted":
                return {
                    who: payload?.invited_name,
                    whoId: payload?.invited_id,
                    message: " accepted your invitation to ",
                    groupName: payload?.group_name,
                    groupId: payload?.group_id
                };
            case "new_event":
                return {
                    message: "New event ",
                    eventTitle: payload?.event_title,
                    groupName: payload?.group_name,
                    groupId: payload?.group_id
                };
            default:
                return { message: "You have a new notification" };
        }
    };

    const handleAction = async (accept) => {
        setIsActing(true);

        try {
            let result;
            const { type, payload } = notification;

            switch (type) {
                case "follow_request":
                    result = await handleFollowRequest({
                        requesterId: payload.requester_id,
                        accept
                    });
                    break;
                case "group_invite":
                    result = await respondToGroupInvite({
                        groupId: payload.group_id,
                        accept
                    });
                    break;
                case "group_join_request":
                    result = await handleJoinRequest({
                        groupId: payload.group_id,
                        requesterId: payload.requester_id,
                        accepted: accept
                    });
                    break;
                default:
                    return;
            }

            if (result?.success) {
                // Optimistically update acted state
                setActed(true);
                // Mark as read
                await markNotificationAsRead(notification.id);
                onUpdate?.(notification.id, { acted: true });
            }
        } catch (error) {
            console.error("Error handling notification action:", error);
        } finally {
            setIsActing(false);
        }
    };

    const handleDelete = async () => {
        setIsDeleting(true);
        try {
            const result = await deleteNotification(notification.id);
            if (result.success) {
                onDelete?.(notification.id);
            }
        } catch (error) {
            console.error("Error deleting notification:", error);
        } finally {
            setIsDeleting(false);
        }
    };

    const handleMarkAsRead = async () => {
        try {
            await markNotificationAsRead(notification.id);
        } catch (error) {
            console.error("Error marking as read:", error);
        }
    };

    const content = constructMessage();

    const isSeen = notification.seen;

    return (
        <div className={`group bg-background border border-(--border) rounded-xl p-4 transition-all hover:border-(--muted)/40 hover:shadow-sm ${isSeen ? "opacity-40" : ""}`}>
            <div className="flex items-start gap-3">
                {/* Icon */}
                <div className="shrink-0 w-10 h-10 bg-(--accent)/10 rounded-full flex items-center justify-center text-(--accent)">
                    {getNotificationIcon()}
                </div>

                {/* Content */}
                <div className="flex-1 min-w-0">
                    <div className="text-sm text-foreground leading-snug">
                        {content.who && (
                            <Link
                                href={`/profile/${content.whoId}`}
                                className="font-semibold text-(--accent) hover:underline"
                            >
                                {content.who}
                            </Link>
                        )}
                        <span>{content.message}</span>
                        {content.groupName && (
                            <Link
                                href={`/groups/${content.groupId}`}
                                className="font-semibold text-(--accent) hover:underline"
                            >
                                {content.groupName}
                            </Link>
                        )}
                        {content.eventTitle && (
                            <>
                                <span className="font-semibold">{content.eventTitle}</span>
                                <span> in </span>
                                <Link
                                    href={`/groups/${content.groupId}?t=events`}
                                    className="font-semibold text-(--accent) hover:underline"
                                >
                                    {content.groupName}
                                </Link>
                            </>
                        )}
                        {content.link && (
                            <Link
                                href={content.link}
                                className="text-(--accent) hover:underline ml-1"
                            >
                                View
                            </Link>
                        )}
                    </div>

                    <p className="text-xs text-(--muted) mt-1">
                        {getRelativeTime(notification.created_at)}
                    </p>

                    {/* Action Buttons for actionable notifications */}
                    {needsAction && (
                        <div className="flex items-center gap-2 mt-3">
                            <button
                                onClick={() => handleAction(true)}
                                disabled={isActing}
                                className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-(--accent) text-white hover:bg-(--accent-hover) rounded-full transition-colors disabled:opacity-50 cursor-pointer"
                            >
                                <Check className="w-3 h-3" />
                                Accept
                            </button>
                            <button
                                onClick={() => handleAction(false)}
                                disabled={isActing}
                                className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium border border-(--border) text-foreground hover:bg-(--muted)/10 rounded-full transition-colors disabled:opacity-50 cursor-pointer"
                            >
                                <X className="w-3 h-3" />
                                Decline
                            </button>
                        </div>
                    )}

                    {/* Acted indicator */}
                    {acted && notification.needs_action && (
                        <p className="text-xs text-(--muted) mt-2 italic">Responded</p>
                    )}
                </div>

                {/* Delete Button */}
                <button
                    onClick={handleDelete}
                    disabled={isDeleting}
                    className="shrink-0 p-2 text-(--muted) hover:text-red-500 hover:bg-red-500/5 rounded-full transition-colors opacity-0 group-hover:opacity-100 disabled:opacity-50 cursor-pointer"
                >
                    <Trash2 className="w-4 h-4" />
                </button>
            </div>
        </div>
    );
}
