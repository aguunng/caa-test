package response

import "time"

type RoomsResponse struct {
	Data struct {
		CustomerRooms []CustomerRoom `json:"customer_rooms"`
	} `json:"data"`
	Meta struct {
		CurrentTotal int    `json:"current_total"`
		CursorAfter  string `json:"cursor_after"`
		CursorBefore string `json:"cursor_before"`
	} `json:"meta"`
	Status int `json:"status"`
}

type CustomerRoom struct {
	ChannelID               int         `json:"channel_id"`
	ContactID               int         `json:"contact_id"`
	ID                      int         `json:"id"`
	IsHandledByBot          bool        `json:"is_handled_by_bot"`
	IsResolved              bool        `json:"is_resolved"`
	IsWaiting               bool        `json:"is_waiting"`
	LastCommentSender       string      `json:"last_comment_sender"`
	LastCommentSenderType   string      `json:"last_comment_sender_type"`
	LastCommentText         string      `json:"last_comment_text"`
	LastCommentTimestamp    time.Time   `json:"last_comment_timestamp"`
	LastCustomerCommentText string      `json:"last_customer_comment_text"`
	LastCustomerTimestamp   time.Time   `json:"last_customer_timestamp"`
	Name                    string      `json:"name"`
	RoomBadge               interface{} `json:"room_badge"`
	RoomID                  string      `json:"room_id"`
	RoomType                string      `json:"room_type"`
	Source                  string      `json:"source"`
	UserAvatarURL           string      `json:"user_avatar_url"`
	UserID                  string      `json:"user_id"`
}
