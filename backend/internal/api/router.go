package api

import("context";"encoding/json";"net/http";"time";"github.com/shido-ui/Worldbuster/backend/internal/inventory";"github.com/shido-ui/Worldbuster/backend/internal/economy";"github.com/shido-ui/Worldbuster/backend/internal/organization";"github.com/shido-ui/Worldbuster/backend/internal/job";"github.com/shido-ui/Worldbuster/backend/internal/progression";"github.com/shido-ui/Worldbuster/backend/internal/social";"github.com/shido-ui/Worldbuster/backend/internal/events";"github.com/shido-ui/Worldbuster/backend/internal/world")

type PlayerContextResolver interface { PlayerID(context.Context,string)(string,error) }

type Router struct{progressionDB *ProgressionAPI; playerResolver PlayerContextResolver; playerContext *PlayerContext;persistentJobs *PersistentJobsAPI;education *EducationAPI;socialPersistent *PersistentSocialAPI;reputation *ReputationAPI;organizationsPersistent *PersistentOrganizationsAPI;territory *TerritoryAPI;market *MarketAPI;missions *MissionsAPI;combat *CombatAPI;equipment *EquipmentAPI;worldEvents *WorldEventsAPI;news *NewsAPI;achievements *AchievementsAPI;admin *AdminAPI;world *world.Service;auth AuthBackend;travel *world.TravelService;inventory *inventory.Service;economy *economy.Service;organizations *organization.Service;jobs *job.Service;progression *progression.Service;social *social.Service;events *events.Service}
func NewRouter(w *world.Service,a AuthBackend,t *world.TravelService,i *inventory.Service,e *economy.Service,o *organization.Service,j *job.Service,p *progression.Service,s *social.Service,ev *events.Service)*Router{return &Router{world:w,auth:a,travel:t,inventory:i,economy:e,organizations:o,jobs:j,progression:p,social:s,events:ev}}
func(r *Router)Handler()http.Handler{
 mux:=http.NewServeMux();mux.HandleFunc("/health",r.health);mux.HandleFunc("/api/v1/world",r.worldState);if r.playerContext!=nil { mux.HandleFunc("/api/v1/player/dashboard",r.playerDashboard) }
 wh:=newWorldHandler(r.travel,r.playerContext);ih:=newInventoryHandler(r.inventory);eh:=newEconomyHandler(r.economy);jh:=newJobHandler(r.jobs);ph:=newProgressionHandler(r.progression);sh:=newSocialHandler(r.social);ehv:=newEventHandler(r.events,r.playerContext);mux.HandleFunc("/api/v1/world/locations",wh.locations);mux.HandleFunc("/api/v1/world/travel",wh.travel);mux.HandleFunc("/api/v1/world/travel/current",wh.currentTravel);mux.HandleFunc("/api/v1/player/state",r.playerState);mux.HandleFunc("/api/v1/player/inventory",r.playerInventory);mux.HandleFunc("/api/v1/player/inventory/add",r.addPlayerInventory);mux.HandleFunc("/api/v1/player/economy",r.playerEconomy);if r.playerContext==nil { mux.HandleFunc("/api/v1/inventory",ih.get); mux.HandleFunc("/api/v1/inventory/add",ih.add); mux.HandleFunc("/api/v1/economy/balance",eh.balance); mux.HandleFunc("/api/v1/economy/credit",eh.credit); mux.HandleFunc("/api/v1/economy/history",eh.history) } else { mux.HandleFunc("/api/v1/inventory",r.playerInventory); mux.HandleFunc("/api/v1/inventory/add",func(w http.ResponseWriter,_ *http.Request){writeJSON(w,http.StatusForbidden,map[string]string{"error":"inventory changes are server-authoritative"})}); mux.HandleFunc("/api/v1/economy/balance",r.playerEconomy); mux.HandleFunc("/api/v1/economy/history",r.playerEconomy); mux.HandleFunc("/api/v1/economy/credit",func(w http.ResponseWriter,_ *http.Request){writeJSON(w,http.StatusForbidden,map[string]string{"error":"economy credits are server-issued"})}) };if r.education!=nil { mux.HandleFunc("/api/v1/education/courses",r.education.Courses); mux.HandleFunc("/api/v1/education/status",r.education.Status); mux.HandleFunc("/api/v1/education/enroll",r.education.Enroll) }; if r.persistentJobs!=nil { mux.HandleFunc("/api/v1/jobs",r.persistentJobs.List); mux.HandleFunc("/api/v1/jobs/status",r.persistentJobs.Status); mux.HandleFunc("/api/v1/jobs/employ",r.persistentJobs.Employ) } else { mux.HandleFunc("/api/v1/jobs",jh.list); mux.HandleFunc("/api/v1/jobs/employ",jh.employ) };mux.HandleFunc("/api/v1/progression",ph.get);mux.HandleFunc("/api/v1/progression/courses",ph.courses);mux.HandleFunc("/api/v1/progression/enroll",ph.enroll);if r.socialPersistent!=nil { mux.HandleFunc("/api/v1/social/send",r.socialPersistent.Send); mux.HandleFunc("/api/v1/social/inbox",r.socialPersistent.Inbox); mux.HandleFunc("/api/v1/social/relationship",r.socialPersistent.Relationship) } else { mux.HandleFunc("/api/v1/social/send",sh.send); mux.HandleFunc("/api/v1/social/inbox",sh.inbox) };if r.reputation!=nil { mux.HandleFunc("/api/v1/reputation",r.reputation.Get) };if r.market!=nil { mux.HandleFunc("/api/v1/market/assets",r.market.Assets); mux.HandleFunc("/api/v1/market/orders",r.market.Orders); mux.HandleFunc("/api/v1/market/sell",r.market.Sell); mux.HandleFunc("/api/v1/market/buy",r.market.Buy); mux.HandleFunc("/api/v1/market/cancel",r.market.Cancel) }; if r.missions!=nil { mux.HandleFunc("/api/v1/missions",r.missions.List); mux.HandleFunc("/api/v1/missions/accept",r.missions.Accept); mux.HandleFunc("/api/v1/missions/progress",func(w http.ResponseWriter,_ *http.Request){writeJSON(w,http.StatusForbidden,map[string]string{"error":"mission progress is server-issued"})}); mux.HandleFunc("/api/v1/missions/complete",r.missions.Complete) }; if r.combat!=nil { mux.HandleFunc("/api/v1/combat/resolve",r.combat.Resolve); mux.HandleFunc("/api/v1/combat/history",r.combat.History) }; if r.equipment!=nil { mux.HandleFunc("/api/v1/equipment/items",r.equipment.Items); mux.HandleFunc("/api/v1/equipment",r.equipment.Equipped); mux.HandleFunc("/api/v1/equipment/equip",r.equipment.Equip) }; if r.worldEvents!=nil { mux.HandleFunc("/api/v1/world-events",r.worldEvents.Recent); mux.HandleFunc("/api/v1/world-events/definitions",r.worldEvents.Definitions) }; if r.news!=nil { mux.HandleFunc("/api/v1/news",r.news.Recent) }; if r.achievements!=nil { mux.HandleFunc("/api/v1/achievements",r.achievements.List); mux.HandleFunc("/api/v1/rankings",r.achievements.Rankings) }; if r.admin!=nil { mux.HandleFunc("/api/v1/admin/world/flags",r.admin.Flags); mux.HandleFunc("/api/v1/admin/world/flags/set",r.admin.SetFlag) };if r.territory!=nil { mux.HandleFunc("/api/v1/territories",r.territory.List); mux.HandleFunc("/api/v1/territory",r.territory.Get); mux.HandleFunc("/api/v1/territories/influence",r.territory.AddInfluence) };if r.organizationsPersistent!=nil { mux.HandleFunc("/api/v1/organizations",r.organizationsPersistent.List); mux.HandleFunc("/api/v1/organizations/join",r.organizationsPersistent.Join); mux.HandleFunc("/api/v1/organizations/reputation",r.organizationsPersistent.Reputation); mux.HandleFunc("/api/v1/organizations/alliance",r.organizationsPersistent.Alliance); mux.HandleFunc("/api/v1/organizations/activity",r.organizationsPersistent.Activity) };mux.HandleFunc("/api/v1/events",ehv.recent);mux.HandleFunc("/api/v1/notifications",ehv.notifications);mux.HandleFunc("/api/v1/notifications/read",ehv.markNotificationRead);mux.HandleFunc("/api/v1/notifications/read-all",ehv.markNotificationsRead)
 if r.progressionDB!=nil { mux.HandleFunc("/api/v1/progression/xp",func(w http.ResponseWriter,_ *http.Request){writeJSON(w,http.StatusForbidden,map[string]string{"error":"progression changes are server-issued"})}) }
 ah:=newAuthHandler(r.auth);mux.HandleFunc("/api/v1/auth/register",ah.register);mux.HandleFunc("/api/v1/auth/login",ah.login);mux.HandleFunc("/api/v1/auth/me",ah.me);mux.HandleFunc("/api/v1/auth/logout",ah.logout)
 return securityHeaders(RateLimitMiddleware(NewRateLimiter(120, time.Minute),mux))
}
func(r *Router)health(w http.ResponseWriter,_ *http.Request){writeJSON(w,http.StatusOK,map[string]any{"status":"ok","service":"worldbuster","time":time.Now().UTC()})}
func(r *Router)worldState(w http.ResponseWriter,_ *http.Request){writeJSON(w,http.StatusOK,r.world.Snapshot())}
func writeJSON(w http.ResponseWriter,status int,value any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_=json.NewEncoder(w).Encode(value)}
func securityHeaders(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Header().Set("X-Content-Type-Options","nosniff");w.Header().Set("X-Frame-Options","DENY");w.Header().Set("Referrer-Policy","strict-origin-when-cross-origin");next.ServeHTTP(w,r)})}


