# 🔌 Arken Universal Adaptor
![Logo](/images/adaptor-logo.png?raw=true)
> _Note: This repository is the work during “SCB 10X DeFi Launchpad Hackathon”. The design of Arken Universal Adaptor is not yet finalized and subject to change._

Arken Universal Adaptor is a standardized interface written in Golang. The main purpose is for blockchain protocol developers to integrate their protocol with Arken to enable real-time charts and Arken DEX aggregator integration with ease.

There are mainly 2 interfaces in Arken Universal Adaptor design, **Chain Adaptor** and **Protocol Adaptor**. Chain Adaptor is an abstraction layer to unify interactions between chains and protocols. Protocol Adaptor is a protocol-specific code to get information from the transaction. E.g. transaction’s trading price, liquidity pool reserves, current exchange rates, etc.

![Architecture](/images/architecture.png?raw=true)

## Interface Terminology
Throughout the documentation, you might find many terms that we use. Here is the compiled terms.

### Chain Adaptor
A chain adaptor interface is an abstraction layer that unifies interaction between a protocol and the blockchain that the protocol is based on. Each blockchain has it own instantiation of the chain interface. Some chains might share the same implementation of the interface, for example, an EVM-based blockchains.

### Protocol Adaptor
A protocol adaptor interface is the definition of methods that are needed for Arken Trading Platform. The platform will invoke those methods and the returned data will be fed into Arken Trading Platform.

### Transaction
A transaction is a unit of action in a block. Single block in blockchain comprises of many transactions. A transaction contain the data that the protocol adaptor will use to process required information.

### Log
A log refers to the set of data emitted by a successful transaction. A transaction can emit multiple logs, or zero log. The log can be useful for extracting information and feed into Arken Trading Platform.

### Contract Query Data
A contract query data refers to the data inquiry in the smart contract. In some implementation of the protocol where log processing or transaction processing is not possible, querying directly from the smart contract in a polling manner is required. A contract query data instance is basically consisted of contract address and query data.

### Event
An event has the data structure that is meaningful to Arken Trading Platform. The platform can use this type data to integrate with the platform’s feature such as real-time charts, best trading route calculation, etc.

## Chain Adaptor
If the chain that the protocol is based on is not in the list of supported chains, a Chain Adaptor is required first.

### Supported Chain Adaptors
* EVM-based chains
* Terra 🆕

### Implementing New Chain Adaptor
If the chain the protocol is based on is not in the supported chain list, a chain adaptor is needed to be implemented first. Not all of the methods are required to be implemented. It depends on which type event that the protocol is using.

### Chain Adaptor Interface References

```golang
type ChainAdaptor interface {
	GetLatestBlock() (BlockNumber, error)
	GetLogs(fromBlock int64, toBlock int64, logFilters []LogFilter) ([]Log, error)
	GetTransactions(fromBlock int64, toBlock int64, transactionFilters []TransactionFilter) ([]Transaction, error)
	QueryContract(contractQueryData interface{}) (ContractResponse, error)
}
```

## Protocol Adaptor
In order to support new protocol, a protocol adaptor is required. There are different approaches to implement a protocol adaptor depending on how the protocol is implemented. We recommend you to look at the event types required by Arken Trading Platform and figure out how to get those information based on your protocol implementation.

### Log-based Events
If the information in the events can be retrieved from logs, we can use log-based implementation of the protocol adaptor. The `GetLogFilters` and `ProcessLog` are required to be implemented. Also, `GetLatestBlock` and `GetLogs` are required in the Chain Adaptor implementation. The example of this type of event is Swap log in Uniswap V2 protocol.

### Transaction-based Events
If the information in the events cannot be retrieved from logs, but can be retrieved from transactions, we can use transaction-based implementation of the protocol adaptor. The `GetTransactionFilters` and `ProcessTransaction` are required to be implemented. Also, `GetLatestBlock` and `GetTransactions` are required in the Chain Adaptor implementation.

### Contract-based Events
If extracting information both from log and transaction is not possible, we have another method of extracting data, which is smart contract query in a polling manner. This type of data retrieval is not recommended and should be saved as last resort due to the drop in performance. The example of this type of event is getting token reserve values in Dodoswap liquidity pool. The `GetPollingContracts` and `ProcessPollingContractResponse` are required to be implemented. Also, `GetLatestBlock` and `QueryContract` are required in the Chain Adaptor implementation.

### Supported Protocol Adaptors
Here are the list of supported protocol adaptors. If your protocol is already on the list, it means that Arken Trading Platform can instantly integrate your platform.
* Uniswap V2
* Dodoswap V1
* Dodoswap V2
* Ellipsis Finance
* TerraSwap 🆕

### Protocol Adaptor Interface References

```golang
type ProtocolAdaptor interface {
	GetLogFilters() ([]LogFilter, error)
	GetTransactionFilters() ([]TransactionFilter, error)
	GetPollingContracts() ([]ContractResponse, error)
	ProcessLog(log Log) ([]Event, error)
	ProcessPollingContractResponse(ContractResponse) ([]Event, error)
	ProcessTransaction(Transaction) ([]Event, error)
}
```
### Feature Integration Matrix
There are several features on Arken Trading Platform. Protocol Adaptor is not required to support all features. Here is the list of events that is required to be implemented in each feature.

| Adaptor Implementation                        | Real-time Chart | Best Route Calculation |
|-----------------------------------------------|:---------------:|:----------------------:|
| Chain Adaptor                                 |        ✅        |            ✅           |
| [Event] TradeTx                               |        ✅        |                        |
| [Event] PoolCreated, PoolUpdated, PoolDeleted |                 |            ✅           |
| Pool Interface Implementation                 |                 |            ✅           |

## Additional Interface References
```golang
type Event struct {
	Type    EventType
	payload interface{}
}

func (e *Event) GetTradeTx() (*TradeTx, error) // For type TradeTx

type TradeTx struct {
	Token        TradeTxToken
	TxHash       string
	LogIndex     int64
	Block        int64
	LpAddress    string
	AmountBase   float64
	AmountUsd    float64
	PriceBaseUsd float64
	Side         TradeTxSide
	Timestamp    int64
	IsRemoved    bool
	IsHistorical bool
}

func (e *Event) GetPool() (*Pool, error) // For type = PoolCreated, PoolUpdated, PoolDeleted

type Pool interface {
	GetChain() string
	GetAddress() string
	GetTokens() []Token
	GetExchangeRate(token Token, amount uint256.Int) uint256.Int
	GetMidPrice(quoteToken Token) (float64, error)
	GetReserve(token Token) (uint256.Int, error)
}
// Struct from Pool interface must be marshallable in JSON format

```
