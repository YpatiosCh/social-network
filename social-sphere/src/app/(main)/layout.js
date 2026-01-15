import Navbar from "@/components/layout/Navbar";
import LiveSocketWrapper from "@/components/providers/LiveSocketWrapper";

export const dynamic = 'force-dynamic';

export default function MainLayout({ children }) {
    return (
        <LiveSocketWrapper>
            <div className="min-h-screen flex flex-col bg-(--muted)/6">
                <Navbar />
                <main className="flex-1 w-full">
                    {children}
                </main>
            </div>
        </LiveSocketWrapper>
    );
}