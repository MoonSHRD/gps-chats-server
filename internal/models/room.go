package models

type Room struct {
	ID             int            `json:"id,omitempty" db:"id"`
	Latitude       float32        `json:"latitude,omitempty" db:"latitude"`
	Longitude      float32        `json:"longitude,omitempty" db:"longitude"`
	CreatedAt      NullableTime   `json:"created_at,omitempty" db:"created_at"`
	TTL            int            `json:"ttl,omitempty" db:"ttl"`
	Categories     []ChatCategory `json:"categories,omitempty" db:"-"`
	RoomID         string         `json:"roomId,omitempty" db:"room_id"`
	ParentGroupID  string         `json:"parentGroupId,omitempty" db:"parent_group_id"`
	EventStartDate NullableTime   `json:"eventStartDate,omitempty" db:"event_start_date"`
	Name           string         `json:"name" db:"name"`
}
