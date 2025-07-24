package dto

type (
	PayloadLimiter struct {
		ClientID string `json:"client_id"`
		Route    string `json:"route" binding:"required"`
	}

	ResponseLimiter struct {
		Status        int `json:"status"`
		Limit         int `json:"limit"`
		Remain        int `json:"remain"`
		ResetInSecond int `json:"reset_in_second"`
	}

	ConfigClientRedis struct {
		ClientID string `json:"client_id"`
		Limit    int    `json:"limit"`
		Window   int    `json:"window"`
	}

	PayloadConfigClient struct {
		ClientID string `json:"client_id"`
		Route    string `json:"route" binding:"required"`
		Limit    int    `json:"limit"`
		Window   int    `json:"window"`
	}
)
