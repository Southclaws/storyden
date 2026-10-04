package settings

type oAuthAutonomousRegistrationModeEnum string

const (
	oAuthAutonomousRegistrationModeDisabled  oAuthAutonomousRegistrationModeEnum = "disabled"
	oAuthAutonomousRegistrationModeProtected oAuthAutonomousRegistrationModeEnum = "protected"
	oAuthAutonomousRegistrationModeApproval  oAuthAutonomousRegistrationModeEnum = "approval"
	oAuthAutonomousRegistrationModeOpen      oAuthAutonomousRegistrationModeEnum = "open"
)
