package fcmservercli

import (
	"encoding/json"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *FcmServerClient) SendNotification(data map[string]interface{}) map[string]interface{} {
	body, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	req := inst.client.R().SetBody(body)
	resp, err := nresty.FormatRestyResponse(req.Post("/send"))
	if err != nil {
		return nil
	}
	var content map[string]interface{}
	_ = json.Unmarshal(resp.Body(), &content)
	return content
}
