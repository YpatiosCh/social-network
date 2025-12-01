package application

import "context"

func (s *PostsService) hasRightToView(ctx context.Context, req hasRightToView) (bool, error) {
	// get the requester id, the parent entity id (so post or event even if the request is for comments)
	// user ids the requester follows and group ids the requester belongs to
	// group and post audience=group: only members can see
	// post audience=everyone: everyone can see (can we check this before all the fetches from users?)
	// post audience=followers: requester can see if they follow creator
	// post audience=selected: requester can see if they are in post audience table
	return false, nil
}
