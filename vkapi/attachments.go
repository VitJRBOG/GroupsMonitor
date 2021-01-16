package vkapi

type attachment struct {
	Type      string `json:"text"`
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	AccessKey string `json:"access_key"`
	URL       string `json:"url"`
}

func parseAttachmentsData(items []interface{}) []attachment {
	var attachments []attachment

	for _, m := range items {
		item := m.(map[string]interface{})

		var a attachment

		a.Type = item["type"].(string)
		if a.Type == "photo" || a.Type == "video" || a.Type == "audio" ||
			a.Type == "doc" || a.Type == "poll" || a.Type == "link" {
			if a.Type == "link" {
				a.URL = item["link"].(map[string]interface{})["url"].(string)
			} else {
				a.OwnerID = int(item[a.Type].(map[string]interface{})["owner_id"].(float64))
				a.ID = int(item[a.Type].(map[string]interface{})["id"].(float64))
				if accessKey, exist := item[a.Type].(map[string]interface{})["access_key"]; exist {
					a.AccessKey = accessKey.(string)
				}
			}
			attachments = append(attachments, a)
		}
	}

	return attachments
}
