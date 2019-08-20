package utils

const (
	REPORT_ID_CYDEX_DAY  = "001" //CYDEX 按天统计报表
	REPORT_ID_CYDEX_HOUR = "002" //CYDEX 按小时统计报表

	REDIS_DB_STATISTICS_USER_INFO_DT = 1
	REDIS_DB_STATISTICS_CYDEX_DAY    = 2
	REDIS_DB_STATISTICS_CYDEX_HOUR   = 3
	REDIS_DB_DATA_SERVICE            = 4
	REDIS_DB_TRANSFER_NODE           = 5
	REDIS_DB_NOTIFY_NODE             = 6
	REDIS_DB_TRANSFER_SERVICE        = 7
	REDIS_DB_USER_SYSTEM             = 8
	REDIS_DB_CYDEX_MANAGER           = 9
	REDIS_DB_ZONE                    = 10
	REDIS_DB_OPERATION_MAINTENANCE   = 11
	REDIS_DB_CALCULATE               = 12

	REDIS_NM_SERVER_B = 20
	REDIS_NM_CLIENT   = 21

	REDIS_DB_SL_CYDEX_MANAGER = 29

	PRODUCT_CODE_CYDEX     = "001"
	PRODUCT_CODE_WEB       = "002" //官网
	PRODUCT_CODE_CATON_NET = "003"

	//月份 1,01,Jan,January
	//日　 2,02,_2
	//时　 3,03,15,PM,pm,AM,am
	//分　 4,04
	//秒　 5,05
	//年　 06,2006
	//周几 Mon,Monday
	//时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
	//时区字母缩写 MST

	// 日期格式 YYYY-MM-DD
	DATE_FORMAT_SHORT = "2006-01-02"
	// 时间格式 YYYY-MM-DD HH:MM:SS
	DATE_FORMAT_LONG = "2006-01-02 15:04:05"
)

type LMData struct {
	Cpu  []*LMCpu `json:"cpu"`
	Fan  []*LMFan `json:"fan"`
	Data string   `json:"data"` //存放sensor生成的原始数据，json字符串
}
type LMCpu struct {
	Temperature int `json:"temerature"` // 温度（根据sensor的检测值算出来的温度）
}
type LMFan struct {
	Rpm int `json:"rpm"` // 风扇转速（根据sensor的检测值算出来的）
}
