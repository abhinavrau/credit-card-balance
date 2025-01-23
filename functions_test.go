package creditcardservice

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestGetBalance(t *testing.T) {
	req, err := http.NewRequest("GET", "/balance/1234", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/balance/{accountNumber}", getBalance)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var card CreditCard
	err = json.Unmarshal(rr.Body.Bytes(), &card)
	if err != nil {
		t.Fatal(err)
	}

	if card.AccountNumber != "1234" {
		t.Errorf("handler returned unexpected account number: got %v want %v", card.AccountNumber, "1234")
	}
}

func TestGetRecentTransactions(t *testing.T) {
	req, err := http.NewRequest("GET", "/transactions/1111", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/transactions/{accountNumber}", getRecentTransactions)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var transactions []Transaction
	err = json.Unmarshal(rr.Body.Bytes(), &transactions)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) != 10 {
		t.Errorf("handler returned unexpected number of transactions: got %v want %v", len(transactions), 10)
	}

	// Check if transactions are sorted by date (most recent first)
	for i := 1; i < len(transactions); i++ {
		if transactions[i].Date.After(transactions[i-1].Date) {
			t.Errorf("transactions are not sorted by date in descending order")
			break
		}
	}
}

func TestGetAccountStatus(t *testing.T) {
	req, err := http.NewRequest("GET", "/status/0987654321", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/status/{accountNumber}", getAccountStatus)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var accountStatus struct {
		Status        string        `json:"status"`
		DeclineReason DeclineReason `json:"decline_reason,omitempty"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &accountStatus)
	if err != nil {
		t.Fatal(err)
	}

	if accountStatus.Status != "declined" {
		t.Errorf("handler returned unexpected status: got %v want %v", accountStatus.Status, "declined")
	}

	if accountStatus.DeclineReason != CreditLimitReached {
		t.Errorf("handler returned unexpected decline reason: got %v want %v", accountStatus.DeclineReason, CreditLimitReached)
	}
}

func TestGetBalanceNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/balance/9999999999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/balance/{accountNumber}", getBalance)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestMain(m *testing.M) {
	// Setup test data
	creditCards = map[string]CreditCard{
		"1234": {
			AccountNumber:   "1234",
			CreditLimit:     5000,
			Balance:         1500,
			LastPaymentDate: time.Now().AddDate(0, 0, -15),
			Status:          "active",
		},
		"0987654321": {
			AccountNumber:   "0987654321",
			CreditLimit:     10000,
			Balance:         10001,
			LastPaymentDate: time.Now().AddDate(0, 0, -5),
			Status:          "declined",
			DeclineReason:   CreditLimitReached,
		},
	}

	transactions = []Transaction{
		{ID: "1", AccountNumber: "1111", Amount: -150.75, Date: time.Now().AddDate(0, 0, -1), Description: "Hotel - Paris"},
		{ID: "2", AccountNumber: "1111", Amount: -89.50, Date: time.Now().AddDate(0, 0, -1), Description: "Restaurant - Paris"},
		{ID: "3", AccountNumber: "1111", Amount: -200.00, Date: time.Now().AddDate(0, 0, -2), Description: "Train Ticket - London to Paris"},
		{ID: "4", AccountNumber: "1111", Amount: -75.25, Date: time.Now().AddDate(0, 0, -2), Description: "Taxi - London"},
		{ID: "5", AccountNumber: "1111", Amount: -1500.00, Date: time.Now().AddDate(0, 0, -3), Description: "Flight - New York to London"},
		{ID: "6", AccountNumber: "1111", Amount: -120.30, Date: time.Now().AddDate(0, 0, -3), Description: "Duty Free Shop - JFK Airport"},
		{ID: "7", AccountNumber: "1111", Amount: -45.00, Date: time.Now().AddDate(0, 0, -4), Description: "Taxi - New York"},
		{ID: "8", AccountNumber: "1111", Amount: -85.75, Date: time.Now().AddDate(0, 0, -4), Description: "Restaurant - New York"},
		{ID: "9", AccountNumber: "1111", Amount: -250.00, Date: time.Now().AddDate(0, 0, -5), Description: "Hotel - New York"},
		{ID: "10", AccountNumber: "1111", Amount: -60.25, Date: time.Now().AddDate(0, 0, -5), Description: "Souvenir Shop - Times Square"},
		{ID: "11", AccountNumber: "1111", Amount: -100.00, Date: time.Now().AddDate(0, 0, -6), Description: "Extra transaction"},
	}

	// Run the tests
	m.Run()
}
