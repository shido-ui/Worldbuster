package simulation

type EconomySignal struct {
 Supply int64
 Demand int64
 BasePrice int64
 CurrentPrice int64
}

type EconomyAdjustment struct {
 NewPrice int64
 Direction int
 Pressure int64
}

func CalculatePriceAdjustment(s EconomySignal) EconomyAdjustment {
 if s.BasePrice <= 0 || s.CurrentPrice <= 0 { return EconomyAdjustment{NewPrice:1} }
 pressure:=s.Demand-s.Supply
 denominator:=s.Supply+s.Demand+10
 delta:=s.CurrentPrice*pressure/denominator
 maxStep:=s.CurrentPrice/10
 if maxStep<1 {maxStep=1}
 if delta>maxStep {delta=maxStep}
 if delta< -maxStep {delta=-maxStep}
 next:=s.CurrentPrice+delta
 if next<1 {next=1}
 if pressure>0 {return EconomyAdjustment{NewPrice:next,Direction:1,Pressure:pressure}}
 if pressure<0 {return EconomyAdjustment{NewPrice:next,Direction:-1,Pressure:pressure}}
 return EconomyAdjustment{NewPrice:next,Direction:0,Pressure:0}
}
