package bean

import (
	"fmt"
	"sort"
)

// GetBit 读取字节数组中某一位的bit
func GetBit(data []byte, pos int) bool {
	if pos < 0 || pos >= len(data)*8 {
		return false
	}
	byteIndex := pos / 8
	bitIndex := pos % 8
	return ((data[byteIndex] >> bitIndex) & 1) != 0
}

// SetBit 在指定的字节数组中将某一位设置为1, 如果超过则扩展数组
func SetBit(data *[]byte, pos int) error {
	if pos < 0 {
		return fmt.Errorf("pos out of range")
	}

	// 计算需要的字节长度
	requiredLength := pos/8 + 1
	if requiredLength > len(*data) {
		// 扩展切片的大小
		*data = append(*data, make([]byte, requiredLength-len(*data))...)
	}

	byteIndex := pos / 8
	bitIndex := pos % 8
	(*data)[byteIndex] |= 1 << bitIndex

	return nil
}

// UnsetBit 将字节数组的某一位设置为0
func UnsetBit(data []byte, pos int) error {
	if pos < 0 || pos >= len(data)*8 {
		return fmt.Errorf("pos out of range")
	}

	byteIndex := pos / 8
	bitIndex := pos % 8
	data[byteIndex] &= ^(1 << bitIndex)

	return nil
}

// SyncFieldsToUpdateDate 在changes类型中, 存的是map, 转到ud中去
func SyncFieldsToUpdateDate(fields map[int]*FieldValue, ud *UpdateData) {
	// Clear the current content of UD.Values
	ud.Values = make([]*FieldValue, 0)

	// Get the keys of the Fields map
	fieldKeys := make([]int, 0, len(fields))
	for k := range fields {
		fieldKeys = append(fieldKeys, k)
	}

	// Sort the keys
	sort.Ints(fieldKeys)

	// Add the FieldValues to UD.Values in the order of the sorted keys
	for _, k := range fieldKeys {
		ud.Values = append(ud.Values, fields[k])
		// Set the corresponding bit in the Masks
		SetBit(&ud.Masks, k)
	}
}

func NewUpdateDataResp(...interface{}) *Response {
	creates := []*AnyObj{}

	return &Response{
		Random:  0,
		Creates: creates,
		Updates: nil,
		Deletes: nil,
	}
}

////////////////////////////////////////
//测试方法

//func (pb *ChangeChatUser) SetId(newVal int64) *ChangeChatUser {
//	pb.Id = newVal
//	return pb
//}

//func updateMaskCase() {
//	chatUser := NewChangeChatUser(100)
//	chatUser.Id = 11
//
//	ud := chatUser.SetSignature("我的签名").SetLevel(11).SetSex(2).Finish()
//	resp := Response{
//		Random:  0,
//		Updates: make([]*UpdateData, 1),
//	}
//	resp.Updates = append(resp.Updates, ud)
//}
