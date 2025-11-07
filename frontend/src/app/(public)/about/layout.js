import "../../globals.css";

export const metadata = {
  title: "SocialSphere - About",
};

export default function RootLayout({ children }) {
  return (
    <main>
      <div>{children}</div>
    </main>
  );
}