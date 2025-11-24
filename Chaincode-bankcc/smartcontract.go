package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Transfer représente une transaction de virement
type Transfer struct {
	ID          string   `json:"id"`
	FromAccount string   `json:"fromAccount"`
	ToAccount   string   `json:"toAccount"`
	Amount      int      `json:"amount"`
	Currency    string   `json:"currency"`
	Timestamp   string   `json:"timestamp"`
	Status      string   `json:"status"`     // pending, approved, executed
	ApprovedBy  []string `json:"approvedBy"` // liste des IDs des approbateurs
}

// SmartContract définit la structure du contrat
type SmartContract struct {
	contractapi.Contract
}

// InitLedger initialise le ledger avec un exemple (optionnel)
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	transfers := []Transfer{
		{
			ID:          "TX001",
			FromAccount: "CUST123",
			ToAccount:   "CUST456",
			Amount:      100,
			Currency:    "EUR",
			Timestamp:   time.Now().Format(time.RFC3339),
			Status:      "pending",
			ApprovedBy:  []string{},
		},
	}

	for _, transfer := range transfers {
		transferJSON, err := json.Marshal(transfer)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(transfer.ID, transferJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}
	return nil
}

// InitTransfer crée une nouvelle demande de virement
func (s *SmartContract) InitTransfer(
	ctx contractapi.TransactionContextInterface,
	id string,
	fromAccount string,
	toAccount string,
	amount int,
	currency string
) error {
	exists, err := s.TransferExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the transfer %s already exists", id)
	}

	transfer := Transfer{
		ID:          id,
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Amount:      amount,
		Currency:    currency,
		Timestamp:   time.Now().Format(time.RFC3339),
		Status:      "pending",
		ApprovedBy:  []string{},
	}

	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, transferJSON)
}

// ApproveTransfer ajoute un approbateur à la transaction
func (s *SmartContract) ApproveTransfer(
	ctx contractapi.TransactionContextInterface,
	id string,
	approverID string
) error {
	transfer, err := s.ReadTransfer(ctx, id)
	if err != nil {
		return err
	}

	// Vérifie que le statut est "pending"
	if transfer.Status != "pending" {
		return fmt.Errorf("transfer %s is not pending", id)
	}

	// Ajoute l'approbateur
	transfer.ApprovedBy = append(transfer.ApprovedBy, approverID)

	transfer.Status = "approved" // On simplifie : 1 approbation suffit

	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, transferJSON)
}

// ExecuteTransfer marque le virement comme exécuté
func (s *SmartContract) ExecuteTransfer(
	ctx contractapi.TransactionContextInterface,
	id string
) error {
	transfer, err := s.ReadTransfer(ctx, id)
	if err != nil {
		return err
	}

	if transfer.Status != "approved" {
		return fmt.Errorf("transfer %s is not approved", id)
	}

	transfer.Status = "executed"
	transfer.Timestamp = time.Now().Format(time.RFC3339)

	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, transferJSON)
}

// ReadTransfer lit une transaction depuis le ledger
func (s *SmartContract) ReadTransfer(
	ctx contractapi.TransactionContextInterface,
	id string
) (*Transfer, error) {
	transferJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if transferJSON == nil {
		return nil, fmt.Errorf("the transfer %s does not exist", id)
	}

	var transfer Transfer
	err = json.Unmarshal(transferJSON, &transfer)
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// TransferExists retourne true si la transaction existe
func (s *SmartContract) TransferExists(
	ctx contractapi.TransactionContextInterface,
	id string
) (bool, error) {
	transferJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return transferJSON != nil, nil
}

// GetAllTransfers retourne toutes les transactions
func (s *SmartContract) GetAllTransfers(
	ctx contractapi.TransactionContextInterface
) ([]*Transfer, error) {
	resultIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	var transfers []*Transfer
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var transfer Transfer
		err = json.Unmarshal(queryResponse.Value, &transfer)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, &transfer)
	}

	return transfers, nil
}

// Main fonction

func main() {
    fmt.Println("🔥🔥🔥 DÉBUT DU MAIN: démarrage du chaincode") // <-- À AJOUTER
    chaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
        fmt.Printf("❌ ERREUR NEWCHAINCODE: %s", err.Error())
        return
    }
    fmt.Println("✅ Chaincode créé, appel à Start()...") // <-- À AJOUTER
    if err := chaincode.Start(); err != nil {
        fmt.Printf("❌ ERREUR START: %s", err.Error())
        return
    }
}

