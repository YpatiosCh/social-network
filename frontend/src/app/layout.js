import './globals.css'
import Footer from '@/components/layout/Footer'

const RootLayout = ({ children }) => {
  return (
    <html>
      <body>
        {children}
      <div className='bg-black'>
        <Footer />
      </div>
      </body>
    </html>
  )
}

export default RootLayout;