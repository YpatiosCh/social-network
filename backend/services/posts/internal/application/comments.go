package application

import (
	"context"
	"fmt"
	ds "social-network/services/posts/internal/db/dbservice"
	"social-network/shared/gen-go/media"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	"social-network/shared/go/models"
	tele "social-network/shared/go/telemetry"
)

func (s *Application) CreateComment(ctx context.Context, req models.CreateCommentReq) (err error) {
	errMsg := fmt.Sprintf("create comment: req: %#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		tele.Error(ctx, "Posts Service: Request to create comment failed validation: @1", "request", req)
		return ce.Wrap(ce.ErrInvalidArgument, err, errMsg)
	}

	accessCtx := accessContext{
		requesterId: req.CreatorId.Int64(),
		entityId:    req.ParentId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		tele.Error(ctx, "Posts Service: hasRightToView for create comment failed with accessCtx @1 ", "access ctx", accessCtx)
		return ce.Wrap(ce.ErrInternal, err, errMsg+": hasRightToView").WithPublic("posts service error")
	}
	if !hasAccess {
		tele.Error(ctx, "Posts Service: Requester @1 has no permission to create comment for post @2", "requester id", req.CreatorId, "post id", req.ParentId)
		return ce.Wrap(ce.ErrPermissionDenied, fmt.Errorf("user has no permission to view or edit this entity"), errMsg).WithPublic("permission denied")
	}
	var commentId int64
	err = s.txRunner.RunTx(ctx, func(q *ds.Queries) error {
		commentId, err = q.CreateComment(ctx, ds.CreateCommentParams{
			CommentCreatorID: req.CreatorId.Int64(),
			ParentID:         req.ParentId.Int64(),
			CommentBody:      req.Body.String(),
		})

		if err != nil {
			tele.Error(ctx, "Posts Service: db create comment failed for req @1 ", "request", req)
			return ce.Wrap(ce.ErrInternal, err, errMsg+": db: create comment").WithPublic("posts service error")
		}

		if req.ImageId != 0 {
			err = q.UpsertImage(ctx, ds.UpsertImageParams{
				ID:       req.ImageId.Int64(),
				ParentID: commentId,
			})
			if err != nil {
				tele.Error(ctx, "Posts Service: db upsert image for comment @1 failed for image id @2 ", "comment id", commentId, "image id", req.ImageId)
				return ce.Wrap(ce.ErrInternal, err, errMsg+": db: upsert image").WithPublic("posts service error")
			}
		}
		return nil
	})
	if err != nil {
		tele.Error(ctx, "Posts Service: create comment tx failed for req @1 ", "request", req)
		return ce.Wrap(nil, err, errMsg)
	}

	//create notification
	userMap, err := s.userRetriever.GetUsers(ctx, ct.Ids{req.CreatorId})
	if err != nil {
		tele.Error(ctx, "Posts Service: create comment: get users failed for user @1 ", "user id", req.CreatorId)
		return nil //return with no error but without creating non-essential notif
	}
	var commenterUsername string
	if u, ok := userMap[req.CreatorId]; ok {
		commenterUsername = u.Username.String()
	}
	basicPost, err := s.db.GetBasicPostByID(ctx, req.ParentId.Int64())
	if err != nil {
		tele.Error(ctx, "Posts Service: get basic post by id failed for post @1 ", "post id", req.ParentId)
		return nil //return with no error but without creating non-essential notif
	}
	err = s.clients.CreatePostComment(ctx, basicPost.CreatorID, req.CreatorId.Int64(), req.ParentId.Int64(), commenterUsername, req.Body.String())
	if err != nil {
		tele.Error(ctx, "Posts Service: create post comment notification failed for comment @1 ", "comment id", commentId)
		return nil //return with no error but without creating non-essential notif
	}
	return nil
}

