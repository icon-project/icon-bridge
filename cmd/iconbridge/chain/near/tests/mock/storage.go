package mock

type Response struct {
	Reponse interface{}
	Error   error
}

type Storage struct {
	LatestChainStatus      Response
	BmcLinkStatusMap       map[string]Response
	AccessKeyMap           map[string]Response
	BlockByHashMap         map[string]Response
	ReceiptProofMap        map[string]Response
	BlockByHeightMap       map[int64]Response
	LightClientBlockMap    map[string]Response
	ContractStateChangeMap map[int64]Response
	TransactionHash        Response
	TransactionResultMap   map[string]Response
	AccountMap             map[string]Response
}
