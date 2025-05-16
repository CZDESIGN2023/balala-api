package search_es

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	v1 "go-cs/api/search/v1"
	esV8 "go-cs/internal/utils/es/v8"
)

type HitItem struct {
	Id     string    `json:"_id"`
	Fields ExtendMap `json:"fields"` // fields取少量字段，性能比_source好 https://www.elastic.co/guide/en/elasticsearch/reference/7.17/search-fields.html
	//Source ExtendMap `json:"_source"`

	//Sort   []any          `json:"sort"`
}

type Group struct {
	Field     string     // 分组字段
	Key       string     // 分组字段值
	DocCount  int        // 数量
	SubGroups []*Group   // 子分组
	Hits      []*HitItem // 仅最后一层分组存在
}

// FlattenGroup 递归展开聚合结果
func FlattenGroup(aggGroups [][]*Group) [][]*Group {
	var f func(*Group) [][]*Group
	f = func(g *Group) [][]*Group {
		if g == nil {
			return nil
		}

		if len(g.SubGroups) == 0 {
			return [][]*Group{{g}}
		}

		var subs [][]*Group
		for _, sub := range g.SubGroups {
			subs = append(subs, f(sub)...)
		}

		var rets [][]*Group
		for _, sub := range subs {
			rets = append(rets, append([]*Group{g}, sub...))
		}

		return rets
	}

	var ret [][]*Group
	for _, aggGroup := range aggGroups {
		for _, g := range aggGroup {
			ret = append(ret, f(g)...)
		}
	}

	return ret
}

// ConvertAgg 转换聚合结果
func ConvertAgg(aggregations esV8.Aggregations) [][]*Group {
	var ret [][]*Group
	for aggName, raw := range aggregations {
		var data map[string]any
		_ = json.Unmarshal(raw, &data)

		buckets, ok := data["buckets"].([]any)
		if !ok {
			continue
		}

		groups := convertAggBuckets(aggName, buckets)
		ret = append(ret, groups)
	}

	return ret
}

// ConvertHits 转换非聚合结果
func ConvertHits(searchHits []*esV8.SearchHit) [][]*Group {
	var hints []*HitItem
	for _, v := range searchHits {
		hints = append(hints, &HitItem{
			Id:     v.Id,
			Fields: ExtendMap(v.Fields),
		})
	}

	g := &Group{
		Field:    "",
		Key:      "",
		DocCount: len(searchHits),
		Hits:     hints,
	}

	return [][]*Group{{g}}
}

// ConvertAggBuckets 转换聚合结果中的buckets
func convertAggBuckets(groupField string, buckets []any) []*Group {
	var ret []*Group
	for _, bucketData := range buckets {
		bucket, ok := bucketData.(map[string]any)
		if !ok {
			continue
		}

		var key string
		var docCount int
		var topHints []*HitItem
		var subGroups []*Group

		for k, v := range bucket {
			switch k {
			case "key":
				key = cast.ToString(v)
			case "doc_count":
				docCount = cast.ToInt(v)
			case "top_hits":
				_topHints, _ := v.(map[string]any)
				topHints = convertTopHits(_topHints)
			default:
				nextGroup, ok := v.(map[string]any)
				if !ok {
					continue
				}

				nextBuckets, ok := nextGroup["buckets"].([]any)
				if !ok {
					continue
				}

				// 处理下一层的数据
				subGroups = convertAggBuckets(k, nextBuckets)
			}
		}

		ret = append(ret, &Group{
			Field:     groupField,
			Key:       key,
			DocCount:  docCount,
			SubGroups: subGroups,
			Hits:      topHints,
		})
	}

	return ret
}

// ConvertTopHits 转换聚合结果中的top_hits
func convertTopHits(topHints map[string]any) []*HitItem {
	var ret []*HitItem

	tmp1, ok := topHints["hits"]
	if !ok {
		return nil
	}

	tmp2, ok := tmp1.(map[string]any)
	if !ok {
		return nil
	}

	hints, ok := tmp2["hits"].([]any)
	if !ok {
		return nil
	}

	for _, v := range hints {
		item, ok := v.(map[string]any)
		if !ok {
			continue
		}

		//source, ok := item["_source"].(map[string]any)
		//if !ok {
		//	continue
		//}

		fields, ok := item["fields"].(map[string]any)
		if !ok {
			continue
		}

		ret = append(ret, &HitItem{
			Id:     cast.ToString(item["_id"]),
			Fields: fields,
			//Source: source,
		})
	}

	return ret
}

func CheckGroupField(list []*v1.GroupBy) error {
	for _, v := range list {
		switch QueryField(v.Field) {
		case UserIdField, WorkObjectIdField, SpaceIdField, PriorityField, DirectorsField, WorkItemFlowIdField, WorkItemStatusIdField, VersionIdField:
		default:
			return fmt.Errorf("分组参数不支持 %v", v.Field)
		}
	}
	return nil
}
