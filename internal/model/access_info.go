package model

type AccessInfo struct {
	ID              int64
	EndpointAddress string
	Role            string
	CreatedAt       string
	UpdatedAt       string
}

type AccessInfoCache struct {
	EndpointAddress string
	Role            string
}
