package userservice

import (
	"context"
	"errors"
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
		return 0, errors.New("invalid date format: expected YYYY-MM-DD")
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

	//start transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	qtx := s.db.(*sqlc.Queries).WithTx(tx)

	// Insert user
	userID, err := qtx.InsertNewUser(ctx, sqlc.InsertNewUserParams{
		Username:      req.Username,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		DateOfBirth:   dob,
		Avatar:        req.Avatar,
		AboutMe:       req.About,
		ProfilePublic: req.Public,
	})
	if err != nil {
		return 0, err
	}

	// Insert auth
	if err := qtx.InsertNewUserAuth(ctx, sqlc.InsertNewUserAuthParams{
		UserID:       userID,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Salt:         "", //I think we can remove this
	}); err != nil {
		return 0, err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return UserId(userID), nil

}

func LoginUser() {
	//called with: username (and/or email?), password
	//returns user_id,username,avatar,profile_public, or error
	//---------------------------------------------------------------------
	//BY EMAIL OR ONLY USERNAME???? (discuss with front)
	//GetUserForLogin (id, username, password hash, salt), check status=active (maybe also email TODO)
	//check password correct
	//if failed login IncrementFailedLoginAttempts
	//if success ResetFailedLoginAttempts
	//issue token (user service or api gateway?)
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

func GetUserProfile() {
	//called with: user_id, viewer_id (to check permission to view)
	//returns id, username, first_name, last_name, date_of_birth, avatar, about_me, profile_public, number of followers, number of following, list of groups
	//---------------------------------------------------------------------
	// check if user has permission to see (public profile or isFollower)
	// GetUserProlife(id)

	// number of followers, following (TODO keep in profile with trigger?)
	// number of groups? (TODO keep in profile with trigger?)
	// which groups

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

func GetAllGroups() {
	//called with nothing
	//returns list of groups containing group_id, group_title, group_description, members_count
	//---------------------------------------------------------------------
	//GetAllGroups
}

func GetUserGroups() {
	//called with user_id
	//returns list of groups containing group_id, group_title, group_description, members_count
	//---------------------------------------------------------------------

	//GetUserGroups
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
