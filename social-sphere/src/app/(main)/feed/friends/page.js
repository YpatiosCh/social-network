import { LogoutButton } from "@/components/LogoutButton";
import { getFriendsPosts } from "@/actions/posts/get-friends-posts";
import PostCard from "@/components/ui/PostCard";
import CreatePost from "@/components/ui/CreatePost";

export const metadata = {
    title: "Friends Feed",
}


export default async function FriendsFeedPage() {
    const posts = await getFriendsPosts({ limit: 10, offset: 0 });
    console.log(posts);

    return (
        <div>
            <div className="pt-15 flex flex-col px-70">
                <CreatePost />
            </div>
            <div className="mt-8 mb-6">
                <h1 className="text-center feed-title">Friends Feed</h1>
                <p className="text-center feed-subtitle">What's happening in your sphere?</p>
            </div>
            <div className="section-divider mb-6" />
            <div className="pt-6 flex flex-col px-70">
                {posts?.length > 0 ? (
                    posts.map((post, index) => {
                        if (posts.length === index + 1) {
                            return (
                                <div key={`${post.ID}-${index}`}>
                                    <PostCard post={post} />
                                </div>
                            );
                        } else {
                            return <PostCard key={`${post.ID}-${index}`} post={post} />;
                        }
                    })
                ) : (
                    <div className="flex flex-col items-center justify-center py-20 animate-fade-in">
                        <p className="text-muted text-center max-w-md">
                            Your friends haven't shared anything yet. Why not be the first to start the conversation?
                        </p>
                    </div>
                )}
            </div>
            <LogoutButton />
        </div>
    );
}