func(r *Router) WithProgressionRepository(api *ProgressionAPI)*Router {
 api.Auth=r.auth
 r.progressionDB=api
 return r
}

func(r *Router) WithPlayerResolver(resolver PlayerContextResolver)*Router { r.playerResolver=resolver; if r.progressionDB!=nil { r.progressionDB.ResolvePlayer=resolver.PlayerID }; return r }

func(r *Router) WithPlayerContext(pc *PlayerContext)*Router { r.playerContext=pc; return r }

func(r *Router) WithPersistentJobs(api *PersistentJobsAPI)*Router { r.persistentJobs=api; return r }

func(r *Router) WithEducation(api *EducationAPI)*Router { r.education=api; return r }

func(r *Router) WithPersistentSocial(api *PersistentSocialAPI)*Router { r.socialPersistent=api; return r }

func(r *Router) WithReputation(api *ReputationAPI)*Router { r.reputation=api; return r }

func(r *Router) WithPersistentOrganizations(api *PersistentOrganizationsAPI)*Router { r.organizationsPersistent=api; return r }

func(r *Router) WithTerritory(api *TerritoryAPI)*Router { r.territory=api; return r }

func(r *Router) WithMarket(api *MarketAPI)*Router { r.market=api; return r }

func(r *Router) WithMissions(api *MissionsAPI)*Router { r.missions=api; return r }

func(r *Router) WithCombat(api *CombatAPI)*Router { r.combat=api; return r }

func(r *Router) WithEquipment(api *EquipmentAPI)*Router { r.equipment=api; return r }

func(r *Router) WithWorldEvents(api *WorldEventsAPI)*Router { r.worldEvents=api; return r }

func(r *Router) WithNews(api *NewsAPI)*Router { r.news=api; return r }

func(r *Router) WithAchievements(api *AchievementsAPI)*Router { r.achievements=api; return r }

func(r *Router) WithAdmin(api *AdminAPI)*Router { r.admin=api; return r }
