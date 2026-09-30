package simulation

type ProductionRecipe struct {
 ID string
 Name string
 BusinessType BusinessType
 MinLevel int
 InputItem string
 InputQuantity int
 OutputItem string
 OutputQuantity int
 LaborCost int64
}

type ProductionResult struct {
 BusinessID string
 RecipeID string
 Produced bool
 InputItem string
 InputQuantity int
 OutputItem string
 OutputQuantity int
 LaborCost int64
 Reason string
}

func CanProduce(b BusinessState, r ProductionRecipe, inputQuantity int64) bool {
 if !b.Active || b.Type != r.BusinessType || b.Level < r.MinLevel || inputQuantity < int64(r.InputQuantity) {
  return false
 }
 return b.Cash >= r.LaborCost
}
