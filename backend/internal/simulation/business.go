package simulation

type BusinessType string

const (
	BusinessRetail     BusinessType = "RETAIL"
	BusinessService    BusinessType = "SERVICE"
	BusinessProduction BusinessType = "PRODUCTION"
)

type BusinessState struct {
	ID           string
	OwnerID      string
	Name         string
	Type         BusinessType
	Cash         int64
	Level        int
	Reputation   int
	Employees    int
	RevenueTotal int64
	ExpenseTotal int64
	Active       bool
}

func BusinessStartupCost(t BusinessType) int64 {
	switch t {
	case BusinessProduction:
		return 250
	case BusinessService:
		return 180
	default:
		return 150
	}
}

func BusinessDailyRevenue(b BusinessState, demand int64) int64 {
	if !b.Active {
		return 0
	}
	base := int64(15 + b.Level*10)
	if demand > 0 {
		base += demand * 2
	}
	if b.Reputation > 0 {
		base += int64(b.Reputation) / 50
	}
	return base
}

func BusinessWage(b BusinessState) int64 {
	return int64(5 + b.Level*2)
}
