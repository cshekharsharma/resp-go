package resp3

type RecordResponse struct {
	Value interface{}
	Code  uint32
}

type ScalarRecord struct {
	Value  interface{}
	Type   uint8
	LAT    int64
	Expiry int64
}
