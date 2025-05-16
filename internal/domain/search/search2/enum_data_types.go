package search2

type DataType string

const (
	Integer     DataType = "integer"
	Text        DataType = "text"         //文本
	RichText    DataType = "rich-text"    //富文本
	User        DataType = "user"         //单选人员
	MultiUser   DataType = "multi-user"   //多选人员
	Select      DataType = "select"       //单选
	MultiSelect DataType = "multi-select" //多选
	Date        DataType = "date"         //日期
	DateRange   DataType = "date-range"   //日期区间
)

func (d DataType) IsBig() bool {
	switch d {
	case RichText:
		return true
	}

	return false
}
