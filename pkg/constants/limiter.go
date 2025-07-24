package constants

type RateLimitStatus string

const (
	RateLimitStatusAllow RateLimitStatus = "Allow"
	RateLimitStatusDeny  RateLimitStatus = "Deny"
)
