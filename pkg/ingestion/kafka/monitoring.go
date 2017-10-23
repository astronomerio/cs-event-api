package kafka

type Stats struct {
	Brokers          map[string]Broker `json:"brokers"`
	MetadataCacheCnt int               `json:"metadata_cache_cnt"`
	MsgCnt           int               `json:"msg_cnt"`
	MsgMax           int               `json:"msg_max"`
	MsgSize          int               `json:"msg_size"`
	MsgSizeMax       int64             `json:"msg_size_max"`
	Name             string            `json:"name"`
	Replyq           int               `json:"replyq"`
	SimpleCnt        int               `json:"simple_cnt"`
	Time             int               `json:"time"`
	Topics           struct {
	} `json:"topics"`
	Ts   int64  `json:"ts"`
	Type string `json:"type"`
}

type Broker struct {
	BufGrow    int `json:"buf_grow"`
	IntLatency struct {
		Avg int `json:"avg"`
		Cnt int `json:"cnt"`
		Max int `json:"max"`
		Min int `json:"min"`
		Sum int `json:"sum"`
	} `json:"int_latency"`
	Name         string `json:"name"`
	Nodeid       int    `json:"nodeid"`
	OutbufCnt    int    `json:"outbuf_cnt"`
	OutbufMsgCnt int    `json:"outbuf_msg_cnt"`
	ReqTimeouts  int    `json:"req_timeouts"`
	Rtt          struct {
		Avg int `json:"avg"`
		Cnt int `json:"cnt"`
		Max int `json:"max"`
		Min int `json:"min"`
		Sum int `json:"sum"`
	} `json:"rtt"`
	Rx           int    `json:"rx"`
	Rxbytes      int    `json:"rxbytes"`
	Rxcorriderrs int    `json:"rxcorriderrs"`
	Rxerrs       int    `json:"rxerrs"`
	Rxpartial    int    `json:"rxpartial"`
	State        string `json:"state"`
	Stateage     int    `json:"stateage"`
	Throttle     struct {
		Avg int `json:"avg"`
		Cnt int `json:"cnt"`
		Max int `json:"max"`
		Min int `json:"min"`
		Sum int `json:"sum"`
	} `json:"throttle"`
	Toppars struct {
	} `json:"toppars"`
	Tx             int `json:"tx"`
	Txbytes        int `json:"txbytes"`
	Txerrs         int `json:"txerrs"`
	Txretries      int `json:"txretries"`
	WaitrespCnt    int `json:"waitresp_cnt"`
	WaitrespMsgCnt int `json:"waitresp_msg_cnt"`
	Wakeups        int `json:"wakeups"`
	ZbufGrow       int `json:"zbuf_grow"`
}
