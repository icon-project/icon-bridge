package algo

type bmcService string
type btsService int
type blacklistSvc int

const (
	FEE_GATHERING bmcService = "FeeGathering"
	LINK          bmcService = "Link"
	UNLINK        bmcService = "Unlink"
	INIT          bmcService = "Init"
)

const (
	REQUEST_COIN_TRANSFER btsService = iota
	REQUEST_COIN_REGISTER
	REPONSE_HANDLE_SERVICE
	BLACKLIST_MESSAGE
	CHANGE_TOKEN_LIMIT
	UNKNOWN_TYPE
)
const (
	ADD_TO_BLACKLIST blacklistSvc = iota
	REMOVE_FROM_BLACKLIST
)

type BMCMessage struct {
	Src     string //  an address of BMC (i.e. btp://1234.PARA/0x1234)
	Dst     string //  an address of destination BMC
	Svc     string //  service name of BSH
	Sn      []byte //  sequence number of BMC
	Message []byte //  serialized Service Message from BSH
}

type ServiceMessage struct {
	ServiceType []byte
	Payload     []byte
}

type ServiceMessagePayload struct { //Used for response events
	Code []byte
	Msg  []byte
}

// BMC Services
type FeeGatheringSvc struct {
	FeeAggregator string
	Services      []string
}

type LinkSvc struct {
	Addresses        []string //  Address of multiple Relays handle for this link network
	RxSeq            []byte
	TxSeq            []byte
	BlockIntervalSrc []byte
	BlockIntervalDst []byte
	MaxAggregation   []byte
	DelayLimit       []byte
	RelayIdx         []byte
	RotateHeight     []byte
	RxHeight         []byte
	RxHeightSrc      []byte
	IsConnected      bool
}

type InitSvc struct {
	Links []string
}

//BTS Services

type CoinTransferSvc struct {
	From   string
	To     string
	Assets []Asset
}

type Asset struct {
	CoinName string
	Value    []byte
}

type BlacklistSvc struct {
	RequestType []byte
	Addresses   []string
	Net         string
}

type TokenLimitSvc struct {
	CoinNames   []string
	TokenLimits [][]byte
	Net         string
}
