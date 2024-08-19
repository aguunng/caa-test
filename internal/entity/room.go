package entity

// Room ...
type Room struct {
	ID        string `json:"id"`
	ChannelID int    `json:"channel_id"`
	Source    string `json:"source"`
}
