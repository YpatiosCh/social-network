'use client'

const Home = () => {
  return (
    <div className="landing-container">
      {/* Grain texture overlay */}
      <div className="grain-texture"></div>

      {/* Hero */}
      <section className="hero-section">
        <div className="hero-grid-full">
          {/* Main glass surface - Full Width */}
          <div className="glass-card-wrapper">
            <div className="glass-glow"></div>
            <div className="glass-card-main">
              <div className="badge-pill">
                SocialSphere
              </div>

              <h1 className="hero-heading">
                Connect<span className="hero-gradient-text">.</span><br />
                Share<span className="hero-gradient-text">.</span><br />
                <span className="hero-gradient-text">Thrive</span>.
              </h1>

              <p className="hero-description">
                Real conversations. Real connections. SocialSphere brings people together through posts, groups, events, and live chat—all in one place.
              </p>

              <div className="hero-cta-group">
                <button className="btn-primary">
                  <span>Get Started →</span>
                </button>

                <button className="btn-secondary">
                  Explore features
                </button>
              </div>
            </div>

            {/* Floating accent line */}
            <div className="floating-accent-line"></div>
          </div>
        </div>
      </section>

      {/* Bento grid - asymmetric chaos */}
      <section className="section-bento">
        <div className="bento-container">
          <div className="bento-grid">
            {/* Large */}
            <div className="bento-card bento-large">
              <div className="bento-card-glow glow-blue"></div>
              <div className="bento-card-content border-left-top">
                <div className="bento-content-flex-col">
                  <div>
                    <div className="bento-icon">
                      👥
                    </div>
                    <h3 className="bento-heading-large">
                      Follow, post,<br/>connect
                    </h3>
                    <p className="bento-description">
                      Share updates with followers, create private or public posts, and build your network organically.
                    </p>
                  </div>
                  <div className="bento-accent-line"></div>
                </div>
              </div>
            </div>

            {/* Medium tall */}
            <div className="bento-card bento-medium">
              <div className="bento-card-glow glow-cyan"></div>
              <div className="bento-card-content border-right-top">
                <div className="emoji-icon">💬</div>
                <h3 className="bento-heading-medium">Live messaging</h3>
                <p className="bento-description">Real-time chat with friends and group conversations. Stay connected instantly.</p>
              </div>
            </div>

            {/* Small */}
            <div className="bento-card bento-small">
              <div className="bento-card-glow glow-indigo"></div>
              <div className="bento-card-content border-bottom-left bento-content-center">
                <div className="text-center">
                  <div className="emoji-icon-large">🔐</div>
                  <div className="bento-label">Private</div>
                </div>
              </div>
            </div>

            {/* Wide short */}
            <div className="bento-card bento-wide">
              <div className="bento-card-glow glow-cyan"></div>
              <div className="bento-card-content border-bottom bento-content-flex">
                <div>
                  <div className="emoji-icon-small">👨‍👩‍👧‍👦</div>
                  <h3 className="bento-heading-small">Groups & Events</h3>
                </div>
              </div>
            </div>

            {/* Tiny */}
            <div className="bento-card bento-tiny">
              <div className="bento-card-glow glow-teal"></div>
              <div className="bento-card-content border-right-bottom bento-content-between">
                <span className="bento-label">MOBILE • DESKTOP</span>
                <div className="emoji-icon-small">📱</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Benefits - diagonal stack */}
      <section className="section-benefits">
        <div className="benefits-container">
          {/* 01 */}
          <div className="benefit-item">
            <div className="benefit-number benefit-number-left">
              01
            </div>
            <div className="benefit-content benefit-content-left">
              <div className="benefit-glow benefit-glow-blue"></div>
              <div className="benefit-card benefit-card-left">
                <div className="benefit-accent-line benefit-accent-blue"></div>
                <h3 className="benefit-heading">
                  Share your world
                </h3>
                <p className="benefit-description">
                  Post updates, photos, and thoughts with friends or the world. Control who sees what with flexible privacy settings—public, friends-only, or custom audiences.
                </p>
              </div>
            </div>
          </div>

          {/* 02 */}
          <div className="benefit-item">
            <div className="benefit-number benefit-number-right">
              02
            </div>
            <div className="benefit-content benefit-content-right">
              <div className="benefit-glow benefit-glow-cyan"></div>
              <div className="benefit-card benefit-card-right">
                <div className="benefit-accent-line benefit-accent-cyan"></div>
                <h3 className="benefit-heading benefit-heading-right">
                  Build communities
                </h3>
                <p className="benefit-description benefit-description-right">
                  Create groups around shared interests. Host events, plan meetups, vote on decisions. Real communities with real people, not just endless scrolling.
                </p>
              </div>
            </div>
          </div>

          {/* 03 */}
          <div className="benefit-item">
            <div className="benefit-number benefit-number-left">
              03
            </div>
            <div className="benefit-content benefit-content-left">
              <div className="benefit-glow benefit-glow-indigo"></div>
              <div className="benefit-card benefit-card-left">
                <div className="benefit-accent-line benefit-accent-indigo"></div>
                <h3 className="benefit-heading">
                  Stay in touch
                </h3>
                <p className="benefit-description">
                  Chat with friends one-on-one or in groups. Get instant notifications when someone follows you, comments on your posts, or invites you to events. Never miss a beat.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA - brutalist */}
      <section className="section-cta">
        <div className="cta-container">
          <div className="cta-wrapper">
            <div className="cta-glow"></div>
            <div className="cta-card">
              <div className="cta-content">

                <h2 className="cta-heading">
                  Ready to<br/>
                  <span className="hero-gradient-text">connect?</span>
                </h2>

                <p className="cta-description">
                  Join a community where you can share, discover, and engage. Follow friends, create groups, host events, and chat in real-time—all in one place.
                </p>

                <div className="cta-button-wrapper">
                  <button className="btn-cta">
                    <span>Create your account</span>
                  </button>

                  <div className="cta-subtext">
                    Free to join • Takes less than a minute
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

    </div>
  )
};

export default Home;
