package model

type AccessInfo struct {
	ID              int64  `db:"id"`
	EndpointAddress string `db:"endpoint_address"`
	Role            string `db:"role"`
	CreatedAt       string `db:"created_at"`
	UpdatedAt       string `db:"updated_at"`
}
