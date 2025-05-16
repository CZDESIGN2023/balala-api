package service

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewDwhService,
	NewDwhAds,
)

type DwhService struct {
	Ads *DwhAds
}

func NewDwhService(ads *DwhAds) *DwhService {
	return &DwhService{
		Ads: ads,
	}
}