func (s *Application) EditComment(ctx context.Context, req models.EditCommentReq) error {
	errMsg := fmt.Sprintf("edit comment: req: %#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		tele.Error(ctx, "Posts Service: Request to edit comment failed validation: @1", "request", req)
		return ce.Wrap(ce.ErrInvalidArgument, err, errMsg)
	}

	accessCtx := accessContext{
		requesterId: req.CreatorId.Int64(),
		entityId:    req.CommentId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		tele.Error(ctx, "Posts Service: hasRightToView for create comment failed with accessCtx @1 ", "access ctx", accessCtx)
		return ce.Wrap(ce.ErrInternal, err, errMsg+": hasRightToView").WithPublic("posts service error")
	}
	if !hasAccess {
		tele.Error(ctx, "Posts Service: Requester @1 has no permission to edit comment @2", "requester id", req.CreatorId, "post id", req.CommentId)
		return ce.Wrap(ce.ErrPermissionDenied, fmt.Errorf("user has no permission to view or edit this entity"), errMsg).WithPublic("permission denied")
	}

	return s.txRunner.RunTx(ctx, func(q *ds.Queries) error {
		rowsAffected, err := q.EditComment(ctx, ds.EditCommentParams{
			CommentBody:      req.Body.String(),
			ID:               req.CommentId.Int64(),
			CommentCreatorID: req.CreatorId.Int64(),
		})
		if err != nil {
			tele.Error(ctx, "Posts Service: db edit comment failed for req @1 ", "request", req)
			return ce.Wrap(ce.ErrInternal, err, errMsg+": db: edit comment").WithPublic("posts service error")
		}
		if rowsAffected != 1 {
			tele.Error(ctx, "Posts Service: db create comment found no comment fitting the criteria req @1 ", "request", req)
			return ce.Wrap(ce.ErrNotFound, err, errMsg+": db: create comment").WithPublic("posts service error")

		}
		if req.ImageId > 0 {
			err := q.UpsertImage(ctx, ds.UpsertImageParams{
				ID:       req.ImageId.Int64(),
				ParentID: req.CommentId.Int64(),
			})
			if err != nil {
				tele.Error(ctx, "Posts Service: db upsert image for comment @1 failed for image id @2 ", "comment id", req.CommentId, "image id", req.ImageId)
				return ce.Wrap(ce.ErrInternal, err, errMsg+": db: upsert image").WithPublic("posts service error")
			}
		}
		if req.DeleteImage {
			rowsAffected, err := q.DeleteImage(ctx, req.CommentId.Int64())
			if err != nil {
				tele.Error(ctx, "Posts Service: db delete image for comment @1 failed for image id @2 ", "comment id", req.CommentId, "image id", req.ImageId)
				return ce.Wrap(ce.ErrInternal, err, errMsg+": db: delete image").WithPublic("posts service error")

			}
			if rowsAffected != 1 {
				tele.Warn(ctx, "EditComment: image @1 for comment @2 could not be deleted: not found.", "image id", req.ImageId, "comment id", req.CommentId)
			}
		}
		return nil
	})

}

func (s *Application) DeleteComment(ctx context.Context, req models.GenericReq) error {
	errMsg := fmt.Sprintf("delete comment: req: %#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		tele.Error(ctx, "Posts Service: Request to delete comment failed validation: @1", "request", req)
		return ce.Wrap(ce.ErrInvalidArgument, err, errMsg)
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		tele.Error(ctx, "Posts Service: hasRightToView for delete comment failed with accessCtx @1 ", "access ctx", accessCtx)
		return ce.Wrap(ce.ErrInternal, err, errMsg+": hasRightToView").WithPublic("posts service error")

	}
	if !hasAccess {
		tele.Error(ctx, "Posts Service: Requester @1 has no permission to delete comment @2", "requester id", req.RequesterId, "comment id", req.EntityId)
		return ce.Wrap(ce.ErrPermissionDenied, fmt.Errorf("user has no permission to view or edit this entity"), errMsg).WithPublic("permission denied")

	}

	rowsAffected, err := s.db.DeleteComment(ctx, ds.DeleteCommentParams{
		ID:               req.EntityId.Int64(),
		CommentCreatorID: req.RequesterId.Int64(),
	})
	if err != nil {
		tele.Error(ctx, "Posts Service: db delete comment failed for req @1 ", "request", req)
		return ce.Wrap(ce.ErrInternal, err, errMsg+": db: delete comment").WithPublic("posts service error")

	}
	if rowsAffected != 1 {
		tele.Warn(ctx, "Posts Service: comment @1 could not be deleted: not found.", "comment id", req.EntityId)
	}
	return nil
}

