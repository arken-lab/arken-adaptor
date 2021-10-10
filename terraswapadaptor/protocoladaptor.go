package terraswapadaptor

import (
	"errors"

	"github.com/arken-lab/arken-adaptor/terraadaptor"
	"github.com/arken-lab/arken-adaptor/universaladaptor"
)

type ProtocolAdaptorImpl struct {
	ContractAddresses []string
}

func NewProtocolAdaptor(contractAddresses []string) universaladaptor.ProtocolAdaptor {
	c := &ProtocolAdaptorImpl{
		ContractAddresses: contractAddresses,
	}
	return c
}

func (p *ProtocolAdaptorImpl) GetTransactionFilters() ([]universaladaptor.TransactionFilter, error) {
	filters := []universaladaptor.TransactionFilter{}
	for _, contractAddr := range p.ContractAddresses {
		filters = append(filters, universaladaptor.TransactionFilter(terraadaptor.TransactionFilter{
			ContractAddress: contractAddr,
		}))
	}
	return filters, nil
}

func (p *ProtocolAdaptorImpl) ProcessTransaction(universaladaptor.Transaction) ([]universaladaptor.Event, error) {
	return nil, errors.New("not implemented")
}

func (p *ProtocolAdaptorImpl) GetLogFilters() ([]universaladaptor.LogFilter, error) {
	return nil, errors.New("not implemented")
}

func (p *ProtocolAdaptorImpl) GetPollingContracts() ([]universaladaptor.ContractResponse, error) {
	return nil, errors.New("not implemented")
}

func (p *ProtocolAdaptorImpl) ProcessLog(log universaladaptor.Log) ([]universaladaptor.Event, error) {
	return nil, errors.New("not implemented")
}

func (p *ProtocolAdaptorImpl) ProcessPollingContractResponse(universaladaptor.ContractResponse) ([]universaladaptor.Event, error) {
	return nil, errors.New("not implemented")
}
