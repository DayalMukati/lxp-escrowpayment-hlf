package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type EscrowContract struct {
	contractapi.Contract
}

type Escrow struct {
	EscrowID   string  `json:"escrowID"`
	PayerID    string  `json:"payerID"`
	PayeeID    string  `json:"payeeID"`
	Amount     float64 `json:"amount"`
	PayerApproved bool `json:"payerApproved"`
	PayeeConfirmed bool `json:"payeeConfirmed"`
	Released   bool    `json:"released"`
}

func (e *EscrowContract) CreateEscrow(ctx contractapi.TransactionContextInterface, escrowID string, payerID string, payeeID string, amount float64) error {
	escrow := Escrow{
		EscrowID:   escrowID,
		PayerID:    payerID,
		PayeeID:    payeeID,
		Amount:     amount,
		PayerApproved: false,
		PayeeConfirmed: false,
		Released:   false,
	}

	escrowAsBytes, err := json.Marshal(escrow)
	if err != nil {
		return fmt.Errorf("failed to marshal escrow: %s", err.Error())
	}

	return ctx.GetStub().PutState(escrowID, escrowAsBytes)
}

func (e *EscrowContract) ApproveEscrow(ctx contractapi.TransactionContextInterface, escrowID string, approverID string) error {
	escrowAsBytes, err := ctx.GetStub().GetState(escrowID)
	if err != nil {
		return fmt.Errorf("failed to get escrow: %s", err.Error())
	}

	if escrowAsBytes == nil {
		return fmt.Errorf("escrow %s does not exist", escrowID)
	}

	escrow := new(Escrow)
	err = json.Unmarshal(escrowAsBytes, escrow)
	if err != nil {
		return fmt.Errorf("failed to unmarshal escrow: %s", err.Error())
	}

	// Approve the escrow for the payer
	if approverID == escrow.PayerID {
		escrow.PayerApproved = true
	} else if approverID == escrow.PayeeID {
		escrow.PayeeConfirmed = true
	} else {
		return fmt.Errorf("approver is neither payer nor payee")
	}

	escrowAsBytes, err = json.Marshal(escrow)
	if err != nil {
		return fmt.Errorf("failed to marshal escrow: %s", err.Error())
	}

	return ctx.GetStub().PutState(escrowID, escrowAsBytes)
}

func (e *EscrowContract) ReleaseEscrow(ctx contractapi.TransactionContextInterface, escrowID string) error {
	escrowAsBytes, err := ctx.GetStub().GetState(escrowID)
	if err != nil {
		return fmt.Errorf("failed to get escrow: %s", err.Error())
	}

	if escrowAsBytes == nil {
		return fmt.Errorf("escrow %s does not exist", escrowID)
	}

	escrow := new(Escrow)
	err = json.Unmarshal(escrowAsBytes, escrow)
	if err != nil {
		return fmt.Errorf("failed to unmarshal escrow: %s", err.Error())
	}

	// Ensure both payer and payee have approved
	if !escrow.PayerApproved || !escrow.PayeeConfirmed {
		return fmt.Errorf("both parties have not approved the escrow")
	}

	escrow.Released = true

	escrowAsBytes, err = json.Marshal(escrow)
	if err != nil {
		return fmt.Errorf("failed to marshal escrow: %s", err.Error())
	}

	return ctx.GetStub().PutState(escrowID, escrowAsBytes)
}

func (e *EscrowContract) QueryEscrowStatus(ctx contractapi.TransactionContextInterface, escrowID string) (*Escrow, error) {
	escrowAsBytes, err := ctx.GetStub().GetState(escrowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow: %s", err.Error())
	}

	if escrowAsBytes == nil {
		return nil, fmt.Errorf("escrow %s does not exist", escrowID)
	}

	escrow := new(Escrow)
	err = json.Unmarshal(escrowAsBytes, escrow)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal escrow: %s", err.Error())
	}

	return escrow, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(EscrowContract))
	if err != nil {
		fmt.Printf("Error creating escrow contract: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting escrow contract: %s", err.Error())
	}
}
