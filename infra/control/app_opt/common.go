package app_opt

import (
	"github.com/watermint/toolbox/domain/common/model/mo_string"
	"github.com/watermint/toolbox/infra/control/app_budget"
	"runtime"
)

const (
	OutputNone      = "none"
	OutputText      = "text"
	OutputMarkdown  = "markdown"
	OutputJson      = "json"
	BudgetLow       = "low"
	BudgetNormal    = "normal"
	BudgetUnlimited = "unlimited"
)

type CommonOpts struct {
	// Automatically open the artifact folder, after successful execution
	AutoOpen bool

	// Limit bandwidth to downloading/uploading contents
	BandwidthKb int

	// Set concurrency of worker execution
	Concurrency int

	// Enable debug mode
	Debug bool

	// Enable experiments
	Experiment string

	// Memory budget
	BudgetMemory mo_string.SelectString

	// Storage budget
	BudgetStorage mo_string.SelectString

	// Set output format
	Output mo_string.SelectString

	// Explicitly set proxy the hostname and the port number
	Proxy mo_string.OptionalString

	// Quiet mode
	Quiet bool

	// Do not store token in the file
	Secure bool

	// Specify workspace path
	Workspace mo_string.OptionalString
}

func (z *CommonOpts) Preset() {
	z.AutoOpen = false
	z.BandwidthKb = 0
	z.Concurrency = runtime.NumCPU()
	z.Debug = false
	z.Experiment = ""
	z.BudgetMemory.SetOptions([]string{BudgetLow, BudgetNormal}, BudgetNormal)
	z.BudgetStorage.SetOptions(app_budget.StorageBudgets, string(app_budget.DefaultBudget))
	z.Output.SetOptions([]string{OutputText, OutputMarkdown, OutputJson, OutputNone}, OutputText)
	z.Proxy = mo_string.NewOptional("")
	z.Quiet = false
	z.Secure = false
	z.Workspace = mo_string.NewOptional("")
}

func Default() CommonOpts {
	com := CommonOpts{}
	com.Workspace = mo_string.NewOptional("")
	com.Proxy = mo_string.NewOptional("")
	com.BudgetMemory = mo_string.NewSelect()
	com.BudgetStorage = mo_string.NewSelect()
	com.Output = mo_string.NewSelect()
	com.Preset()
	return com
}
