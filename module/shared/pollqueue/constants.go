package pollqueue

type PollRetryType string

const (
	NORMAL_RETRY    PollRetryType = "normal-retry"
	IMMEDIATE_RETRY PollRetryType = "immediate-retry"
	DELAYED_RETRY   PollRetryType = "delayed-retry"
	NEVER_RETRY     PollRetryType = "never-retry"
)
