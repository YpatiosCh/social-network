package mapping

import (
	pb "social-network/shared/gen-go/chat"
	md "social-network/shared/go/models"
)

func MapPMToProto(m md.PM) *pb.PrivateMessage {
	return &pb.PrivateMessage{
		Id:             m.Id.Int64(),
		ConversationId: m.ConversationID.Int64(),
		Sender:         MapUserToProto(m.Sender),
		MessageText:    string(m.MessageText),
		CreatedAt:      m.CreatedAt.ToProto(),
		UpdatedAt:      m.UpdatedAt.ToProto(),
		DeletedAt:      m.DeletedAt.ToProto(),
	}
}

func MapGetPMsResp(res md.GetPMsResp) *pb.GetPrivateMessagesResponse {
	msgs := make([]*pb.PrivateMessage, 0, len(res.Messages))
	for _, m := range res.Messages {
		msgs = append(msgs, MapPMToProto(m))
	}

	return &pb.GetPrivateMessagesResponse{
		HaveMore: res.HaveMore,
		Messages: msgs,
	}
}
