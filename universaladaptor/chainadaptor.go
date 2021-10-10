package universaladaptor

type BlockNumber int64
type LogFilter interface{}
type Log interface{}
type TransactionFilter interface{}
type Transaction interface{}
type ContractResponse interface{}

type ChainAdaptor interface {
	GetLatestBlock() (BlockNumber, error)
	GetLogs(fromBlock int64, toBlock int64, logFilters []LogFilter) ([]Log, error)
	GetTransactions(fromBlock int64, toBlock int64, transactionFilters []TransactionFilter) ([]Transaction, error)
	QueryContract(contractQueryData interface{}) (ContractResponse, error)
}
