package data

import (
	"context"
	"encoding/json"
	"errors"
	"go-cs/internal/data/convert"
	"go-cs/internal/utils"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"

	domain "go-cs/internal/domain/work_item"
	repo "go-cs/internal/domain/work_item/repo"
)

type spaceWorkItemEsRepo struct {
	baseRepo
	cache bool
	index string
}

func NewSpaceWorkItemEsRepo(data *Data, logger log.Logger) repo.WorkItemEsRepo {
	moduleName := "SpaceWorkItemRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &spaceWorkItemEsRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
		cache: true,
		index: data.conf.Es.Index,
	}

	return repo
}

func (p *spaceWorkItemEsRepo) CreateWorkItemEs(ctx context.Context, workItem *domain.WorkItem) error {
	index := p.index

	workItemEs := convert.WorkItemEntityToEs(workItem)
	esDocJSON, _ := json.Marshal(workItemEs)

	p.log.Debugf("es update doc: %s", string(esDocJSON))
	esRes, err := p.data.es.Index(index, strings.NewReader(string(esDocJSON)), p.data.es.Index.WithDocumentID(cast.ToString(workItem.Id)))
	if err != nil {
		return err
	}

	defer esRes.Body.Close()

	//{"error":{"root_cause":[{"type":"x_content_parse_exception","reason":"[1:2] [UpdateRequest] unknown field [id]"}],"type":"x_content_parse_exception","reason":"[1:2] [UpdateRequest] unknown field [id]"},"status":400}
	var r map[string]interface{}
	if err := json.NewDecoder(esRes.Body).Decode(&r); err != nil {
		p.log.Debugf("Error parsing the response body: %s", err)
		return err
	}

	if esRes.IsError() {
		p.log.Debugf("update request returned an error: %s", esRes.String())
		return errors.New(esRes.String())
	}

	if esRes.StatusCode == http.StatusCreated || esRes.StatusCode == http.StatusOK {
		return nil
	}

	p.log.Debug("update request failed: %s", esRes.String())
	return errors.New(esRes.String())
}
