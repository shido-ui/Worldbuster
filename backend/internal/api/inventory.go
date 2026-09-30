package api

import (
 "encoding/json"
 "net/http"
)

func(r *Router) playerInventory(w http.ResponseWriter,req *http.Request){
 if r.playerContext==nil{writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return}
 _,p,err:=r.playerContext.Resolve(req.Context(),req)
 if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
 inv,err:=r.playerContext.Inventory.Get(req.Context(),p.ID)
 if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"inventory unavailable"});return}
 writeJSON(w,http.StatusOK,inv)
}

func(r *Router) addPlayerInventory(w http.ResponseWriter,req *http.Request){
 if r.playerContext==nil{writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return}
 _,p,err:=r.playerContext.Resolve(req.Context(),req)
 if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
 var body struct{ItemID string `json:"itemId"`; Quantity int `json:"quantity"`}
 if err:=json.NewDecoder(req.Body).Decode(&body);err!=nil||body.ItemID==""||body.Quantity<=0{writeJSON(w,http.StatusBadRequest,map[string]string{"error":"invalid request"});return}
 if err:=r.playerContext.Inventory.Add(req.Context(),p.ID,body.ItemID,body.Quantity);err!=nil{writeJSON(w,http.StatusBadRequest,map[string]string{"error":err.Error()});return}
 inv,err:=r.playerContext.Inventory.Get(req.Context(),p.ID);if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"inventory unavailable"});return}
 writeJSON(w,http.StatusOK,inv)
}
