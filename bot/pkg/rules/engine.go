package rules

import (
	"context"
	"sort"
)

type Engine struct {
	rules []Rule
}

func NewEngine(rules []Rule) *Engine {
	// orden por prioridad
	sort.SliceStable(rules, func(i, j int) bool { return rules[i].Priority > rules[j].Priority })
	return &Engine{rules: rules}
}

func (e *Engine) Eval(ctx context.Context, env Envelope) (ActionResult, bool) {
	for _, r := range e.rules {
		match := true
		for _, p := range r.WhenAll {
			if !p(env) {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		res, err := r.Then(ctx, env)
		if err == nil && res.Handled {
			return res, true
		}
		if r.StopChain {
			return ActionResult{Handled: false}, false
		}
	}
	return ActionResult{}, false
}
