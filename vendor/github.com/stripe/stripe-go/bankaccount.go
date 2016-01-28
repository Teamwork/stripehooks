package stripe

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// BankAccountStatus is the list of allowed values for the bank account's status.
// Allowed values are "new", "verified", "validated", "errored".
type BankAccountStatus string

// BankAccountParams is the set of parameters that can be used when creating or updating a bank account.
type BankAccountParams struct {
	Params

	// The identifier of the parent account under which bank accounts are
	// nested.
	AccountID string

	// A token referencing an external account like one returned from
	// Stripe.js.
	Token string

	// Information on an external account to reference. Only used if `Token`
	// is not provided.
	Account, Country, Currency, Routing string

	Default  bool
	Customer string
}

// BankAccountListParams is the set of parameters that can be used when listing bank accounts.
type BankAccountListParams struct {
	ListParams
	AccountID string
}

// BankAccount represents a Stripe bank account.
type BankAccount struct {
	ID          string            `json:"id"`
	Name        string            `json:"bank_name"`
	Country     string            `json:"country"`
	Currency    Currency          `json:"currency"`
	LastFour    string            `json:"last4"`
	Fingerprint string            `json:"fingerprint"`
	Status      BankAccountStatus `json:"status"`
	Routing     string            `json:"routing_number"`
	Deleted     bool              `json:"deleted"`
	Customer    *Customer         `json:"customer"`
}

// BankAccountList is a list object for bank accounts.
type BankAccountList struct {
	ListMeta
	Values []*BankAccount `json:"data"`
}

func (ba *BankAccount) Display() string {
	return fmt.Sprintf("Bank account ending in %s", ba.LastFour)
}

// AppendDetails adds the bank account's details to the query string values.
func (b *BankAccountParams) AppendDetails(values *url.Values) {
	values.Add("bank_account[country]", b.Country)
	values.Add("bank_account[routing_number]", b.Routing)
	values.Add("bank_account[account_number]", b.Account)

	if len(b.Currency) > 0 {
		values.Add("bank_account[currency]", b.Currency)
	}
}

// UnmarshalJSON handles deserialization of a BankAccount.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (b *BankAccount) UnmarshalJSON(data []byte) error {
	type bankAccount BankAccount
	var bb bankAccount
	err := json.Unmarshal(data, &bb)
	if err == nil {
		*b = BankAccount(bb)
	} else {
		// the id is surrounded by "\" characters, so strip them
		b.ID = string(data[1 : len(data)-1])
	}

	return nil
}
