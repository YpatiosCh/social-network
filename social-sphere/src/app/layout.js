import "./globals.css";

export const metadata = {
  title: "SocialSphere",
  description: "SocialSphere",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
