package creditcardservice

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"
)

type DeclineReason string

const (
	CreditLimitReached    DeclineReason = "You met your credit limit"
	TravelUsage           DeclineReason = "Traveled to a new city where you never used your card before"
	LargePurchaseFlagged  DeclineReason = "Your large purchase was flagged"
	IncorrectPaymentInfo  DeclineReason = "You entered incorrect payment information"
	MissedPayments        DeclineReason = "You have missed payments"
	ExpiredOrDeactivated  DeclineReason = "You're using an expired or deactivated card"
	CardHold              DeclineReason = "Your card has a hold on it"
)

type CreditCard struct {
	AccountNumber   string        `json:"account_number"`
	CreditLimit     float64       `json:"credit_limit"`
	Balance         float64       `json:"balance"`
	LastPaymentDate time.Time     `json:"last_payment_date"`
	Status          string        `json:"status"`
	DeclineReason   DeclineReason `json:"decline_reason,omitempty"`
}

type Transaction struct {
	ID            string    `json:"id"`
	AccountNumber string    `json:"account_number"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
	Description   string    `json:"description"`
}

var creditCards = map[string]CreditCard{
	"1234": {
		AccountNumber:   "1234",
		CreditLimit:     5000,
		Balance:         1500,
		LastPaymentDate: time.Now().AddDate(0, 0, -15),
		Status:          "active",
	},
	"0987": {
		AccountNumber:   "0987",
		CreditLimit:     10000,
		Balance:         10001,
		LastPaymentDate: time.Now().AddDate(0, 0, -5),
		Status:          "declined",
		DeclineReason:   CreditLimitReached,
	},
	"1111": {
		AccountNumber:   "1111",
		CreditLimit:     8000,
		Balance:         7000,
		LastPaymentDate: time.Now().AddDate(0, 0, -20),
		Status:          "declined",
		DeclineReason:   TravelUsage,
	},
	"4444": {
		AccountNumber:   "4444",
		CreditLimit:     15000,
		Balance:         14000,
		LastPaymentDate: time.Now().AddDate(0, 0, -10),
		Status:          "declined",
		DeclineReason:   LargePurchaseFlagged,
	},
	"7777": {
		AccountNumber:   "7777",
		CreditLimit:     6000,
		Balance:         5500,
		LastPaymentDate: time.Now().AddDate(0, 0, -30),
		Status:          "declined",
		DeclineReason:   MissedPayments,
	},
	"0000": {
		AccountNumber:   "0000",
		CreditLimit:     7000,
		Balance:         6000,
		LastPaymentDate: time.Now().AddDate(0, -1, 0),
		Status:          "declined",
		DeclineReason:   ExpiredOrDeactivated,
	},
}

var transactions = []Transaction{
	{ID: "1", AccountNumber: "1234", Amount: -100, Date: time.Now().AddDate(0, 0, -1), Description: "Restaurant"},
	{ID: "2", AccountNumber: "1234", Amount: -50, Date: time.Now().AddDate(0, 0, -2), Description: "Gas Station"},
	{ID: "3", AccountNumber: "0987", Amount: -500, Date: time.Now().AddDate(0, 0, -1), Description: "Electronics"},
	{ID: "4", AccountNumber: "4444", Amount: -5000, Date: time.Now().AddDate(0, 0, -1), Description: "Luxury Purchase"},
	{ID: "5", AccountNumber: "1111", Amount: -150.75, Date: time.Now().AddDate(0, 0, -1), Description: "Hotel - Paris"},
	{ID: "6", AccountNumber: "1111", Amount: -89.50, Date: time.Now().AddDate(0, 0, -1), Description: "Restaurant - Paris"},
	{ID: "7", AccountNumber: "1111", Amount: -200.00, Date: time.Now().AddDate(0, 0, -2), Description: "Train Ticket - London to Paris"},
	{ID: "8", AccountNumber: "1111", Amount: -75.25, Date: time.Now().AddDate(0, 0, -2), Description: "Taxi - London"},
	{ID: "9", AccountNumber: "1111", Amount: -1500.00, Date: time.Now().AddDate(0, 0, -3), Description: "Flight - New York to London"},
	{ID: "10", AccountNumber: "1111", Amount: -120.30, Date: time.Now().AddDate(0, 0, -3), Description: "Duty Free Shop - JFK Airport"},
	{ID: "11", AccountNumber: "1111", Amount: -45.00, Date: time.Now().AddDate(0, 0, -4), Description: "Taxi - New York"},
	{ID: "12", AccountNumber: "1111", Amount: -85.75, Date: time.Now().AddDate(0, 0, -4), Description: "Restaurant - New York"},
	{ID: "13", AccountNumber: "1111", Amount: -250.00, Date: time.Now().AddDate(0, 0, -5), Description: "Hotel - New York"},
	{ID: "14", AccountNumber: "1111", Amount: -60.25, Date: time.Now().AddDate(0, 0, -5), Description: "Souvenir Shop - Times Square"},
}

// CreditCardService handles all credit card related HTTP requests
func CreditCardService(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/balance/"):
		getBalance(w, r)
	case strings.HasPrefix(r.URL.Path, "/transactions/"):
		getRecentTransactions(w, r)
	case strings.HasPrefix(r.URL.Path, "/status/"):
		getAccountStatus(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	accountNumber := strings.TrimPrefix(r.URL.Path, "/balance/")
	card, ok := creditCards[accountNumber]
	if !ok {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(card)
}

func getRecentTransactions(w http.ResponseWriter, r *http.Request) {
	accountNumber := strings.TrimPrefix(r.URL.Path, "/transactions/")
	var recentTransactions []Transaction
	for _, t := range transactions {
		if t.AccountNumber == accountNumber {
			recentTransactions = append(recentTransactions, t)
		}
	}
	sort.Slice(recentTransactions, func(i, j int) bool {
		return recentTransactions[i].Date.After(recentTransactions[j].Date)
	})
	if len(recentTransactions) > 10 {
		recentTransactions = recentTransactions[:10]
	}
	json.NewEncoder(w).Encode(recentTransactions)
}

func getAccountStatus(w http.ResponseWriter, r *http.Request) {
	accountNumber := strings.TrimPrefix(r.URL.Path, "/status/")
	card, ok := creditCards[accountNumber]
	if !ok {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	status := struct {
		Status        string        `json:"status"`
		DeclineReason DeclineReason `json:"decline_reason,omitempty"`
	}{
		Status:        card.Status,
		DeclineReason: card.DeclineReason,
	}
	json.NewEncoder(w).Encode(status)
}
