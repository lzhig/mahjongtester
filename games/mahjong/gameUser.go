package mahjong

type roomT struct {
	ante     int32  // 房间底注
	id       int32  // 房间ID
	seat     int32  // 房间位置
	anteType int32  // 房间类型
	currency string // 货币类型
}

type playerDataT struct {
	uid       uint32
	name      string
	timestamp int64
	extra     string
	token     string

	uin uint32

	rooms map[int32]*roomT
}
