package settings

type AutoSwitch struct {
	TriggerEndpoint string
	AfterSec        int
}

type Setting struct {
	OutputDir     string
	RecordOnStart bool
	AutoStop      *AutoSwitch
	AutoStart     *AutoSwitch
}
