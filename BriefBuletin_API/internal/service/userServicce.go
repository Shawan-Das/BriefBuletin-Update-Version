package service

import (
	"context"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	auth "github.com/rest/api/internal/dbmodel/db_query"
	"github.com/rest/api/internal/util"

	"github.com/rest/api/internal/model"
)

// var _usLogger = logrus.New()

// /api/auth/create - create user
func (s *RESTService) createUser(c *gin.Context) APIResponse {
	var input model.CreateUserInput
	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	// Validate required fields
	if input.Email == "" || input.Password == "" || input.Phone == "" {
		return BuildResponse400("Email, password, and phone are required")
	}

	// Set default values
	userName := input.UserName
	if userName == "" {
		userName = input.Email // Use email as default username
	}

	// Check if user already exists - check all three unique fields
	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	// Check email uniqueness
	_, err := qtx.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return BuildResponse400("User with this email already exists")
	}

	// Check username uniqueness
	_, err = qtx.GetUserByUserName(ctx, userName)
	if err == nil {
		return BuildResponse400("User with this username already exists")
	}

	// Check phone uniqueness
	_, err = qtx.GetUserByPhone(ctx, input.Phone)
	if err == nil {
		return BuildResponse400("User with this phone number already exists")
	}
	role := input.Role
	if role == "" {
		role = "USER" // Default role
	}

	// Hash password
	hashedPassword := s.getHashOf(input.Password)
	// otp & exp
	otp := util.GenerateSixDigits()
	currentTime := time.Now()
	otpExp := currentTime.Add(5 * time.Minute)
	// Create user
	createParams := auth.CreateUserParams{
		UserName: userName,
		Email:    input.Email,
		Phone:    input.Phone,
		Pass:     hashedPassword,
		Otp:      otp,
		OtpExp:   ToPGTimestamp(otpExp),
		Role:     role,
	}

	err = qtx.CreateUser(ctx, createParams)
	if err != nil {
		_asLogger.Errorf("Error creating user: %v", err)
		return BuildResponse500("Failed to create user", err.Error())
	}
	// send otp mail
	// TODO: without error handling
	// request3rdPartyEmailAPI(input.Email, userName, otp, true)
	smtpService := SmtpService{}
	go smtpService.SetNewPassWordMail(input.Email, input.UserName, otp, true)

	return BuildResponse200("An OTP has been sent to your email. Please verify yourself", nil)
}

