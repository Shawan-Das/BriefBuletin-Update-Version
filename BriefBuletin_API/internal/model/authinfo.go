package model

import (
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
)

// const _AuthInfoTable = "hrm.authentication_info"
// const _ACLInfoTable = "hrm.acl_info"

var _ValidRoles = map[string]bool{
	"HR":          true,
	"DEPT_HEAD":   true,
	"EMP_MANAGER": true,
	"USER":        true,
	"SUPER_ADMIN": true,
}

// AuthenticationInfo reprents the row in hrm.authentication_info table
type AuthenticationInfo struct {
	Empcode        string      `json:"empcode"` //Mandatory field
	EmaiID         string      `json:"email,"`  //Key
	Password       string      `json:"password"`
	Role           interface{} `json:"role"`
	ExtraInfo      interface{} `json:"extrainfo,omitempty"`
	PasswordExpiry string      `json:"pwdexpiry"`
	PasswordStatus string      `json:"pwdstat"`
	Fname          string      `json:"empname,omitempty"`
	Lname          string      `json:"lastname,omitempty"`
}

// IsValidRole return true if the input role map is valid
func IsValidRole(roleMap map[string]interface{}) bool {
	if len(roleMap) == 0 {
		return false
	}
	for role := range roleMap {
		if _, isFound := _ValidRoles[role]; !isFound {
			return false
		}
	}
	return true
}

// ACLInfo contains entry of acl_info table
type ACLInfo struct {
	Action  string `json:"action"`
	EmpCode string `json:"empcode"`
	Role    string `json:"roles"`
}

// BuildACLInfo converts db results
func BuildACLInfo(objs []interface{}) []ACLInfo {
	jb, _ := json.Marshal(objs)
	authData := make([]ACLInfo, 0)
	err := json.Unmarshal(jb, &authData)
	if err != nil {
		return nil
	}
	return authData
}

type AuthDataInput struct {
	Email       string `json:"email,omitempty"` // For backward compatibility
	Login       string `json:"login,omitempty"` // Can be username, email, or phone
	Password    string `json:"pwd"`
	NewPassword string `json:"newPwd,omitempty"`
	Phone       string `json:"phone,omitempty"`
	UserName    string `json:"userName,omitempty"`
	Role        string `json:"role,omitempty"`
}

// AuthorizationClaims JWTTokenClaims
type AuthorizationClaims struct {
	UserID   int32  `json:"user_id"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type AuthServiceConfig struct {
	JWTKey     *string  `json:"jwtKey"`
	BypassAuth []string `json:"bypassAuth"`
}


type UserRoleInfo struct {
	UserID   int               `json:"user_id"`
	Email    string            `json:"email"`
	UserName string            `json:"username"`
	Role     string            `json:"role"`
	RoleMap  map[string]bool   `json:"role_map"`
}