package userservice

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db   sqlc.Querier  // interface, can be *sqlc.Queries or mock
	pool *pgxpool.Pool // needed to start transactions
}

// NewUserService constructs a new UserService
func NewUserService(db sqlc.Querier, pool *pgxpool.Pool) *UserService {
	return &UserService{
		db:   db,
		pool: pool,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, req RegisterUserRequest) (UserId, error) {

	// convert date
	dobTime, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return 0, ErrInvalidDateFormat
	}

	dob := pgtype.Date{
		Time:  dobTime,
		Valid: true,
	}

	//hash password
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return 0, err
	}

	var newId int64

	err = s.runTx(ctx, func(q *sqlc.Queries) error {

		// Insert user
		userId, err := q.InsertNewUser(ctx, sqlc.InsertNewUserParams{
			Username:      req.Username,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			DateOfBirth:   dob,
			Avatar:        req.Avatar,
			AboutMe:       req.About,
			ProfilePublic: req.Public,
		})
		if err != nil {
			return err //TODO check how to return correct error
		}
		newId = userId

		// Insert auth
		return q.InsertNewUserAuth(ctx, sqlc.InsertNewUserAuthParams{
			UserID:       userId,
			Email:        req.Email,
			PasswordHash: passwordHash,
		})
	})

	if err != nil {
		return 0, err //TODO check how to return correct error
	}

	return UserId(newId), nil

}

func (s *UserService) LoginUser(ctx context.Context, req LoginReq) (User, error) {
	var u User

	err := s.runTx(ctx, func(q *sqlc.Queries) error {
		row, err := q.GetUserForLogin(ctx, req.Identifier)
		if err != nil {
			return err
		}

		u = User{
			UserId:   row.ID,
			Username: row.Username,
			Avatar:   row.Avatar,
			Public:   row.ProfilePublic,
		}

		if !checkPassword(row.PasswordHash, req.Password) {
			q.IncrementFailedLoginAttempts(ctx, row.ID)
			return err
		}
		q.ResetFailedLoginAttempts(ctx, u.UserId)
		return nil
	})

	if err != nil {
		return User{}, ErrWrongCredentials
	}

	//TODO what happens when eg failed login attempts > 3? Add logic?

	return u, nil
}

type BasicUserInfo struct {
	UserName      string
	Avatar        string
	PublicProfile bool
}

func GetBasicUserInfo(ctx context.Context, userID int64) (resp BasicUserInfo, err error) {
	//called with: user_id
	//returns username, avatar, profile_public(bool)
	//---------------------------------------------------------------------
	// GetUserBasic(id)
	return BasicUserInfo{UserName: "Mitsos", Avatar: "M", PublicProfile: true}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req UserProfileRequest) (UserProfileResponse, error) {
	var profile UserProfileResponse
	err := s.runTx(ctx, func(q *sqlc.Queries) error { //TODO consider not using a transaction for everything
		// TODO helper: check if user has permission to see (public profile or isFollower)
		row, err := q.GetUserProfile(ctx, req.UserId)
		if err != nil {
			return err
		}

		dob := time.Time{}
		if row.DateOfBirth.Valid {
			dob = row.DateOfBirth.Time
		}

		profile = UserProfileResponse{
			UserId:      row.ID,
			Username:    row.Username,
			FirstName:   row.FirstName,
			LastName:    row.LastName,
			DateOfBirth: dob,
			Avatar:      row.Avatar,
			About:       row.AboutMe,
			Public:      row.ProfilePublic,
		}
		profile.FollowersCount, err = q.GetFollowerCount(ctx, profile.UserId)
		if err != nil {
			return err
		}
		profile.FollowingCount, err = q.GetFollowingCount(ctx, profile.UserId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return UserProfileResponse{}, err
	}

	profile.Groups, err = s.GetUserGroups(ctx, profile.UserId)
	if err != nil {
		return UserProfileResponse{}, err
	}

	return profile, nil

	// THIS CAN BE HANDLED BY THE API GATEWAY (and different call from front):
	// from forum service get all posts paginated (and number of posts)
	// and within all posts check each one if viewer has permission
}

func UpdateUserProfile() {
	//called with user_id and any of: username (TODO), first_name, last_name, date_of_birth, avatar, about_me
	//returns full profile
	//request needs to come from same user
	//---------------------------------------------------------------------

	//UpdateUserProflie
	//TODO check how to not update all fields but only changes (nil pointers?)
}

func UpdateUserPassword() {
	//called with user_id, old password, new password_hash, salt
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------
	//UpdateUserPassword
}

func UpdateUserEmail() {
	//called with user_id, new email
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------
	//UpdateUserEmail
}

func (s *UserService) GetAllGroups(ctx context.Context) ([]Group, error) {
	rows, err := s.db.GetAllGroups(ctx)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0, len(rows))
	for _, r := range rows {
		groups = append(groups, Group{
			GroupId:          r.ID,
			GroupTitle:       r.GroupTitle,
			GroupDescription: r.GroupDescription,
			MembersCount:     r.MembersCount,
		})
	}

	return groups, nil
}

func (s *UserService) GetUserGroups(ctx context.Context, userId int64) ([]Group, error) {
	rows, err := s.db.GetUserGroups(ctx, userId)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0, len(rows))
	for _, r := range rows {
		groups = append(groups, Group{
			GroupId:          r.GroupID,
			GroupTitle:       r.GroupTitle,
			GroupDescription: r.GroupDescription,
			MembersCount:     r.MembersCount,
			Role:             r.Role,
		})
	}

	return groups, nil
}

