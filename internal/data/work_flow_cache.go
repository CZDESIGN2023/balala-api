package data

import (
	"context"
	"time"

	domain "go-cs/internal/domain/work_flow"

	goCache "github.com/Code-Hex/go-generics-cache"
)

func (r *workFlowRepo) GetWorkFlowTemplateFormMemoryCache(ctx context.Context, id int64) (*domain.WorkFlowTemplate, error) {

	fromCache, ok := r.flowTpltCacheAccessObj.Get(id)
	if ok {
		return fromCache, nil
	}

	formDb, err := r.GetFlowTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	r.flowTpltCacheAccessObj.Set(id, formDb, goCache.WithExpiration(time.Hour*48))

	return formDb, nil
}
