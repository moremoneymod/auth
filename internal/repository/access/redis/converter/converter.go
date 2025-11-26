package converter

import (
	serv "github.com/moremoneymod/auth/internal/model"
	cache "github.com/moremoneymod/auth/internal/repository/access/redis/model"
)

func ToAccessInfoFromCache(info *cache.AccessInfo) *serv.AccessInfoCache {
	return &serv.AccessInfoCache{
		Role:            info.Role,
		EndpointAddress: info.EndpointAddress,
	}
}