func GetGroupInfo() {
	//called with group_id
	//returns group_id, group_title, group_description, members_count (owner?)
	//---------------------------------------------------------------------

	//GetGroupInfo

	//different calls for chat and posts (API GATEWAY)
}

func GetGroupMembers() {
	//called with group_id
	//returns list of members containing user_id, username, avatar, profile_public(bool), role(owner or member), joined_at
	//---------------------------------------------------------------------

	//getGroupMembers
}

func SeachByUsers() {
	//called with search term
	//returns list of users containing user_id, username, avatar, profile_public(bool)
	//---------------------------------------------------------------------

	//SearchUsers by username or name
}

func SearchGroup() {
	//called with search term
	//returns list of groups containing group_id, group_title, group_description, members_count
	//---------------------------------------------------------------------

	//SeachGroupsFuzzy
}

func GetFollowers() {
	//called with user_id
	//returns list of users containing user_id, username, avatar, profile_public(bool)
	//---------------------------------------------------------------------

	//GetFollowers() TODO FIX RETURNS
}

func GetFollowing() {
	//called with user_id
	//returns list of users containing user_id, username, avatar, profile_public(bool)
	//---------------------------------------------------------------------
	//GetFollowing()
}

func InviteToGroup() {
	//called with group_id,user_id
	//returns success or error
	//request needs to come from group owner or group member
	//---------------------------------------------------------------------

	//SendGroupInvite
}

func RequestJoinGroup() {
	//called with group_id,user_id
	//returns success or error
	//---------------------------------------------------------------------

	//SendGroupJoinRequest
}

func HandleGroupInvite() {
	//called with group_id,user_id, bool (accept or decline)
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------

	//yes or no
	//AcceptGroupInvite & addUserToGroup
	//DeclineGroupInvite
}

func CancelGroupInvite() {
	//called with group_id,user_id (who is invited), sender_id
	//returns success or error
	//request needs to come from sender
	//---------------------------------------------------------------------

	//CancelGroupInvite
}

func HandleGroupJoinRequest() {
	//called with group_id,user_id (who requested to join),owner_id(who responds), bool (accept or decline)
	//returns success or error
	//request needs to come from group owner
	//---------------------------------------------------------------------

	//yes or no
	//AcceptGroupJoinRequest & addUserToGroup
	//RejectGroupJoinRequest
}

func CancelGroupJoinRequest() {
	//called with group_id,user_id, bool (accept or decline)
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------

	//CancelGroupJoinRequest
}

func LeaveGroup() {
	//called with group_id,user_id
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------

	//initiated by user
	//LeaveGroup
}

func RemoveFromGroup() {
	//called with group_id,user_id (who is removed), owner_id(making the request)
	//returns success or error
	//request needs to come from group owner
	//---------------------------------------------------------------------

	//initiated by owner
	//LeaveGroup
}

func CreateGroup() {
	//called with owner_id, group_title, group_description
	//returns group_id
	//---------------------------------------------------------------------

	//CreateGroup
	//AddGroupOwnerAsMember
}

func DeleteGroup() { //low priorty
	//called with group_id, owner_id
	//returns success or error
	//request needs to come from owner
	//---------------------------------------------------------------------

	//initiated by ownder
	//SoftDeleteGroup
}

func TranferGroupOwnerShip() { //low priority
	//called with group_id,previous_owner_id, new_owner_id
	//returns success or error
	//request needs to come from previous owner (or admin - not implemented)
	//---------------------------------------------------------------------

}

// HashPassword hashes a password using bcrypt.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword compares a hashed password with a plain-text password.
func checkPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func issueToken() { //here or in gateway?

}

func checkToken() { //in shared?

}

func DeleteUser() { //low priority
	//called with: user_id
	//returns success or error
	//request needs to come from same user or admin (not implemented)
	//---------------------------------------------------------------------
	//softDeleteUser(id)
}

func BanUser() { //low priority
	//called with: user_id
	//returns success or error
	//request needs to come from admin (not implemented)
	//---------------------------------------------------------------------
	//BanUser(id, liftBanDate)
	//TODO logic to automatically unban after liftBanDate
}

func UnbanUser() { //low priority
	//called with: user_id
	//returns success or error
	//request needs to come from admin (not implemented) or automatic on expiration
	//---------------------------------------------------------------------
	//UnbanUser(id)
}
