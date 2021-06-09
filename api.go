package TerraformStation

type TerraformStationService interface {
	TFCommand(input *TFCommandInput) (r *TFCommandResult, err error)
}
