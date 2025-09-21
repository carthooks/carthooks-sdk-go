package carthooks

import "strconv"

type UrlSets struct {
	//3种尺寸， 原始尺寸， 128x128px,  26x26px
	FullSizeUrl string `json:"full_size_url"` //原始尺寸
	ThumbUrl    string `json:"thumb_url"`     //128x128px
	IconUrl     string `json:"icon_url"`      //26x26px
}

type ApiImageResult struct {
	Url      *UrlSets `json:"url"`
	Meta     any      `json:"meta"`
	Expired  int      `json:"expired"`
	FileSize int64    `json:"file_size"`
	Created  int      `json:"created"`
}

type RecordFormat struct {
	ID        uint                   `json:"id"`
	Title     string                 `json:"title"`
	CreatedAt int64                  `json:"created_at"`
	UpdatedAt int64                  `json:"updated_at"`
	Creator   uint                   `json:"creator"`
	Fields    map[string]interface{} `json:"fields"`
}

type EventMessage struct {
	Version string           `json:"version"`
	Meta    EventMessageMeta `json:"meta"`
	Payload any              `json:"payload"`
}

type EventCode string

const (
	EventCodeRecordCreated EventCode = "collection.item.created"
	EventCodeRecordUpdated EventCode = "collection.item.updated"
)

type EventMessageMeta struct {
	TenantID     uint      `json:"tenant_id"`
	CollectionID uint      `json:"collection_id"`
	Event        EventCode `json:"event"`
	TriggerType  string    `json:"trigger_type"`
	TriggerName  string    `json:"trigger_name,omitempty"`
}

func (e *EventMessageMeta) ToMap() map[string]string {
	return map[string]string{
		"tenant_id":     strconv.FormatUint(uint64(e.TenantID), 10),
		"collection_id": strconv.FormatUint(uint64(e.CollectionID), 10),
		"trigger_type":  e.TriggerType,
		"trigger_name":  e.TriggerName,
	}
}
