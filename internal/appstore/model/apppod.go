package model

type AppPod struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// AppID 关联的应用ID
	AppID int64 `json:"app_id"`

	// AppPodID PodID
	AppPodID int64 `json:"app_pod_id"`
}
