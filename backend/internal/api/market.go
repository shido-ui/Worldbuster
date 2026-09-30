package api

import("encoding/json";"net/http";"github.com/shido-ui/Worldbuster/backend/internal/store")

type MarketAPI struct{Repo store.MarketRepository;Context *PlayerContext}
func(a *MarketAPI) Assets(w http.ResponseWriter,r *http.Request){x,err:=a.Repo.Assets(r.Context());if err!=nil{writeJSON(w,500,map[string]string{"error":"market unavailable"});return};writeJSON(w,200,map[string]any{"assets":x})}
func(a *MarketAPI) Orders(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("assetId");if id==""{writeJSON(w,400,map[string]string{"error":"assetId required"});return};x,err:=a.Repo.Orders(r.Context(),id);if err!=nil{writeJSON(w,500,map[string]string{"error":"orders unavailable"});return};if a.Context!=nil{if session,_,resolveErr:=a.Context.Resolve(r.Context(),r);resolveErr==nil{for i:=range x{x[i].IsOwner=x[i].SellerAccountID==session.AccountID}}};writeJSON(w,200,map[string]any{"orders":x})}
func(a *MarketAPI) Sell(w http.ResponseWriter,r *http.Request){s,p,err:=a.Context.Resolve(r.Context(),r);_ = s;if err!=nil{writeJSON(w,401,map[string]string{"error":"authentication required"});return};_ = p;var q struct{AssetID string `json:"assetId"`;Quantity int64 `json:"quantity"`;UnitPrice int64 `json:"unitPrice"`};if json.NewDecoder(r.Body).Decode(&q)!=nil||q.AssetID==""||q.Quantity<=0||q.UnitPrice<=0{writeJSON(w,400,map[string]string{"error":"assetId, positive quantity and positive unitPrice required"});return};var accountID string;account,err:=a.Context.Economy.GetByAccountID(r.Context(),s.AccountID);if err!=nil{writeJSON(w,500,map[string]string{"error":"economy unavailable"});return};accountID=account.ID;_ = accountID;x,err:=a.Repo.PlaceSell(r.Context(),s.AccountID,q.AssetID,q.Quantity,q.UnitPrice);if err!=nil{writeJSON(w,400,map[string]string{"error":"order rejected"});return};writeJSON(w,201,x)}
func(a *MarketAPI) Buy(w http.ResponseWriter,r *http.Request){s,_,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,401,map[string]string{"error":"authentication required"});return};var q struct{OrderID string `json:"orderId"`;Quantity int64 `json:"quantity"`};if json.NewDecoder(r.Body).Decode(&q)!=nil||q.OrderID==""||q.Quantity<=0{writeJSON(w,400,map[string]string{"error":"orderId and positive quantity required"});return};x,err:=a.Repo.Buy(r.Context(),s.AccountID,q.OrderID,q.Quantity);if err!=nil{writeJSON(w,400,map[string]string{"error":"purchase rejected"});return};writeJSON(w,200,x)}

func(a *MarketAPI) Cancel(w http.ResponseWriter,r *http.Request){
 s,_,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,401,map[string]string{"error":"authentication required"});return}
 var q struct{OrderID string `json:"orderId"`}
 if json.NewDecoder(r.Body).Decode(&q)!=nil||q.OrderID==""{writeJSON(w,400,map[string]string{"error":"orderId required"});return}
 x,err:=a.Repo.Cancel(r.Context(),s.AccountID,q.OrderID);if err!=nil{writeJSON(w,400,map[string]string{"error":"cancellation rejected"});return}
 writeJSON(w,200,x)
}
