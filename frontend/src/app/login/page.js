'use client'

import { useState } from 'react'

const LoginPage  = () => {
  const [formData, setFormData] = useState({
    loginName: '',
    password: '',
  })

  const handleLogin = async (event) => {
    event.preventDefault()

    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      })

      if (!response.ok) {
        throw new Error('Login failed')
      }

      const data = await response.json()
      console.log(data)
    } catch (error) {
      console.error('Login error', error)
    }
  }

  return(
  <div className="landing-container">
    {/* Grain texture overlay */}
    <div className="grain-texture"></div>
    <section className="hero-section">
      <div className="hero-grid-full">
        {/* Main glass surface - Full Width */}
        <div className="glass-card-wrapper">
          <div className="glass-glow"></div>
          <div className="glass-card-main">
            

            <h1 className="hero-heading text-center">
              Login to Connect<span className="hero-gradient-text">.</span><br />
             
            </h1>

            <form className="mt-6 text-center" onSubmit={handleLogin}>
                <input
                  type="text"
                  name="loginName"
                  value={formData.loginName}
                  placeholder="Enter your login name"
                  onChange={(event) => setFormData({ ...formData, loginName: event.target.value })}
                  className="mt-2 w-1/5 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-sm text-white placeholder-white/60 focus:border-white focus:outline-none"
                />
                <br />
                <input
                  type="password"
                  name="password"
                  value={formData.password}
                  placeholder="Enter your password"
                  onChange={(event) => setFormData({ ...formData, password: event.target.value })}
                  className="mt-2 w-1/5 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-sm text-white placeholder-white/60 focus:border-white focus:outline-none"
                />

            <div className="hero-cta-group mt-6 flex justify-center gap-4 py-2">
              <button className="btn-primary" type="submit">
                <span> Login</span>
              </button>

              <button className="btn-secondary">
                Exit
              </button>
            </div>
            </form>

            <div className="hero-cta-group mt-6 flex justify-center gap-4 py-2">
              <button className="btn-secondary">
                Not registered? <a href="/register">Go to Register</a>
              </button>
            </div>
          </div>

          {/* Floating accent line */}
          <div className="floating-accent-line"></div>
        </div>
      </div>
    </section>
  </div>
  )
}

export default LoginPage
