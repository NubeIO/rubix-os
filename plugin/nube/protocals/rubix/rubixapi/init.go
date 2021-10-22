package rubixapi

type RestClient struct {
	ClientToken string
}

// New returns a new instance of FlowClient.
func New() *RestClient {
	return &RestClient{""}
}
