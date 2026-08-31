package prompts

import _ "embed"

//go:embed filter.md
var Filter string

//go:embed planner.md
var Planner string

//go:embed selector.md
var Selector string

//go:embed plan_selector.md
var PlanSelector string
