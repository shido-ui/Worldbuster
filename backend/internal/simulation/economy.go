package simulation

type EconomicState struct {
	Balance          int64
	LifetimeIncome   int64
	LifetimeSpending int64
}

func (e EconomicState) CanAfford(amount int64) bool { return amount > 0 && e.Balance >= amount }
func (e *EconomicState) Earn(amount int64) {
	if e == nil || amount <= 0 {
		return
	}
	e.Balance += amount
	e.LifetimeIncome += amount
}
func (e *EconomicState) Spend(amount int64) bool {
	if e == nil || amount <= 0 || e.Balance < amount {
		return false
	}
	e.Balance -= amount
	e.LifetimeSpending += amount
	return true
}
func economicCost(a ActionType) int64 {
	switch a {
	case ActionStudy:
		return 2
	case ActionSocialize:
		return 1
	case ActionTravel:
		return 5
	}
	return 0
}
