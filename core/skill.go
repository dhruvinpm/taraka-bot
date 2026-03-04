package core

import "context"

type Skill interface {
	Name() string
	Execute(ctx context.Context, params map[string]interface{}) error
}

type HuntSkill struct {
	engine *Engine
}

func (s *HuntSkill) Name() string { return "hunt" }

func (s *HuntSkill) Execute(ctx context.Context, params map[string]interface{}) error {
	country, _ := params["country"].(string)
	niche, _ := params["niche"].(string)
	return s.engine.RunHunt(ctx, country, niche)
}

type AnalyzeSkill struct {
	engine *Engine
}

func (s *AnalyzeSkill) Name() string { return "analyze" }

func (s *AnalyzeSkill) Execute(ctx context.Context, params map[string]interface{}) error {
	return s.engine.RunAnalysis(ctx)
}

type DiplomatSkill struct {
	engine *Engine
}

func (s *DiplomatSkill) Name() string { return "diplomat" }

func (s *DiplomatSkill) Execute(ctx context.Context, params map[string]interface{}) error {
	return s.engine.RunOutreach(ctx)
}
