package settings

type AutoSwitch struct {
	TriggerEndpoint string
	AfterSec        int
}

type Setting struct {
	// URL prefix to serve the web UI and APIs. Default: "/isumid"
	Prefix        string
	OutputDir     string
	RecordOnStart bool
	AutoStop      *AutoSwitch
	AutoStart     *AutoSwitch
}
