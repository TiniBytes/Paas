package model

type MiddleVersion struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment"`

	// MiddleTypeID 关联的中间件类型ID
	MiddleTypeID int64 `json:"middle_type_id"`

	// MiddleDockerImage 镜像地址
	MiddleDockerImage string `json:"middle_docker_image"`

	// MiddleVersion 镜像版本
	MiddleVersion string `json:"middle_version"`
}
