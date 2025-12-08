"use client";

import { use, useState, useEffect, useCallback } from "react";
import ProfileHeader from "@/components/features/profile/profile-header";
import { Lock } from "lucide-react";
import { fetchUserProfile } from "@/services/profile/profile-actions";
import { fetchUserPosts } from "@/services/posts/posts";
import { getUserByID } from "@/mock-data/users";
import FeedList from "@/components/feed/feed-list";
import CreatePost from "@/components/ui/create-post";


export default function ProfilePage({ params }) {
    const { id } = use(params);
    const [loading, setLoading] = useState(true);
    const [user, setUser] = useState(null);
    const [initialPosts, setInitialPosts] = useState([]);
    const [postsLoaded, setPostsLoaded] = useState(false);

    // mock data
    const currentUser = getUserByID("1");

    // Data Fetching
    useEffect(() => {
        const loadUser = async () => {
            try {
                const data = await fetchUserProfile(id);
                setUser(data);
            } catch (error) {
                console.error("Failed to fetch user:", error);
            } finally {
                setLoading(false);
            }
        };

        loadUser();
    }, [id]);

    useEffect(() => {
        if (user) {
            const loadPosts = async () => {
                try {
                    const posts = await fetchUserPosts(user.ID, 0, 5);
                    setInitialPosts(posts);
                } catch (error) {
                    console.error("Failed to fetch posts:", error);
                } finally {
                    setPostsLoaded(true);
                }
            };
            loadPosts();
        }
    }, [user]);

    const fetchPosts = useCallback(async (offset, limit) => {
        if (!user) return [];
        return await fetchUserPosts(user.ID, offset, limit);
    }, [user]);

    console.log("user", user);
    console.log("initialPosts", initialPosts);
    console.log("postsLoaded", postsLoaded);

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[50vh]">
                <div className="w-8 h-8 border-4 border-(--foreground) border-t-transparent rounded-full animate-spin" />
            </div>
        );
    }

    if (!user) {
        return (
            <div className="flex items-center justify-center min-h-[50vh]">
                <div className="w-8 h-8 border-4 border-(--foreground) border-t-transparent rounded-full animate-spin" />
            </div>
        );
    }

    // Check if profile is private and viewer is not following (and not owner)
    const isOwnProfile = user.ID === currentUser.ID; // Mock check
    const isPrivateView = !user.publicProf && !user.isFollower && !isOwnProfile;


    return (
        <div className="w-full py-8 animate-in fade-in duration-500">
            <div className="max-w-7xl mx-auto px-6">
                <ProfileHeader user={user} isOwnProfile={isOwnProfile} />

                <div className="flex gap-6 mt-6">
                    {/* Left Sidebar - Spacer to match feed alignment */}
                    <aside className="hidden xl:block w-48 shrink-0" />

                    {/* Main Content */}
                    <main className="flex-1 max-w-2xl mx-auto min-w-0">
                        {isPrivateView ? (
                            <div className="flex flex-col items-center justify-center py-24 text-center bg-(--muted)/5 rounded-2xl border border-(--muted)/10">
                                <div className="w-16 h-16 rounded-full bg-(--muted)/10 flex items-center justify-center mb-4">
                                    <Lock className="w-8 h-8 text-(--muted)" />
                                </div>
                                <h2 className="text-xl font-bold mb-2">This profile is private</h2>
                                <p className="text-(--muted) max-w-md">
                                    Follow this account to see their photos and videos.
                                </p>
                            </div>
                        ) : (
                            <div className="mt-8">
                                {postsLoaded ? (
                                    <div>
                                        <CreatePost onPostCreated={CreatePost} />
                                        <FeedList initialPosts={initialPosts} fetchPosts={fetchPosts} />
                                    </div>
                                ) : (
                                    <div className="flex justify-center p-4">
                                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
                                    </div>
                                )}
                            </div>
                        )}
                    </main>

                    {/* Right Sidebar - Reserved for future widgets */}
                    <aside className="hidden lg:block w-80 shrink-0" />
                </div>
            </div>
        </div>
    );
}
