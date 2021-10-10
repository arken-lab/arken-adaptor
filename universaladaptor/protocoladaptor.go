package universaladaptor

type ProtocolAdaptor interface {
	GetLogFilters() ([]LogFilter, error)
	GetTransactionFilters() ([]TransactionFilter, error)
	GetPollingContracts() ([]ContractResponse, error)
	ProcessLog(log Log) ([]Event, error)
	ProcessPollingContractResponse(ContractResponse) ([]Event, error)
	ProcessTransaction(Transaction) ([]Event, error)
}
