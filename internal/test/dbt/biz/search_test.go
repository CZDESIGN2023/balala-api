package biz

import (
	"context"
	v1 "go-cs/api/search/v1"
	"go-cs/internal/test/dbt"
	"go-cs/pkg/pprint"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
)

func TestSearchWorkItem(t *testing.T) {
	req := &v1.SearchSpaceWorkItemGroupInfoRequestV2{
		//SpaceId: 87,
		//ConditionGroup: &v1.ConditionGroup{
		//	Conjunction: "AND",
		//	Conditions: []*v1.Condition{
		//		{Field: "last_status_at", Values: []string{"2024/03/11 00:00:00", "2024/03/17 00:00:00"}, Operator: "BETWEEN"},
		//	},
		//},
		//Sorts: []*v1.Sort{
		//	{Field: "priority", Order: "DESC"},
		//	{Field: "work_object_id", Order: "DESC"},
		//},
		Groups: []*v1.GroupBy{
			{Field: "work_object_id"},
			//{Field: "priority"},
			//{Field: "directors"},
		},
	}

	info, err := dbt.UC.SearchUsecase.SearchGroupInfo(context.Background(), 42, req)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(info)
}

func TestSearchWorkItemByIds(t *testing.T) {
	ids := []int64{472, 1134, 1138, 1139, 1454, 1584, 1586, 1587, 1588, 1589, 1590}

	info, err := dbt.UC.SearchUsecase.SearchMySpaceWorkItemsByIdV2(context.Background(), 42, ids)
	if err != nil {
		t.Error(err)
	}

	t.Skip()
	marshal, err := protojson.Marshal(info)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(string(marshal))
}

func Benchmark_SearchMySpaceWorkItemsByIdV2(b *testing.B) {
	ids := []int64{595, 596, 594, 633, 550, 551, 540, 534, 508, 448, 404, 533, 399, 434, 401, 400, 391, 513, 451, 394}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = dbt.UC.SearchUsecase.SearchMySpaceWorkItemsByIdV2(context.Background(), 42, ids)
	}
}

func Benchmark_SearchWorkItem(b *testing.B) {
	req := &v1.SearchSpaceWorkItemGroupInfoRequestV2{
		SpaceId: 161,
		ConditionGroup: &v1.ConditionGroup{
			Conjunction: "AND",
			//Conditions: []*v1.Condition{
			//	{Field: "user_id", Values: []string{"42"}, Operator: "EQ"},
			//	{Field: "directors", Values: []string{"21", "42"}, Operator: "IN"},
			//},
		},
		//Sorts:     nil,
		Groups: []*v1.GroupBy{
			{Field: "work_object_id"},
		},
	}

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = dbt.UC.SearchUsecase.SearchGroupInfo(context.Background(), 42, req)
	}
}

func Test_SearchWorkItemEs(t *testing.T) {
	req := &v1.SearchSpaceWorkItemGroupInfoRequestV2{
		SpaceId: 87,
		ConditionGroup: &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				//{Field: "node_plan_time", Values: []string{"2020/11/11 00:00:00", "2024/11/11 00:00:00"}, Operator: "BETWEEN"},
				{Field: "node_status", Values: []string{"3"}, Operator: "IN"},
			},
		},
		//Sorts: []*v1.Sort{
		//	{Field: "work_object_id", Order: "DESC"},
		//	{Field: "version_id", Order: "DESC"},
		//	{Field: "work_object_id", Order: "DESC"},
		//},
		Groups: []*v1.GroupBy{
			{Field: "work_object_id"},
			//{Field: "space_id"},
			{Field: "priority"},
		},
	}

	res, err := dbt.UC.SearchUsecase.SearchGroupInfo(context.Background(), 42, req)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(res)
}
