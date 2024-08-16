package response

import "time"

type AgentDetailResponse struct {
	Data struct {
		Agent struct {
			AvatarURL            string        `json:"avatar_url"`
			CreatedAt            string        `json:"created_at"`
			CurrentCustomerCount int           `json:"current_customer_count"`
			Email                string        `json:"email"`
			ForceOffline         bool          `json:"force_offline"`
			ID                   int           `json:"id"`
			IsAvailable          bool          `json:"is_available"`
			IsReqOtpReset        interface{}   `json:"is_req_otp_reset"`
			LastLogin            time.Time     `json:"last_login"`
			Name                 string        `json:"name"`
			SdkEmail             string        `json:"sdk_email"`
			SdkKey               string        `json:"sdk_key"`
			Type                 int           `json:"type"`
			TypeAsString         string        `json:"type_as_string"`
			UserChannels         []UserChannel `json:"user_channels"`
			UserRoles            []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"user_roles"`
		} `json:"agent"`
	} `json:"data"`
	Status int `json:"status"`
}

type UserChannel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
