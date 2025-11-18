'use client'

import { useState } from 'react'

const LoginPage  = () => {
  const [formData, setFormData] = useState({
    loginName: '',
    email: '',
    Name: '',
    Surname: '',
    Phone: '',
    NickName: '',
    BirthDate: '',
    Gender: '',
    Address: '',
    City: '',
    Country: '',
    ZipCode: '',
    password: '',
    confirmPassword: '',
  })

  const handleLogin = async (event) => {
    event.preventDefault()

    const requiredFields = ['loginName', 'password', 'confirmPassword']
    const missingFields = requiredFields.filter((field) => !formData[field].trim())

    if (missingFields.length > 0) {
      window.alert('Please fill in all required fields.')
      return
    }

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

            <form className="mt-8 space-y-6 text-left" onSubmit={handleLogin}>
              <div className="grid grid-cols-2 gap-6 md:grid-cols-3">
                <label className="flex flex-col text-sm text-white md">
                  <span className="font-semibold">Login Name <span className="text-red-500">*</span></span>
                  <input
                    type="text"
                    name="loginName"
                    value={formData.loginName}
                    placeholder="Select login name"
                    required
                    onChange={(event) => setFormData({ ...formData, loginName: event.target.value })}
                    className="mt-2 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white">
                  <span className="font-semibold">Email<span className="text-red-500">*</span></span>
                  <input
                    type="email"
                    name="email"
                    value={formData.email}
                    placeholder="Enter email"
                    onChange={(event) => setFormData({ ...formData, email: event.target.value })}
                    className="mt-2 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white">
                  <span className="font-semibold">Name</span>
                  <input
                    type="text"
                    name="Name"
                    value={formData.Name}
                    placeholder="Enter name"
                    onChange={(event) => setFormData({ ...formData, Name: event.target.value })}
                    className="mt-2 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white">
                  <span className="font-semibold">Surname</span>
                  <input
                    type="text"
                    name="Surname"
                    value={formData.Surname}
                    placeholder="Enter surname"
                    onChange={(event) => setFormData({ ...formData, Surname: event.target.value })}
                    className="mt-2 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white">
                  <span className="font-semibold">Phone</span>
                  <input
                    type="tel"
                    name="Phone"
                    value={formData.Phone}
                    placeholder="Enter phone number"
                    onChange={(event) => setFormData({ ...formData, Phone: event.target.value })}
                    className="mt-2 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white">
                  <span className="font-semibold">Nickname</span>
                  <input
                    type="text"
                    name="NickName"
                    value={formData.NickName}
                    placeholder="Enter nickname"
                    onChange={(event) => setFormData({ ...formData, NickName: event.target.value })}
                    className="mt-2 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white">
                  <span className="font-semibold">Birth Date</span>
                  <input
                    type="date"
                    name="BirthDate"
                    value={formData.BirthDate}
                    onChange={(event) => setFormData({ ...formData, BirthDate: event.target.value })}
                    className="mt-2 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white">
                  <span className="font-semibold">Gender</span>
                  <select
                    name="Gender"
                    value={formData.Gender}
                    onChange={(event) => setFormData({ ...formData, Gender: event.target.value })}
                    className="mt-2 w-32 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-sm text-white focus:border-white focus:outline-none"
                  >
                    <option value="" className="text-black">Select</option>
                    <option value="male" className="text-black">Male</option>
                    <option value="female" className="text-black">Female</option>
                    <option value="other" className="text-black">Other</option>
                  </select>
                </label>

                <label className="flex flex-col text-sm text-white md:col-span-3">
                  <span className="font-semibold">Address</span>
                  <textarea
                    name="Address"
                    value={formData.Address}
                    placeholder="Enter address"
                    rows="3"
                    onChange={(event) => setFormData({ ...formData, Address: event.target.value })}
                    className="mt-2 w-3/5 rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white md:col-span-1">
                  <span className="font-semibold">City</span>
                  <input
                    type="text"
                    name="City"
                    value={formData.City}
                    placeholder="Enter city"
                    onChange={(event) => setFormData({ ...formData, City: event.target.value })}
                    className="mt-2 w-full rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white md:col-span-1">
                  <span className="font-semibold">Zip Code</span>
                  <input
                    type="text"
                    name="ZipCode"
                    value={formData.ZipCode}
                    placeholder="Enter zip code"
                    onChange={(event) => setFormData({ ...formData, ZipCode: event.target.value })}
                    className="mt-2 w-full rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white md:col-span-1">
                  <span className="font-semibold">Country</span>
                  <input
                    type="text"
                    name="Country"
                    value={formData.Country}
                    placeholder="Enter country"
                    onChange={(event) => setFormData({ ...formData, Country: event.target.value })}
                    className="mt-2 w-full rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white md:col-span-1">
                  <span className="font-semibold">Password <span className="text-red-500">*</span></span>
                  <input
                    type="password"
                    name="password"
                    value={formData.password}
                    placeholder="Select password"
                    required
                    onChange={(event) => setFormData({ ...formData, password: event.target.value })}
                    className="mt-2 w-full rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>

                <label className="flex flex-col text-sm text-white md:col-span-1">
                  <span className="font-semibold">Confirm Password <span className="text-red-500">*</span></span>
                  <input
                    type="password"
                    name="confirmPassword"
                    value={formData.confirmPassword}
                    placeholder="Confirm password"
                    required
                    onChange={(event) => setFormData({ ...formData, confirmPassword: event.target.value })}
                    className="mt-2 w-full rounded-lg border border-white/40 bg-white/10 px-3 py-2 text-white placeholder-white/60 focus:border-white focus:outline-none"
                  />
                </label>
              </div>

              <div className="hero-cta-group mt-6 flex justify-center gap-4 py-2">
                <button className="btn-primary" type="submit">
                  <span> Register</span>
                </button>

                <button className="btn-secondary" type="button">
                  Exit
                </button>
              </div>
            </form>
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
