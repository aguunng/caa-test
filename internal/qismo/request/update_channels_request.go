package request

type UserChannel struct {
	ChannelID int    `json:"channel_id"`
	Source    string `json:"source"`
}

type AgentUpdatedRequest struct {
	Email     string        `json:"email"`
	Name      string        `json:"name"`
	Channels  []UserChannel `json:"channels"`
	UserRoles []string      `json:"role_ids"`
}