// /api/auth/login - login (supports username, email, or phone)
func (s *RESTService) validateLogin(c *gin.Context) APIResponse {
	var input model.LoginInput
	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	if input.Login == "" || input.Password == "" {
		return BuildResponse400("Login identifier (username/email/phone) and password are required")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	// Try to find user by username, email, or phone
	user, err := qtx.GetUserByLogin(ctx, input.Login)
	if err != nil {
		_asLogger.Errorf("Error getting user: %v", err)
		return BuildResponse400("Invalid login credentials or password")
	}
	// check if user is valid
	if !user.UserValid {
		otp := util.GenerateSixDigits()
		currentTime := time.Now()
		otpExp := currentTime.Add(5 * time.Minute)
		// request3rdPartyEmailAPI(user.Email, user.UserName, otp, true)
		smtpService := SmtpService{}
		go smtpService.SetNewPassWordMail(user.Email, user.UserName, otp, true)
		_ = qtx.SendNewOtp(ctx, auth.SendNewOtpParams{
			Otp:    otp,
			OtpExp: ToPGTimestamp(otpExp),
			Email:  user.Email,
		})
		parts := strings.Split(user.Email, "@")
		var mail string
		if len(parts) != 2 {
			mail = user.Email[:3] + "**********"
		} else {
			mail = parts[0][:2] + "********@" + parts[1]
		}
		return BuildResponse403("An OTP has been sent to your email "+mail+". Please verify yourself to login.", nil)
	}
	// Check password
	hashedPassword := s.getHashOf(input.Password)
	if user.Pass != hashedPassword {
		return BuildResponse404("Invalid login credentials or password", false)
	}

	// Check if password is valid
	if !user.PssValid {
		return BuildResponse400("Password is not valid. Please reset your password")
	}

	// Create JWT token
	jwtToken := s.createJWTToken(user.UserID, user.Email, user.UserName, user.Role)

	response := BuildResponse200("Login successful", map[string]interface{}{
		"user_id":   user.UserID,
		"user_name": user.UserName,
		"email":     user.Email,
		"role":      user.Role,
		"token":     jwtToken,
		// "phone":     user.Phone,
	})
	response.Token = &jwtToken

	return response
}

// /api/auth/resetpwd - reset password
func (s *RESTService) resetPassword(c *gin.Context) APIResponse {
	var input model.AuthDataInput
	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	if input.Email == "" || input.NewPassword == "" {
		return BuildResponse400("Email and new password are required")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	// Check if user exists
	_, err := qtx.GetUserByEmail(ctx, input.Email)
	if err != nil {
		_asLogger.Errorf("Error getting user: %v", err)
		return BuildResponse404("User not found", false)
	}

	// Hash new password
	hashedPassword := s.getHashOf(input.NewPassword)

	// Update password
	updateParams := auth.UpdatePasswordParams{
		Pass:     hashedPassword,
		PssValid: true,
		Email:    input.Email,
	}

	err = qtx.UpdatePassword(ctx, updateParams)
	if err != nil {
		_asLogger.Errorf("Error updating password: %v", err)
		return BuildResponse500("Failed to reset password", err.Error())
	}

	return BuildResponse200("Password reset successfully", nil)
}

func (s *RESTService) sendOtp(c *gin.Context) APIResponse {
	var i model.OtpVerify
	if !parseInput(c, &i) {
		return BuildResponse400("Invalid input provided")
	}
	if i.Login == "" {
		return BuildResponse400("please enter user name/email/phone number of your account")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	user, err := qtx.GetUserByLogin(ctx, i.Login)
	if user.Email == "" {
		return BuildResponse404("User not found", false)
	} else if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}

	parts := strings.Split(user.Email, "@")
	var mail string
	if len(parts) != 2 {
		mail = user.Email[:3] + "**********"
	} else {
		mail = parts[0][:2] + "********@" + parts[1]
	}

	otp := util.GenerateSixDigits()
	currentTime := time.Now()
	otpExp := currentTime.Add(5 * time.Minute)
	smtpService := SmtpService{}
	newUser := !user.UserValid
	// request3rdPartyEmailAPI(user.Email, user.UserName, otp, newUser)
	go smtpService.SetNewPassWordMail(user.Email, user.UserName, otp, newUser)
	_ = qtx.SendNewOtp(ctx, auth.SendNewOtpParams{
		Otp:    otp,
		OtpExp: ToPGTimestamp(otpExp),
		Email:  user.Email,
	})
	return BuildResponse200("An OTP has been sent to your email "+mail+".", nil)
}

func (s *RESTService) userVerify(c *gin.Context) APIResponse {
	var i model.OtpVerify
	if !parseInput(c, &i) {
		return BuildResponse400("Invalid input provided")
	}
	if i.Login == "" || i.Otp == "" {
		return BuildResponse400("username/otp required")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	user, err := qtx.GetUserByLogin(ctx, i.Login)
	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	// check valid otp
	if user.Otp != i.Otp {
		return BuildResponse400("Invalid OTP. Please check your email")
	}
	// check time expire
	if user.OtpExp.Time.Before(time.Now()) {
		return BuildResponse400("OTP expired! Please click on resend to get new OTP")
	}
	err = qtx.ActivateUser(ctx, user.Email)
	if err != nil {
		return BuildResponse500("Unable to validate. Please try again.", err.Error())
	}
	return BuildResponse200("Congratulations!", nil)
}

func (s *RESTService) validateOtp(c *gin.Context) APIResponse {
	var i model.OtpVerify
	if !parseInput(c, &i) {
		return BuildResponse400("Invalid input provided")
	}
	if i.Login == "" || i.Otp == "" {
		return BuildResponse400("username/otp required")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	user, err := qtx.GetUserByLogin(ctx, i.Login)
	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	// check valid otp
	if user.Otp != i.Otp {
		return BuildResponse400("Invalid OTP. Please check your email")
	}
	// check time expire
	if user.OtpExp.Time.Before(time.Now()) {
		return BuildResponse400("OTP expired! Please click on resend to get new OTP")
	}
	return BuildResponse200("Verified!", nil)
}

// /api/auth/update - update user
func (s *RESTService) updateUser(c *gin.Context) APIResponse {
	var input model.UpdateUserInput
	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	// Validate required fields
	if input.UserID == 0 || input.Email == "" || input.Phone == "" || input.UserName == "" {
		return BuildResponse400("User ID, email, phone, and username are required")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	// Check if user exists
	currentUser, err := qtx.GetUserById(ctx, input.UserID)
	if err != nil {
		_asLogger.Errorf("Error getting user: %v", err)
		return BuildResponse404("User not found", false)
	}

	// Check email uniqueness (excluding current user)
	if input.Email != currentUser.Email {
		_, err = qtx.GetUserByEmail(ctx, input.Email)
		if err == nil {
			return BuildResponse400("User with this email already exists")
		}
	}

	// Check username uniqueness (excluding current user)
	if input.UserName != currentUser.UserName {
		_, err = qtx.GetUserByUserName(ctx, input.UserName)
		if err == nil {
			return BuildResponse400("User with this username already exists")
		}
	}

	// Check phone uniqueness (excluding current user)
	if input.Phone != currentUser.Phone {
		_, err = qtx.GetUserByPhone(ctx, input.Phone)
		if err == nil {
			return BuildResponse400("User with this phone number already exists")
		}
	}

	// Set default role if not provided
	role := input.Role
	if role == "" {
		role = currentUser.Role // Keep existing role if not provided
	}

	// Update user
	updateParams := auth.UpdateUserParams{
		UserName: input.UserName,
		Email:    input.Email,
		Phone:    input.Phone,
		Role:     role,
		UserID:   input.UserID,
	}

	err = qtx.UpdateUser(ctx, updateParams)
	if err != nil {
		_asLogger.Errorf("Error updating user: %v", err)
		return BuildResponse500("Failed to update user", err.Error())
	}

	return BuildResponse200("User updated successfully", nil)
}

// /api/auth/users - get all users (requires authentication)
func (s *RESTService) getAllUsers(c *gin.Context) APIResponse {
	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	users, err := qtx.GetAllUsers(ctx)
	if err != nil {
		_asLogger.Errorf("Error getting users: %v", err)
		return BuildResponse500("Failed to retrieve users", err.Error())
	}

	// Transform to response format
	userList := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		userList = append(userList, map[string]interface{}{
			"code":  user.UserID,
			"name":  user.UserName,
			"email": user.Email,
		})
	}

	return BuildResponse200("Users retrieved successfully", userList)
}
