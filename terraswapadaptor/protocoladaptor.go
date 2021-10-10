package terraswapadaptor

import (
	"errors"
	"strconv"

	"github.com/arken-lab/arken-adaptor/terraadaptor"
	"github.com/arken-lab/arken-adaptor/universaladaptor"
)

const (
	UUSD_DENOM       = "uusd"
	DECIMALS_DIVISOR = 1000000
)

type ProtocolAdaptorImpl struct {
	ContractAddresses   []string
	EventRequiredFields []string
}

type SwapEvent struct {
	ContractAddress  string
	Action           string
	Sender           string
	Receiver         string
	OfferAsset       string
	AskAsset         string
	OfferAmount      float64
	ReturnAmount     float64
	TaxAmount        float64
	SpreadAmount     float64
	CommissionAmount float64
	LogIndex         int64
}

func NewProtocolAdaptor(contractAddresses []string) universaladaptor.ProtocolAdaptor {
	c := &ProtocolAdaptorImpl{
		ContractAddresses: contractAddresses,
		EventRequiredFields: []string{
			"contract_address",
			"action",
			"sender",
			"receiver",
			"offer_asset",
			"ask_asset",
			"offer_amount",
			"return_amount",
			"tax_amount",
			"spread_amount",
			"commission_amount",
		},
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

func (p *ProtocolAdaptorImpl) ProcessTransaction(uniTx universaladaptor.Transaction) ([]universaladaptor.Event, error) {
	tx, ok := uniTx.(terraadaptor.Transaction)
	if !ok {
		return nil, errors.New("interface incorrect")
	}
	events := []SwapEvent{}
	for _, log := range tx.Logs {
		for evIdx, ev := range log.Events {
			if p.isSwapEvent(ev) {
				swapEvent := SwapEvent{}
				for _, att := range ev.Attributes {
					switch att.Key {
					case "offer_amount":
						swapEvent.OfferAmount, _ = strconv.ParseFloat(att.Value, 64)
					case "return_amount":
						swapEvent.ReturnAmount, _ = strconv.ParseFloat(att.Value, 64)
					case "tax_amount":
						swapEvent.TaxAmount, _ = strconv.ParseFloat(att.Value, 64)
					case "spread_amount":
						swapEvent.SpreadAmount, _ = strconv.ParseFloat(att.Value, 64)
					case "commission_amount":
						swapEvent.CommissionAmount, _ = strconv.ParseFloat(att.Value, 64)
					case "action":
						swapEvent.Action = att.Value
					case "contract_address":
						swapEvent.ContractAddress = att.Value
					case "sender":
						swapEvent.Sender = att.Value
					case "receiver":
						swapEvent.Receiver = att.Value
					case "offer_asset":
						swapEvent.OfferAsset = att.Value
					case "ask_asset":
						swapEvent.AskAsset = att.Value
					}
				}
				swapEvent.LogIndex = int64(evIdx)
				events = append(events, swapEvent)
			}
		}
	}

	universalEvents := []universaladaptor.Event{}

	for _, ev := range events {
		if ev.AskAsset != UUSD_DENOM && ev.OfferAsset != UUSD_DENOM {
			continue
		}
		if ev.AskAsset == UUSD_DENOM {
			amountBase := ev.OfferAmount / DECIMALS_DIVISOR
			amountUsd := ev.ReturnAmount / DECIMALS_DIVISOR
			if amountBase != 0 {
				priceBaseUsd := amountUsd / amountBase
				timestamp, err := strconv.ParseInt(tx.Timestamp, 10, 64)
				if err != nil {
					return nil, err
				}
				universalEvent := universaladaptor.NewTradeTxEvent(&universaladaptor.TradeTx{
					Token: universaladaptor.TradeTxToken{
						Chain:   terraadaptor.CHAIN_NAME,
						Address: ev.OfferAsset,
					},
					TxHash:       tx.TxHash,
					AmountBase:   amountBase,
					AmountUsd:    amountUsd,
					PriceBaseUsd: priceBaseUsd,
					LogIndex:     ev.LogIndex,
					Block:        tx.Height,
					LpAddress:    ev.ContractAddress,
					Side:         universaladaptor.TradeTxSell,
					Timestamp:    timestamp,
					IsRemoved:    false,
					IsHistorical: false,
				})
				universalEvents = append(universalEvents, universalEvent)
			}
		}
		if ev.OfferAsset == UUSD_DENOM {
			amountBase := ev.ReturnAmount / DECIMALS_DIVISOR
			amountUsd := ev.OfferAmount / DECIMALS_DIVISOR
			if amountBase != 0 {
				priceBaseUsd := amountUsd / amountBase
				timestamp, err := strconv.ParseInt(tx.Timestamp, 10, 64)
				if err != nil {
					return nil, err
				}
				universalEvent := universaladaptor.NewTradeTxEvent(&universaladaptor.TradeTx{
					Token: universaladaptor.TradeTxToken{
						Chain:   terraadaptor.CHAIN_NAME,
						Address: ev.AskAsset,
					},
					TxHash:       tx.TxHash,
					AmountBase:   amountBase,
					AmountUsd:    amountUsd,
					PriceBaseUsd: priceBaseUsd,
					LogIndex:     ev.LogIndex,
					Block:        tx.Height,
					LpAddress:    ev.ContractAddress,
					Side:         universaladaptor.TradeTxSell,
					Timestamp:    timestamp,
					IsRemoved:    false,
					IsHistorical: false,
				})
				universalEvents = append(universalEvents, universalEvent)
			}
		}
	}

	return universalEvents, nil
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

func (p *ProtocolAdaptorImpl) isSwapEvent(event terraadaptor.TxEvent) bool {
	if event.Type != "wasm" {
		return false
	}
	for _, field := range p.EventRequiredFields {
		found := false
		for _, attribute := range event.Attributes {
			if attribute.Key == field {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	hasActionSwap := false
	for _, attribute := range event.Attributes {
		if attribute.Key == "action" && attribute.Value == "swap" {
			hasActionSwap = true
			break
		}
	}
	return hasActionSwap
}
