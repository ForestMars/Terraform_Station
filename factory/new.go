package factory

import (
	"github.com/go-gorm/gorm"
	"github.com/ForestMars/TerraformStation"
	"github.com/ForestMars/TerraformStation/internal"
)

var _ TerraformStation.TerraformStationService = (*internal.TerraformStationImpl)(nil)

func New(db *gorm.DB, cfg *TerraformStation.Config) (service *internal.TerraformStationImpl) {
	var err error
	service, err = internal.New(db, cfg)
	if err != nil {
		panic(err)
	}
	return
}
