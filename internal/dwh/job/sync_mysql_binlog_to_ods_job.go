package job

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/tidwall/sjson"
	"go-cs/internal/dwh/pkg"
	"go-cs/pkg/bus"
	"reflect"
	"runtime/debug"
	"slices"
)

// 从源库binlog中，导入到ods
type SyncMysqlBinlogToODSJob struct {
	ctx             *pkg.JobContext
	id              string
	binLogSourceCfg *mysql.Config
}

func NewSyncMysqlBinlogToODSJob(
	id string,
	ctx *pkg.JobContext,
) *SyncMysqlBinlogToODSJob {

	binLogSourceCfg, err := mysql.ParseDSN(ctx.Data.ExternalDataSourceConf().Database.Dsn)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &SyncMysqlBinlogToODSJob{
		id:              id,
		ctx:             ctx,
		binLogSourceCfg: binLogSourceCfg,
	}
}

func (job *SyncMysqlBinlogToODSJob) Name() string {
	return "SyncMysqlBinlogToODSJob"
}

func (job *SyncMysqlBinlogToODSJob) Id() string {
	return "SyncMysqlBinlogToODSJob:" + uuid.NewString()
}

func (job *SyncMysqlBinlogToODSJob) Run() {
	job.start()
}

func (job *SyncMysqlBinlogToODSJob) start() {

	listenerId := "dwh_to_ods_job"
	bus.On("canal.RowsEvent", listenerId, func(args ...any) {
		// 只能同步处理，不然顺序错误会导致统计结果错误.
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("handle event %s panic: %v, %s", "canal.RowsEvent", err, debug.Stack())
			}
		}()

		var rvs []reflect.Value
		for _, v := range args {
			rvs = append(rvs, reflect.ValueOf(v))
		}
		reflect.ValueOf(job.canalOnRowHandler).Call(rvs)
	})

}

func (job *SyncMysqlBinlogToODSJob) canalOnRowHandler(e *canal.RowsEvent) {

	if job.binLogSourceCfg.DBName != e.Table.Schema {
		return
	}

	//只记录增删改
	if !slices.Contains([]string{canal.UpdateAction, canal.DeleteAction, canal.InsertAction}, e.Action) {
		return
	}

	switch e.Table.Name {
	case "space":
		job.saveToOds(e, "ods_space_d")
	case "space_work_item_v2":
		job.saveToOds(e, "ods_witem_d")
	case "work_item_status":
		job.saveToOds(e, "ods_witem_status_d")
	case "space_work_object":
		job.saveToOds(e, "ods_object_d")
	case "space_work_version":
		job.saveToOds(e, "ods_version_d")
	case "user":
		job.saveToOds(e, "ods_user_d")
	case "space_member":
		job.saveToOds(e, "ods_member_d")
	case "space_work_item_flow_v2":
		job.saveToOds(e, "ods_witem_flow_node_d")
	}

}

func (job *SyncMysqlBinlogToODSJob) saveToOds(e *canal.RowsEvent, odsTableName string) {
	fmt.Println("***********")
	_, afterData := convertRowValueToMap(e)
	if odsTableName == "ods_witem_d" {
		for _, datum := range afterData {
			var err error
			doc, ok := datum["doc"].(string)

			fmt.Println("====================", doc)
			if !ok {
				continue
			}

			// 删除doc中的describe和remark字段
			doc, err = sjson.Set(doc, "describe", "")
			if err != nil {
				continue
			}

			doc, err = sjson.Set(doc, "remark", "")
			if err != nil {
				continue
			}

			datum["doc"] = doc
		}
	}

	for i := 0; i < len(afterData); i++ {
		afterData := afterData[i]
		if e.Action == canal.DeleteAction {
			if cast.ToInt64(afterData["deleted_at"]) == 0 {
				afterData["deleted_at"] = e.Header.Timestamp
			}
		}

		afterData["_op_ts"] = e.Header.Timestamp

		if odsTableName == "ods_witem_d" {
			if afterData["id"] == uint64(17764) {
				fmt.Println(afterData)
			}
		}

		err := job.ctx.Data.Db().Table(odsTableName).Create(afterData).Error
		if err != nil {
			fmt.Println(err)
		}
	}
}

func convertRowValueToMap(e *canal.RowsEvent) ([]map[string]interface{}, []map[string]interface{}) {

	beforeData := make([]map[string]interface{}, 0)
	afterData := make([]map[string]interface{}, 0)

	switch e.Action {
	case canal.InsertAction, canal.DeleteAction:
		for i := 0; i < len(e.Rows); i++ {
			afterRow := make(map[string]interface{})
			for j := 0; j < len(e.Table.Columns); j++ {
				afterRow[e.Table.Columns[j].Name] = e.Rows[i][j]
			}
			afterData = append(afterData, afterRow)
		}

	case canal.UpdateAction:
		for i := 0; i < len(e.Rows); i = i + 2 {
			beforeRow := make(map[string]interface{})
			afterRow := make(map[string]interface{})
			for j := 0; j < len(e.Table.Columns); j++ {
				beforeRow[e.Table.Columns[j].Name] = e.Rows[i][j]
				afterRow[e.Table.Columns[j].Name] = e.Rows[i+1][j]
			}
			beforeData = append(beforeData, beforeRow)
			afterData = append(afterData, afterRow)
		}
	}

	return beforeData, afterData
}
