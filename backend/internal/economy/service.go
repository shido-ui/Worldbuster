package economy

import ("errors";"sync";"time";"crypto/rand";"encoding/hex")
var ErrInsufficientFunds=errors.New("insufficient funds")
var ErrInvalidAmount=errors.New("amount must be positive")
var ErrUnknownAccount=errors.New("unknown account")
var ErrIDEntropyFailure=errors.New("ledger id entropy failure")
var randomRead=rand.Read
type Service struct{mu sync.RWMutex;accounts map[string]*Account;ledger map[string][]LedgerEntry}
func NewService()*Service{return &Service{accounts:map[string]*Account{},ledger:map[string][]LedgerEntry{}}}
func(s *Service)EnsureAccount(id string)Account{s.mu.Lock();defer s.mu.Unlock();if a,ok:=s.accounts[id];ok{return *a};a:=&Account{ID:id,Currency:"WBX"};s.accounts[id]=a;return *a}
func(s *Service)Balance(id string)(Account,error){s.mu.RLock();defer s.mu.RUnlock();a,ok:=s.accounts[id];if !ok{return Account{},ErrUnknownAccount};return *a,nil}
func(s *Service)Credit(id string,amount int64,reason,ref string)(LedgerEntry,error){if amount<=0{return LedgerEntry{},ErrInvalidAmount};s.mu.Lock();defer s.mu.Unlock();a,ok:=s.accounts[id];if !ok{a=&Account{ID:id,Currency:"WBX"};s.accounts[id]=a};a.Balance+=amount;entry,err:=s.entry(a,amount,reason,ref);if err!=nil{return LedgerEntry{},err};return entry,nil}
func(s *Service)Debit(id string,amount int64,reason,ref string)(LedgerEntry,error){if amount<=0{return LedgerEntry{},ErrInvalidAmount};s.mu.Lock();defer s.mu.Unlock();a,ok:=s.accounts[id];if !ok{return LedgerEntry{},ErrUnknownAccount};if a.Balance<amount{return LedgerEntry{},ErrInsufficientFunds};a.Balance-=amount;entry,err:=s.entry(a,-amount,reason,ref);if err!=nil{return LedgerEntry{},err};return entry,nil}
func(s *Service)entry(a *Account,amount int64,reason,ref string)(LedgerEntry,error){b:=make([]byte,12);if _,err:=randomRead(b);err!=nil{return LedgerEntry{},ErrIDEntropyFailure};e:=LedgerEntry{ID:hex.EncodeToString(b),AccountID:a.ID,Amount:amount,BalanceAfter:a.Balance,Reason:reason,ReferenceID:ref,CreatedAt:time.Now().UTC()};s.ledger[a.ID]=append(s.ledger[a.ID],e);return e,nil}
func(s *Service)History(id string)[]LedgerEntry{s.mu.RLock();defer s.mu.RUnlock();return append([]LedgerEntry(nil),s.ledger[id]...)}