func (s *Application) GetCommentsByParentId(ctx context.Context, req models.EntityIdPaginatedReq) ([]models.Comment, error) {
	errMsg := fmt.Sprintf("get comments by parent id: req: %#v", req)

	if err := ct.ValidateStruct(req); err != nil {
		tele.Error(ctx, "Posts Service: Request to get comments by parent id failed validation: @1", "request", req)
		return nil, ce.Wrap(ce.ErrInvalidArgument, err, errMsg)
	}

	accessCtx := accessContext{
		requesterId: req.RequesterId.Int64(),
		entityId:    req.EntityId.Int64(),
	}

	hasAccess, err := s.hasRightToView(ctx, accessCtx)
	if err != nil {
		tele.Error(ctx, "Posts Service: hasRightToView for comments for post @1 failed with accessCtx @2 ", "post id", req.EntityId, "access ctx", accessCtx)
		return nil, ce.Wrap(ce.ErrInternal, err, errMsg+": hasRightToView").WithPublic("posts service error")
	}
	if !hasAccess {
		tele.Error(ctx, "Posts Service: Requester @1 has no permission to view comments for post @2", "requester id", req.RequesterId, "post id", req.EntityId)
		return nil, ce.Wrap(ce.ErrPermissionDenied, fmt.Errorf("user has no permission to view comments of this entity"), errMsg).WithPublic("permission denied")

	}

	rows, err := s.db.GetCommentsByPostId(ctx, ds.GetCommentsByPostIdParams{
		ParentID: req.EntityId.Int64(),
		UserID:   req.RequesterId.Int64(),
		Limit:    req.Limit.Int32(),
		Offset:   req.Offset.Int32(),
	})
	if err != nil {
		tele.Error(ctx, "Posts Service: dbget comments by post id failed for req @1 ", "request", req)
		return nil, ce.Wrap(ce.ErrInternal, err, errMsg+": db: get comments by post id").WithPublic("posts service error")

	}
	comments := make([]models.Comment, 0, len(rows))
	userIDs := make(ct.Ids, 0, len(rows))
	commentImageIds := make(ct.Ids, 0, len(rows))

	for _, r := range rows {
		uid := ct.Id(r.CommentCreatorID)
		userIDs = append(userIDs, uid)

		comments = append(comments, models.Comment{
			CommentId: ct.Id(r.ID),
			ParentId:  req.EntityId,
			Body:      ct.CommentBody(r.CommentBody),
			User: models.User{
				UserId: ct.Id(r.CommentCreatorID),
			},
			ReactionsCount: int(r.ReactionsCount),
			CreatedAt:      ct.GenDateTime(r.CreatedAt.Time),
			UpdatedAt:      ct.GenDateTime(r.UpdatedAt.Time),
			LikedByUser:    r.LikedByUser,
			ImageId:        ct.Id(r.Image),
		})
		if r.Image > 0 {
			commentImageIds = append(commentImageIds, ct.Id(r.Image))
		}
	}

	if len(comments) == 0 {
		return comments, nil
	}

	userMap, err := s.userRetriever.GetUsers(ctx, userIDs)
	if err != nil {
		tele.Error(ctx, "Posts Service: get comments by parent id: get users failed for user ids @1 ", "user ids", userIDs)
		return nil, ce.Wrap(ce.ErrInternal, err, errMsg+": user retriever").WithPublic("error retrieving users info")
	}

	var imageMap map[int64]string
	if len(commentImageIds) > 0 {
		imageMap, _, err = s.mediaRetriever.GetImages(ctx, commentImageIds, media.FileVariant_MEDIUM)
	}
	if err != nil {
		tele.Error(ctx, "Posts Service: get comments by parent id: get images failed for image ids @1 ", "image ids", commentImageIds)
		return nil, ce.Wrap(ce.ErrInternal, err, errMsg+": user retriever").WithPublic("error retrieving users info")

	}

	for i := range comments {
		uid := comments[i].User.UserId
		if u, ok := userMap[uid]; ok {
			comments[i].User = u
		}
		comments[i].ImageUrl = imageMap[comments[i].ImageId.Int64()]
	}

	return comments, nil
}
