package application

import "context"

//FRONT: return full event or just id?
func (s *PostsService) CreateEvent(ctx context.Context, req CreateEventReq) (Event, error) {
	//check creator is member of group (API Gateway)
	return Event{}, nil
}

func (s *PostsService) DeleteEvent(ctx context.Context, req GenericReq) error {
	//check requester is creator of event (and member of the group? what happens if they're not any more?)
	return nil
}

//FRONT: return full event or just error?
func (s *PostsService) EditEvent(ctx context.Context, req EditEventReq) (Event, error) {
	//check requester is creator of event (and member of the group? what happens if they're not any more?)
	return Event{}, nil
}

func (s *PostsService) GetEventsByGroupId(ctx context.Context, req GenericPaginatedReq) ([]Event, error) {
	//check requester is member of group (API Gateway)
	return nil, nil
}

//FRONT: After a response the going/not going count has changed. Do you want me to return the whole event, just the count, or just error?
func (s *PostsService) RespondToEvent(ctx context.Context, req RespondToEventReq) error {
	//check requester is member of group (API Gateway)
	return nil
}

//FRONT: Same as above. Count changes, what do you need me to return?
func (s *PostsService) RemoveEventResponse(ctx context.Context, rec GenericReq) error {
	//check requester is member of group (API Gateway)
	return nil
}

//IF needed, should run periodically as go routine
func (s *PostsService) updateStillValid(ctx context.Context) error {
	//checks event date to today's date and changes stillValid bool for expired events
	return nil
}
