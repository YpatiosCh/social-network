export default function FriendsFeedPage() {
    // Mock data based on the requested struct: Username, Content, CreatedAt, NumOfComments
    const posts = [
        {
            username: "watermelon_musk",
            content: "Sunday mornings are for slow breakfasts and even slower jazz. 🎷🥐 There's something magical about the quiet before the city wakes up.",
            createdAt: "10m ago",
            numOfComments: 3,
            numOfHearts: 1345,
            avatar: "/elon.jpeg"
        },
        {
            username: "trumpet",
            content: "Finally hiked the trail I've been looking at for months. The view from the top was absolutely worth the struggle. Nature has a way of resetting your perspective.",
            createdAt: "3h ago",
            numOfComments: 15,
            numOfHearts: 6,
            avatar: "/trump.jpeg"
        },
        {
            username: "kimpossible",
            content: "Does anyone else feel like time is moving exceptionally fast lately? I swear it was January just yesterday.",
            createdAt: "5h ago",
            numOfComments: 42,
            numOfHearts: 145,
            avatar: "/kim.jpeg"
        },
        {
            username: "Xi_aomi",
            content: "Small wins matter. Fixed a bug that's been bugging me (pun intended) for a week. Celebrating with a donut.",
            createdAt: "1d ago",
            numOfComments: 7,
            numOfHearts: 12,
            avatar: "/xi.jpeg"
        }
    ];

    return (
        <div className="feed-container">
            <div className="feed-header">
                <h1 className="feed-title">Friends Feed</h1>
                <p className="feed-subtitle">Updates from your friends</p>
            </div>

            <div className="flex flex-col">
                {posts.map((post, i) => (
                    <div key={i} className="post-card">
                        {/* Avatar Column */}
                        <div className="post-avatar-container">
                            <img src={post.avatar} alt="Post Avatar" className="post-avatar" />
                        </div>

                        {/* Content Column */}
                        <div className="post-content-container">
                            {/* Header */}
                            <div className="post-header">
                                <h3 className="post-username">
                                    @{post.username}
                                </h3>
                                <span className="post-timestamp">{post.createdAt}</span>
                            </div>

                            {/* Content */}
                            <p className="post-text">
                                {post.content}
                            </p>

                            {/* Footer / Actions */}
                            <div className="post-actions">
                                {/* Reaction Button (Only one as requested) */}
                                <button className="action-btn action-btn-heart group/heart">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="icon-heart">
                                        <path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" />
                                    </svg>
                                    <span className="text-sm font-medium">{post.numOfHearts}</span>
                                </button>

                                {/* Comments */}
                                <button className="action-btn action-btn-comment group/comment">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="icon-comment">
                                        <path strokeLinecap="round" strokeLinejoin="round" d="M12 20.25c4.97 0 9-3.694 9-8.25s-4.03-8.25-9-8.25S3 7.444 3 12c0 2.104.859 4.023 2.273 5.48.432.447.74 1.04.586 1.641a4.483 4.483 0 0 1-.923 1.785A5.969 5.969 0 0 0 6 21c1.282 0 2.47-.402 3.445-1.087.81.22 1.668.337 2.555.337Z" />
                                    </svg>
                                    <span className="text-sm font-medium">{post.numOfComments}</span>
                                </button>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}