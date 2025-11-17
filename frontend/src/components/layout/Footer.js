import Link from 'next/link'

export default function Footer() {
  return (
    <footer className="footer">
      <div className="footer-container">
        <div className="footer-grid">
          {/* Brand Column */}
          <div className="footer-brand">
            <div className="footer-logo">SocialSphere</div>
            <p className="footer-tagline">
              Connect with friends, share moments, and build communities. The social network that brings people together.
            </p>
          </div>

          {/* Product Links */}
          <div className="footer-column">
            <h3 className="footer-title">Product</h3>
            <Link href="/#features" className="footer-link">Features</Link>
            <Link href="/mobile" className="footer-link">Web App</Link>
            <Link href="/desktop" className="footer-link">Desktop App</Link>
          </div>

          {/* Company Links */}
          <div className="footer-column">
            <h3 className="footer-title">Company</h3>
            <Link href="/about" className="footer-link">About Us</Link>
            <Link href="/careers" className="footer-link">Careers</Link>
            <Link href="/contact" className="footer-link">Contact</Link>
          </div>

          {/* Legal Links */}
          <div className="footer-column">
            <h3 className="footer-title">Legal</h3>
            <Link href="/privacy" className="footer-link">Privacy Policy</Link>
            <Link href="/terms" className="footer-link">Terms of Service</Link>
            <Link href="/cookies" className="footer-link">Cookie Policy</Link>
            <Link href="/guidelines" className="footer-link">Community Guidelines</Link>
          </div>
        </div>

        {/* Bottom Bar */}
        <div className="footer-bottom items-center">
          <p className="footer-copyright items-center">
            © 2025 SocialSphere. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  )
}
