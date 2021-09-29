package crud_server_http

type ReadOptionsHTTP struct {
	Path      string
	CGIParams string
	PageNum   uint64
	AllCnt    uint64
}
