package consts

import "go-cs/api/comm"

var platformName = map[comm.ThirdPlatformCode]string{
	comm.ThirdPlatformCode_pf_IM:     "IMChat",
	comm.ThirdPlatformCode_pf_QL:     "轻聊",
	comm.ThirdPlatformCode_pf_Halala: "Halala",
}

func GetPlatformName(code comm.ThirdPlatformCode) string {
	return platformName[code]
}
