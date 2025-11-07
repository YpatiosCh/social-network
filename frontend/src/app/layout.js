import "./globals.css";
import BaseNav from "../components/layout/BaseNav"

export const metadata = {
  title: "SocialSphere - Welcome",
};

export default function RootLayout({ children }) {

  let user = ""

  return (
    <html lang="en">
      <body>
        <BaseNav user={user}/>
        {children}
        </body>
    </html>
  );
}
