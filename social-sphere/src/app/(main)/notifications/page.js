import { getNotifs } from "@/actions/notifs/get-user-notifs";
import NotificationsContent from "@/components/notifications/NotificationsContent";

export const metadata = {
    title: "Notifications",
};

export default async function NotificationsPage() {
    const initialNotifications = await getNotifs({ limit: 20, offset: 0 });

    return <NotificationsContent initialNotifications={initialNotifications} />;
}
