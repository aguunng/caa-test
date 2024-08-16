package response

import "time"

type AgentsRoomResponse struct {
	Data struct {
		Agents []struct {
			AvatarURL            string      `json:"avatar_url"`
			CreatedAt            string      `json:"created_at"`
			CurrentCustomerCount int         `json:"current_customer_count"`
			Email                string      `json:"email"`
			ForceOffline         bool        `json:"force_offline"`
			ID                   int         `json:"id"`
			IsAvailable          bool        `json:"is_available"`
			IsReqOtpReset        interface{} `json:"is_req_otp_reset"`
			LastLogin            time.Time   `json:"last_login"`
			Name                 string      `json:"name"`
			SdkEmail             string      `json:"sdk_email"`
			SdkKey               string      `json:"sdk_key"`
			Type                 int         `json:"type"`
			TypeAsString         string      `json:"type_as_string"`
			UserChannels         []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"user_channels"`
			UserRoles []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"user_roles"`
		} `json:"agents"`
	} `json:"data"`
	Meta struct {
		After      interface{} `json:"after"`
		Before     interface{} `json:"before"`
		PerPage    int         `json:"per_page"`
		TotalCount interface{} `json:"total_count"`
	} `json:"meta"`
	Status int `json:"status"`
}
