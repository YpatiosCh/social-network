import { getConv } from "@/actions/chat/get-conv";
import { getMessages } from "@/actions/chat/get-messages";
import MessagesContent from "@/components/messages/MessagesContent";
import { getProfileInfo } from "@/actions/profile/get-profile-info";

export default async function ConversationPage({ params }) {
    const { id } = await params;

    // Fetch conversations list
    const convsResult = await getConv({ first: true, limit: 50 });
    let conversations = convsResult.success ? convsResult.data : [];

    // Find the selected conversation from the list
    const selectedConversation = conversations.find(
        (conv) => conv.Interlocutor?.id === id
    );

    // Fetch messages for the selected conversation if found
    let initialMessages = [];
    if (selectedConversation) {
        const messagesResult = await getMessages({
            interlocutorId: selectedConversation.Interlocutor?.id,
            limit: 50,
        });
        if (messagesResult.success && messagesResult.data?.Messages) {
            console.log("Not calling for data");
            // Messages come newest first, reverse for display
            initialMessages = messagesResult.data.Messages.reverse();
        }
    } else {
        console.log("calling for data");
        const user = await getProfileInfo(id);
        const newConv = {
            Interlocutor: {
                id: user.user_id,
                username: user.username,
                avatar_url: user.avatar_url
            }
        }

        conversations = [newConv, ...conversations];

        console.log("USER JUST FETCHED: ", user);
        console.log("NEW CONV: ", newConv);
        console.log("CONVERSATIONS: ", conversations);
    }

    console.log("INITIAL MSGS: ", initialMessages);

    return (
        <MessagesContent
            initialConversations={conversations}
            initialSelectedId={id}
            initialMessages={initialMessages}
        />
    );
}
