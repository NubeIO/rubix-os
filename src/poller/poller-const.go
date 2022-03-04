package poller

type PollPriority int

const (
	PRIORITY_HIGH   PollPriority = iota //Only Read Point Value.
	PRIORITY_NORMAL                     //Write the value on COV, don't Read.
	PRIORITY_LOW                        //Write the value on every poll (poll rate defined by setting).
)

type PollRate int

const (
	RATE_FAST   PollRate = iota //Only Read Point Value.
	RATE_NORMAL                 //Write the value on COV, don't Read.
	RATE_SLOW                   //Write the value on every poll (poll rate defined by setting).
)
