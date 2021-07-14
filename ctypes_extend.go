package snap7go

type TS7DataItemGo struct {
	Area     int32
	WordLen  int32
	Result   int32
	DBNumber int32
	Start    int32
	Amount   int32

	/*
		读取和写值的数据
		如果是读这个字段不需要初始化
	*/

	Pdata []byte
}

func (g *TS7DataItemGo) ToC() TS7DataItem {
	return TS7DataItem{
		Area:     g.Area,
		WordLen:  g.Area,
		Result:   g.Result,
		DBNumber: g.DBNumber,
		Start:    g.Start,
		Amount:   g.Amount,
		Pdata:    &g.Pdata[0],
	}
}
