package internal

import (
	"github.com/go-gorm/gorm"
	"github.com/ForestMars/TerraformStation"
)

type TerraformStationImpl struct {
	db             *gorm.DB
	cfg            *TerraformStation.Config
	optionalParam1 string
	optionalParam2 string
}

func New(db *gorm.DB, cfg *TerraformStation.Config) (r *TerraformStationImpl, err error) {
	r = &TerraformStationImpl{
		db:  db,
		cfg: cfg,
	}
	return
}

func (impl *TerraformStationImpl) OptionalParam1(param1 string) *TerraformStationImpl {
	impl.optionalParam1 = param1
	return impl
}

func (impl *TerraformStationImpl) OptionalParam2(param2 string) *TerraformStationImpl {
	impl.optionalParam2 = param2
	return impl
}

func (impl *TerraformStationImpl) TFCommand(input *TerraformStation.TFCommandInput) (r *TerraformStation.TFCommandResult, err error) {
	r = &TerraformStation.TFCommandResult{Result: "result"}
	return
}
