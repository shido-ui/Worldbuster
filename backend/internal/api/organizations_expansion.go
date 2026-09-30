package api

import ("encoding/json";"net/http";)

func(a *PersistentOrganizationsAPI) Alliance(w http.ResponseWriter,r *http.Request){
 _,p,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,401,map[string]string{"error":"authentication required"});return}
 var q struct{OrganizationID string `json:"organizationId"`;TargetOrganizationID string `json:"targetOrganizationId"`;Status string `json:"status"`}
 if json.NewDecoder(r.Body).Decode(&q)!=nil||q.OrganizationID==""||q.TargetOrganizationID==""{writeJSON(w,400,map[string]string{"error":"organizationId and targetOrganizationId required"});return}
 if q.Status=="" {err=a.Repo.ProposeAlliance(r.Context(),q.OrganizationID,q.TargetOrganizationID,p.ID)} else {err=a.Repo.SetAllianceStatus(r.Context(),q.OrganizationID,q.TargetOrganizationID,p.ID,q.Status)}
 if err!=nil{writeJSON(w,400,map[string]string{"error":"alliance operation rejected"});return};writeJSON(w,200,map[string]any{"status":"ok"})
}

func(a *PersistentOrganizationsAPI) Activity(w http.ResponseWriter,r *http.Request){
 _,p,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,401,map[string]string{"error":"authentication required"});return}
 id:=r.URL.Query().Get("organizationId");if id==""{writeJSON(w,400,map[string]string{"error":"organizationId required"});return}
 ok,err:=a.Repo.IsMember(r.Context(),id,p.ID);if err!=nil{writeJSON(w,500,map[string]string{"error":"membership unavailable"});return};if !ok{writeJSON(w,403,map[string]string{"error":"organization membership required"});return}
 score,err:=a.Repo.UpdateActivity(r.Context(),id,1);if err!=nil{writeJSON(w,400,map[string]string{"error":"activity update rejected"});return}
 writeJSON(w,200,map[string]any{"organizationId":id,"activityScore":score})
}

