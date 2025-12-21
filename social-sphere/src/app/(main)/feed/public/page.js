import { LogoutButton } from "@/components/LogoutButton";
import { getPublicPosts } from "@/actions/posts/get-public-posts";
import PostCard from "@/components/ui/PostCard";
import CreatePost from "@/components/ui/CreatePost";

export const metadata = {
    title: "Public Feed",
}

export default async function PublicFeedPage() {
    // call backend for public posts
    const limit = 10;
    const offset = 0;
    const posts = await getPublicPosts({ limit, offset });

    return (
        <div>
            <div className="pt-15 flex flex-col px-70">
                <CreatePost />
            </div>
            <div className="mt-8 mb-6">
                <h1 className="text-center feed-title">Public Feed</h1>
                <p className="text-center feed-subtitle">What's happening in global sphere?</p>
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
                            Be the first ever to share something on the public sphere!
                        </p>
                    </div>
                )}
            </div>
            <LogoutButton />
        </div>
    );
}