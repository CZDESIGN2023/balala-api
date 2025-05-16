package data

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/conf"
	"go-cs/internal/consts"
	search22 "go-cs/internal/domain/search/search2"
	"go-cs/internal/utils"
	esV8 "go-cs/internal/utils/es/v8"
	"go-cs/pkg/stream"
	"net/http"
	"slices"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"

	search_repo "go-cs/internal/domain/search/repo"
)

type searchRepo struct {
	baseRepo
	esIndex              string
	maxInnerResultWindow int
}

func NewSearchRepo(data *Data, conf *conf.Bootstrap, logger log.Logger) search_repo.SearchRepo {
	moduleName := "SearchRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &searchRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},

		esIndex: conf.Data.Es.Index,
	}
}

func (p *searchRepo) QueryWorkItem(ctx context.Context, spaceIds []int64, group *search22.ConditionGroup, selectFields string) ([]*search22.Model, error) {
	var list []*search22.Model

	query, args := search22.BuildCondition(group)

	tx := p.data.RoDB(ctx).
		Where("space_id IN ?", spaceIds)
	if group != nil {
		tx.Where(query, args...)
	}

	err := tx.Select(selectFields).Find(&list).Error
	if err != nil {
		p.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (p *searchRepo) QueryWorkFlowNode(ctx context.Context, spaceIds []int64, group *search22.ConditionGroup) ([]*db.SpaceWorkItemFlowV2, error) {
	var list []*db.SpaceWorkItemFlowV2

	query, args := search22.BuildCondition(group)

	tx := p.data.RoDB(ctx).
		Where("space_id IN ?", spaceIds)
	if group != nil {
		tx.Where(query, args...)
	}

	err := tx.Find(&list).Error
	if err != nil {
		p.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (p *searchRepo) GetWorkItemIdsByQueryWorkFlowNode(ctx context.Context, spaceIds []int64, group *search22.ConditionGroup) ([]int64, error) {
	var list []int64

	query, args := search22.BuildCondition(group)

	tx := p.data.RoDB(ctx).Model(&db.SpaceWorkItemFlowV2{}).
		Where("space_id IN ?", spaceIds)
	if group != nil {
		tx.Where(query, args...)
	}

	err := tx.Distinct("work_item_id").Pluck("work_item_id", &list).Error
	if err != nil {
		p.log.Error(err)
		return nil, err
	}

	return stream.Unique(list), nil
}

func (p *searchRepo) PendingWorkItem(ctx context.Context, userId int64, spaceIds []int64) ([]int64, error) {
	type item struct {
		Id           int64
		LastStatusAt int64
		CreatedAt    int64
	}
	var list []item

	userIdStr := cast.ToString(userId)

	err := p.data.RoDB(ctx).Table("space_work_item_v2 t").
		Select("t.id, t.last_status_at, t.created_at").
		Joins("INNER JOIN work_item_status s ON t.work_item_status_id = s.id").
		Where("t.space_id IN ?", spaceIds).
		Where("s.status_type != ?", consts.WorkItemStatusType_Archived).
		Where("? MEMBER OF(doc->'$.directors')", userIdStr).
		Find(&list).Error

	slices.SortFunc(list, func(a, b item) int {
		av := cast.ToInt64(a.LastStatusAt)
		bv := cast.ToInt64(b.LastStatusAt)

		if av == 0 {
			av = a.CreatedAt
		}
		if bv == 0 {
			bv = b.CreatedAt
		}

		v := cmp.Compare(bv, av)
		if v != 0 {
			return v
		}

		return cmp.Compare(b.Id, a.Id)
	})

	if err != nil {
		return nil, err
	}

	return stream.Map(list, func(v item) int64 {
		return v.Id
	}), nil
}

func (p *searchRepo) QueryWorkItemEs(ctx context.Context, searchSource *esV8.SearchSource) (*esV8.SearchResult, error) {
	err := p.initIndexConfig(ctx)
	if err != nil {
		return nil, err
	}

	es := p.data.es
	indexName := p.esIndex

	//pprint.Println(searchSource)

	marshal, err := json.Marshal(searchSource)
	if err != nil {
		fmt.Println(err)
	}

	esRes, err := es.Search(
		es.Search.WithIndex(indexName),
		es.Search.WithBody(bytes.NewReader(marshal)),
	)

	if err != nil {
		return nil, err
	}
	defer esRes.Body.Close()

	r := &esV8.SearchResult{}
	// 解析查询结果
	if err := json.NewDecoder(esRes.Body).Decode(&r); err != nil {
		p.log.Errorf("Error parsing the response body: %s", err)
		return nil, err
	}

	if esRes.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if esRes.IsError() {
		p.log.Errorf("Search request returned an error: %s", esRes.String())
		return nil, errors.New(esRes.String())
	}

	if esRes.StatusCode != http.StatusOK {
		p.log.Debug("Search request failed: %s", esRes.String())
		return nil, errors.New(esRes.String())
	}

	return r, nil
}

func (p *searchRepo) initIndexConfig(ctx context.Context) error {
	if p.maxInnerResultWindow > 0 { // 已经设置过
		return nil
	}

	es := p.data.es
	indexName := p.esIndex

	settings, err := es.Indices.PutSettings(
		strings.NewReader(utils.ToJSON(map[string]any{
			"index.max_inner_result_window": 20000,
			"index.refresh_interval":        "50ms",
			"index.number_of_replicas":      0,
		})),
		es.Indices.PutSettings.WithIndex(indexName),
	)
	if err != nil {
		p.log.Errorf("Error put settings: %s", err)
		return err
	}
	defer settings.Body.Close()

	if !settings.IsError() {
		p.maxInnerResultWindow = 10000
	}

	return nil
}

func (p *searchRepo) QueryWorkItemEsByPid(ctx context.Context, pid []int64) ([]*search22.Model, error) {
	es := p.data.es
	indexName := p.esIndex

	var searchSource = esV8.NewSearchSource()

	pidAny := stream.Map(pid, func(v int64) any {
		return v
	})

	searchSource.
		FetchSource(false).
		Fields("pid").
		Query(esV8.NewBoolQuery().Filter(esV8.NewTermsQuery("pid", pidAny...)))

	marshal, err := json.Marshal(searchSource)
	if err != nil {
		return nil, err
	}

	esRes, err := es.Search(
		es.Search.WithIndex(indexName),
		es.Search.WithBody(bytes.NewReader(marshal)),
	)

	if err != nil {
		return nil, err
	}
	defer esRes.Body.Close()

	r := &esV8.SearchResult{}
	// 解析查询结果
	if err := json.NewDecoder(esRes.Body).Decode(&r); err != nil {
		p.log.Errorf("Error parsing the response body: %s", err)
		return nil, err
	}

	if esRes.IsError() {
		p.log.Errorf("Search request returned an error: %s", esRes.String())
		return nil, errors.New(esRes.String())
	}

	if esRes.StatusCode != http.StatusOK {
		p.log.Debug("Search request failed: %s", esRes.String())
		return nil, errors.New(esRes.String())
	}

	list := stream.Map(r.Hits.Hits, func(v *esV8.SearchHit) *search22.Model {
		_pids, _ := v.Fields.Float64s("pid")
		var pid int64
		if len(_pids) > 0 {
			pid = int64(_pids[0])
		}

		return &search22.Model{
			Id:  cast.ToInt64(v.Id),
			Pid: pid,
		}
	})

	return list, nil
}

func (p *searchRepo) SearchByName(ctx context.Context, spaceId int64, keyword string) ([]*search22.Model, error) {
	var list []*search22.Model
	err := p.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ? and work_item_name LIKE ?", spaceId, "%"+keyword+"%").
		Select(search22.SelectByColumn("id", "work_item_name")).
		Order("id DESC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
