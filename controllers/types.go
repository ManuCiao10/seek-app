package controllers

import "time"

type UserPostLogin struct {
	ID       string `bson:"_id"` // bson tag is used for MongoDB
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserPostSignup struct {
	ID        string    `bson:"_id"` // bson tag is used for MongoDB
	Fullname  string    `form:"fullname" json:"fullname" binding:"required"`
	Username  string    `form:"username" json:"username" binding:"required"`
	Email     string    `form:"email" json:"email" binding:"required"`
	Password  string    `form:"password" json:"password" binding:"required"`
	ExpiresAt time.Time `json:"expiresAt"` // expire time for the token
}

type UserGoogle struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

// All Discord scope constants that can be used.
const (
	ScopeIdentify                   = "identify"
	ScopeBot                        = "bot"
	ScopeEmail                      = "email"
	ScopeGuilds                     = "guilds"
	ScopeGuildsJoin                 = "guilds.join"
	ScopeConnections                = "connections"
	ScopeGroupDMJoin                = "gdm.join"
	ScopeMessagesRead               = "messages.read"
	ScopeRPC                        = "rpc"                    // Whitelist only
	ScopeRPCAPI                     = "rpc.api"                // Whitelist only
	ScopeRPCNotificationsRead       = "rpc.notifications.read" // Whitelist only
	ScopeWebhookIncoming            = "webhook.Incoming"
	ScopeApplicationsBuildsUpload   = "applications.builds.upload" // Whitelist only
	ScopeApplicationsBuildsRead     = "applications.builds.read"
	ScopeApplicationsStoreUpdate    = "applications.store.update"
	ScopeApplicationsEntitlements   = "applications.entitlements"
	ScopeRelationshipsRead          = "relationships.read" // Whitelist only
	ScopeActivitiesRead             = "activities.read"    // Whitelist only
	ScopeActivitiesWrite            = "activities.write"   // Whitelist only
	ScopeApplicationsCommands       = "applications.commands"
	ScopeApplicationsCommandsUpdate = "applications.commands.update"
)
