package river

import (
	"go-cs/api/comm"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"go-cs/pkg/stream/tuple"
	"runtime/debug"
	"slices"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/spf13/cast"
)

type posSaver struct {
	pos   mysql.Position
	force bool
}

type MyEventHandler struct {
	canal.DummyEventHandler
	r *River
}

func (h *MyEventHandler) OnRotate(header *replication.EventHeader, e *replication.RotateEvent) error {
	pos := mysql.Position{
		Name: string(e.NextLogName),
		Pos:  uint32(e.Position),
	}

	h.r.syncCh <- posSaver{pos, true}

	return h.r.ctx.Err()
}

func (h *MyEventHandler) OnDDL(header *replication.EventHeader, nextPos mysql.Position, _ *replication.QueryEvent) error {
	h.r.syncCh <- posSaver{nextPos, true}
	return h.r.ctx.Err()
}

func (h *MyEventHandler) OnXID(header *replication.EventHeader, nextPos mysql.Position) error {
	h.r.syncCh <- posSaver{nextPos, false}
	return h.r.ctx.Err()
}

// 定义表名常量，避免硬编码
const (
	TableNameSpaceMember      = "space_member"
	TableNameSpaceWorkItemV2  = "space_work_item_v2"
	TableNameSpaceWorkObject  = "space_work_object"
	TableNameSpaceTag         = "space_tag"
	TableNameSpace            = "space"
	TableNameUser             = "user"
	TableNameSpaceWorkVersion = "space_work_version"
)

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	defer func() {
		if err := recover(); err != nil {
			h.r.log.Error("panic", err, string(debug.Stack()))
		}
	}()

	bus.Emit("canal.RowsEvent", e)

	if e.Action == canal.InsertAction {
		return nil
	}

	if !h.r.matchTable(e.Table.Schema, e.Table.Name) {
		return nil
	}

	// 提取主键
	if len(e.Table.PKColumns) != 1 {
		h.r.log.Error("主键数量不正确", e.Table.Name)
		return nil
	}

	var oldRows [][]any
	for i := 0; i < len(e.Rows); i += 2 {
		oldRows = append(oldRows, e.Rows[i])
	}

	pkIdx := e.Table.PKColumns[0]

	var event comm.CanalEvent
	var eventData any

	switch e.Table.Name {
	case TableNameSpaceMember:
		event = comm.CanalEvent_ce_SpaceMember
		spaceIdIdx := slices.IndexFunc(e.Table.Columns, func(column schema.TableColumn) bool {
			return column.Name == "space_id"
		})

		userIdIdx := slices.IndexFunc(e.Table.Columns, func(column schema.TableColumn) bool {
			return column.Name == "user_id"
		})

		if spaceIdIdx < 0 || userIdIdx < 0 {
			return nil
		}

		eventData = stream.Map(oldRows, func(v []any) tuple.Pair[int64, int64] {
			return tuple.T2(
				cast.ToInt64(v[spaceIdIdx]),
				cast.ToInt64(v[userIdIdx]),
			)
		})
	case TableNameSpaceWorkItemV2:
		event = comm.CanalEvent_ce_SpaceWorkItem
		eventData = stream.Map(oldRows, func(v []any) int64 {
			return cast.ToInt64(v[pkIdx])
		})
	case TableNameSpaceWorkObject:
		event = comm.CanalEvent_ce_SpaceWorkObject
		eventData = stream.Map(oldRows, func(v []any) int64 {
			return cast.ToInt64(v[pkIdx])
		})
	case TableNameSpaceTag:
		event = comm.CanalEvent_ce_SpaceTag
		eventData = stream.Map(oldRows, func(v []any) int64 {
			return cast.ToInt64(v[pkIdx])
		})
	case TableNameSpace:
		event = comm.CanalEvent_ce_Space
		eventData = stream.Map(oldRows, func(v []any) int64 {
			return cast.ToInt64(v[pkIdx])
		})
	case TableNameUser:
		event = comm.CanalEvent_ce_User
		eventData = stream.Map(oldRows, func(v []any) int64 {
			return cast.ToInt64(v[pkIdx])
		})
	case TableNameSpaceWorkVersion:
		event = comm.CanalEvent_ce_SpaceWorkVersion
		eventData = stream.Map(oldRows, func(v []any) int64 {
			return cast.ToInt64(v[pkIdx])
		})
	default:
		return nil
	}

	h.r.log.Debug(event, eventData)

	bus.Emit(event, eventData)

	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}
