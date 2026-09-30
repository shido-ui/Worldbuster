package economy

import "time"

type Money struct { Amount int64 `json:"amount"`; Currency string `json:"currency"` }
type LedgerEntry struct { ID string `json:"id"`; AccountID string `json:"accountId"`; Amount int64 `json:"amount"`; BalanceAfter int64 `json:"balanceAfter"`; Reason string `json:"reason"`; ReferenceID string `json:"referenceId,omitempty"`; CreatedAt time.Time `json:"createdAt"` }
type Account struct { ID string `json:"id"`; Balance int64 `json:"balance"`; Currency string `json:"currency"` }
