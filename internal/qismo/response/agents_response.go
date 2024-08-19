package response

type AgentsResponse struct {
	Data struct {
		Agents struct {
			CurrentPage int `json:"current_page"`
			Data        []struct {
				ID                  int         `json:"id"`
				Name                string      `json:"name"`
				Email               string      `json:"email"`
				AuthenticationToken string      `json:"authentication_token"`
				CreatedAt           string      `json:"created_at"`
				UpdatedAt           string      `json:"updated_at"`
				SdkEmail            string      `json:"sdk_email"`
				SdkKey              string      `json:"sdk_key"`
				IsAvailable         bool        `json:"is_available"`
				Type                int         `json:"type"`
				AvatarURL           string      `json:"avatar_url"`
				AppID               int         `json:"app_id"`
				IsVerified          bool        `json:"is_verified"`
				NotificationsRoomID string      `json:"notifications_room_id"`
				BubbleColor         interface{} `json:"bubble_color"`
				QismoKey            string      `json:"qismo_key"`
				DirectLoginToken    interface{} `json:"direct_login_token"`
				LastLogin           string      `json:"last_login"`
				ForceOffline        bool        `json:"force_offline"`
				DeletedAt           interface{} `json:"deleted_at"`
				IsTocAgree          bool        `json:"is_toc_agree"`
				TotpToken           interface{} `json:"totp_token"`
				IsReqOtpReset       interface{} `json:"is_req_otp_reset"`
				LatestService       struct {
					ID                    int         `json:"id"`
					UserID                int         `json:"user_id"`
					RoomLogID             int         `json:"room_log_id"`
					AppID                 int         `json:"app_id"`
					RoomID                string      `json:"room_id"`
					Notes                 interface{} `json:"notes"`
					ResolvedAt            interface{} `json:"resolved_at"`
					IsResolved            bool        `json:"is_resolved"`
					CreatedAt             string      `json:"created_at"`
					UpdatedAt             string      `json:"updated_at"`
					FirstCommentID        interface{} `json:"first_comment_id"`
					LastCommentID         interface{} `json:"last_comment_id"`
					RetrievedAt           string      `json:"retrieved_at"`
					FirstCommentTimestamp interface{} `json:"first_comment_timestamp"`
					DeletedAt             interface{} `json:"deleted_at"`
					SessionCreatedAt      interface{} `json:"session_created_at"`
				} `json:"latest_service"`
				AssignedRules        []interface{} `json:"assigned_rules"`
				CurrentCustomerCount int           `json:"current_customer_count"`
				TotalResolved        int           `json:"total_resolved"`
				TotalCustomers       int           `json:"total_customers"`
				AssignedAgentRoles   []struct {
					ID            int    `json:"id"`
					Name          string `json:"name"`
					IsDefaultRole bool   `json:"is_default_role"`
				} `json:"assigned_agent_roles"`
				IsSupervisor bool `json:"is_supervisor"`
			} `json:"data"`
			FirstPageURL string      `json:"first_page_url"`
			From         int         `json:"from"`
			LastPage     int         `json:"last_page"`
			LastPageURL  string      `json:"last_page_url"`
			NextPageURL  interface{} `json:"next_page_url"`
			Path         string      `json:"path"`
			PerPage      int         `json:"per_page"`
			PrevPageURL  interface{} `json:"prev_page_url"`
			To           int         `json:"to"`
			Total        int         `json:"total"`
		} `json:"agents"`
		Total       int `json:"total"`
		CurrentPage int `json:"current_page"`
	} `json:"data"`
}
