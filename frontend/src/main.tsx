import React,{useEffect,useMemo,useRef,useState}from"react";
import{createRoot}from"react-dom/client";
import{AuthPanel}from"./auth";
import"./styles.css";

type Account={id:string;username:string};
type State={profile:{displayName:string;level:number;xp:number;cash:number;energy:number;strength:number;defense:number;speed:number;intelligence:number;endurance:number};economy:{balance:number;currency:string};inventory:{stacks:{itemId:string;quantity:number}[]}};
type World={tick:number;onlineCount:number;day:number;time:string;status:string};
type Job={id:string;name:string;department:string;baseSalary:number;requiredLevel:number;requiredStat:number;requiredEducation:number};
type Course={id:string;name:string;durationMinutes:number;educationGain:number;requiredLevel:number};
type EventItem={id:string;type:string;actorId:string;targetId?:string;createdAt:string;payload?:Record<string,unknown>};
type ConnectionStatus="LOADING"|"CONNECTED"|"AUTHENTICATION_REQUIRED"|"SERVER_ERROR"|"OFFLINE";
type Mission={id:string;code:string;title:string;description:string;missionType:string;minLevel:number;rewardCash:number;rewardXP:number;targetType:string;targetValue:number};
type Asset={id:string;symbol?:string;name:string;lastPrice?:number};
type NewsItem={id:string;title:string;summary?:string;body?:string;publishedAt?:string};
type WorldEvent={id:string;type:string;title?:string;description?:string;createdAt?:string};
type Achievement={id:string;code?:string;name?:string;title?:string;description?:string;unlocked?:boolean;completed?:boolean};
type Order={id:string;assetId:string;side:string;quantity:number;unitPrice:number;status?:string};
type Location={id:string;name?:string;description?:string};
const worldFallback:World={tick:0,onlineCount:0,day:1,time:"00:00",status:"CONNECTING"};

async function fetchJSON(url:string,init:RequestInit={},signal?:AbortSignal){
 const response=await fetch(url,{...init,credentials:"include",signal});
 const data=await response.json().catch(()=>null);
 if(!response.ok){const error=new Error(data?.error??"HTTP "+response.status);(error as any).status=response.status;throw error}
 return data;
}
const nav=["Overview","Character","Inventory","Economy","World","Missions","Market","Jobs","Education","Social","Reputation","Organizations","Achievements","News"];
const eventLabel=(e:EventItem)=>e.type.replaceAll("_"," ").toLowerCase().replace(/^./,c=>c.toUpperCase());
const formatTime=(value:string)=>{const d=new Date(value);return Number.isNaN(d.getTime())?"recent":d.toLocaleTimeString([], {hour:"2-digit",minute:"2-digit"})};

function App(){
 const[view,setView]=useState("Overview"),[missions,setMissions]=useState<Mission[]>([]),[assets,setAssets]=useState<Asset[]>([]),[orders,setOrders]=useState<Order[]>([]),[news,setNews]=useState<NewsItem[]>([]),[worldEvents,setWorldEvents]=useState<WorldEvent[]>([]),[achievements,setAchievements]=useState<Achievement[]>([]),[locations,setLocations]=useState<Location[]>([]),[account,setAccount]=useState<Account|null>(null),[checking,setChecking]=useState(true),[world,setWorld]=useState(worldFallback),[state,setState]=useState<State|null>(null),[jobs,setJobs]=useState<Job[]>([]),[employment,setEmployment]=useState<Job|null>(null),[courses,setCourses]=useState<Course[]>([]),[training,setTraining]=useState<any>(null),[education,setEducation]=useState(0),[courseBusy,setCourseBusy]=useState(false),[jobBusy,setJobBusy]=useState(false),[message,setMessage]=useState(""),[inbox,setInbox]=useState<any[]>([]),[reputation,setReputation]=useState<any>(null),[organizations,setOrganizations]=useState<any[]>([]),[events,setEvents]=useState<EventItem[]>([]),[connection,setConnection]=useState<ConnectionStatus>("LOADING"),[connectionError,setConnectionError]=useState(""),[mobileNav,setMobileNav]=useState(false);
 const refreshSequence=useRef(0);
 async function optional<T>(url:string,signal?:AbortSignal):Promise<T|null>{try{return await fetchJSON(url,{},signal) as T}catch(error){if((error as any)?.name==="AbortError")throw error;return null}}
 const refresh=async(signal?:AbortSignal)=>{
  const sequence=++refreshSequence.current;
  try{
   const results=await Promise.all([
    fetchJSON("/api/v1/world",{},signal),fetchJSON("/api/v1/player/state",{},signal),fetchJSON("/api/v1/jobs",{},signal),
    fetchJSON("/api/v1/jobs/status",{},signal),fetchJSON("/api/v1/education/courses",{},signal),fetchJSON("/api/v1/education/status",{},signal),
    fetchJSON("/api/v1/social/inbox",{},signal),fetchJSON("/api/v1/reputation",{},signal),fetchJSON("/api/v1/organizations",{},signal),
    fetchJSON("/api/v1/events?limit=20",{},signal),
    optional<any>("/api/v1/missions",signal),optional<any>("/api/v1/market/assets",signal),optional<any>("/api/v1/market/orders?assetId=WBX",signal),
    optional<any>("/api/v1/news?limit=12",signal),optional<any>("/api/v1/world-events",signal),optional<any>("/api/v1/achievements",signal),optional<any>("/api/v1/world/locations",signal)
   ]);
   if(sequence!==refreshSequence.current)return;
   const[w,s,j,e,coursesData,trainingData,si,rep,orgs,eventsData,missionData,assetData,orderData,newsData,worldEventData,achievementData,locationData]=results as any[];
   setWorld(w);setState(s);setJobs(j.jobs??j);setEmployment(e.job??null);setCourses(coursesData.courses??coursesData);
   setEducation(trainingData.education??0);setTraining(trainingData.training??null);setInbox(si.messages??[]);setReputation(rep);
   setOrganizations(orgs.organizations??orgs);setEvents(Array.isArray(eventsData)?eventsData:(eventsData.events??[]));setMissions(missionData?.missions??[]);setAssets(assetData?.assets??[]);setOrders(orderData?.orders??[]);setNews(newsData?.news??[]);setWorldEvents(worldEventData?.events??[]);setAchievements(achievementData?.achievements??[]);setLocations(locationData?.locations??[]);
   setConnection("CONNECTED");setConnectionError("");
  }catch(error){
   if((error as any)?.name==="AbortError")return;
   if(sequence!==refreshSequence.current)return;
   const status=(error as any)?.status;
   if(status===401){setAccount(null);setState(null);setConnection("AUTHENTICATION_REQUIRED");setConnectionError("Session expired. Please sign in again.")}
   else if(status>=500){setConnection("SERVER_ERROR");setConnectionError("World server is unavailable.")}
   else if(error instanceof TypeError){setConnection("OFFLINE");setConnectionError("Network connection unavailable.")}
   else{setConnection("SERVER_ERROR");setConnectionError((error as Error).message)}
  }
 };
 useEffect(()=>{let cancelled=false;fetchJSON("/api/v1/auth/me").then(data=>{if(!cancelled)setAccount({id:data.accountId,username:"PLAYER"})}).catch(error=>{if(!cancelled&&(error as any)?.status===401)setConnection("AUTHENTICATION_REQUIRED")}).finally(()=>{if(!cancelled)setChecking(false)});return()=>{cancelled=true}},[]);
 useEffect(()=>{if(!account)return;const controller=new AbortController();void refresh(controller.signal);const id=setInterval(()=>void refresh(controller.signal),10000);return()=>{controller.abort();clearInterval(id)}},[account]);
 useEffect(()=>{const onKey=(e:KeyboardEvent)=>{if(e.key==="/"&&document.activeElement?.tagName!=="INPUT"){e.preventDefault();setView("Overview");document.getElementById("command-search")?.focus()}};window.addEventListener("keydown",onKey);return()=>window.removeEventListener("keydown",onKey)},[]);
 const enroll=async(courseId:string)=>{setCourseBusy(true);try{const d=await fetchJSON("/api/v1/education/enroll",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({courseId})});setTraining(d.training);setMessage("Training started.");await refresh()}catch(error){setMessage((error as Error).message)}finally{setCourseBusy(false)}};
 const hire=async(jobId:string)=>{setJobBusy(true);setMessage("");try{const d=await fetchJSON("/api/v1/jobs/employ",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({jobId})});setEmployment(d.job);setMessage("Employment confirmed.");await refresh()}catch(error){setMessage((error as Error).message)}finally{setJobBusy(false)}};
 const latestEvents=useMemo(()=>events.slice(0,6),[events]);
 const acceptMission=async(id:string)=>{try{await fetchJSON("/api/v1/missions/accept",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({missionId:id})});setMessage("Mission accepted.");await refresh()}catch(error){setMessage((error as Error).message)}};
 if(checking)return <div className="auth-stage"><div className="boot"><div className="brand-mark">W</div><span>WORLDBUSTER</span><b>ESTABLISHING SECURE SESSION</b><div className="boot-line"/></div></div>;
 if(!account)return <div className="auth-stage"><AuthPanel onAuthenticated={setAccount}/></div>;
 const statusClass=connection.toLowerCase();
 return <div className="app">
  <header className="topbar"><button className="mobile-menu" aria-label="Open navigation" onClick={()=>setMobileNav(v=>!v)}>☰</button><div className="brand"><span className="brand-mark">W</span><span>WORLDBUSTER</span></div><div className={"world-status "+statusClass}><i/><span>{connection==="CONNECTED"?"LIVE WORLD":connection.replaceAll("_"," ")}</span><b>{connection==="CONNECTED"?"DAY "+world.day:"NO LIVE WORLD"}</b></div></header>
  <main className="game-shell">
   <aside className={"side "+(mobileNav?"open":"")}><div className="side-title">COMMAND</div>{nav.map(item=><button key={item} className={"nav "+(view===item?"active":"")} onClick={()=>{setView(item);setMobileNav(false)}}><span>{item}</span>{item==="Social"&&inbox.length>0&&<em>{inbox.length}</em>}</button>)}</aside>
   <section className="content">
    <div className="command-bar"><div className="kicker">PERSISTENT WORLD · {state?.profile.displayName??account.username}</div><label className="command-search"><span>⌕</span><input id="command-search" placeholder="Jump to a system…" onKeyDown={e=>{if(e.key==="Enter"){const q=e.currentTarget.value.toLowerCase();const hit=nav.find(x=>x.toLowerCase().includes(q));if(hit)setView(hit);e.currentTarget.value=""}}}/><kbd>/</kbd></label></div>
    {connection!=="CONNECTED"&&<div className={"status-banner "+statusClass}><i/><span>{connectionError||"Synchronizing with the authoritative server…"}</span><button onClick={()=>void refresh()}>RETRY</button></div>}
    {message&&<div className="toast" role="status">{message}<button onClick={()=>setMessage("")}>×</button></div>}
    <div className="title-row"><div><div className="kicker">COMMAND CENTER</div><h1>THE WORLD IS MOVING.</h1><p>Your state is authoritative. Actions create persistent consequences.</p></div><button className="sync" onClick={()=>void refresh()}><span>↻</span> SYNC</button></div>
    {state&&<section className="cards">
      <article className="metric"><span>LEVEL</span><strong>{state.profile.level}</strong><small>{state.profile.xp} XP</small></article>
      <article className="metric"><span>ENERGY</span><strong>{state.profile.energy}</strong><small>AVAILABLE</small></article>
      <article className="metric"><span>BALANCE</span><strong>{state.economy.balance.toLocaleString()}</strong><small>{state.economy.currency}</small></article>
      <article className="metric accent"><span>WORLD ONLINE</span><strong>{world.onlineCount.toLocaleString()}</strong><small>DAY {world.day} · TICK {world.tick}</small></article>
    </section>}
    {view==="Overview"&&state&&<><section className="dashboard-grid">
      <article className="panel hero-panel"><div className="panel-head"><span>WORLD CLOCK</span><small>SERVER AUTHORITY</small></div><div className="clock">{world.time}</div><div className="world-line"><span>DAY {world.day}</span><i/><span>TICK {world.tick}</span><i/><span>{world.status}</span></div></article>
      <article className="panel"><div className="panel-head"><span>CHARACTER</span><small>PROFILE</small></div><div className="profile-row"><div className="avatar">W</div><div><b>{state.profile.displayName}</b><small>Level {state.profile.level} · {state.profile.xp} XP</small></div></div><div className="mini-stats">{[["STR",state.profile.strength],["DEF",state.profile.defense],["SPD",state.profile.speed],["INT",state.profile.intelligence],["END",state.profile.endurance]].map(([k,v])=><div key={String(k)}><span>{k}</span><b>{v}</b></div>)}</div></article>
     </section><section className="dashboard-grid lower"><article className="panel"><div className="panel-head"><span>LIVE ACTIVITY</span><small>RECENT EVENTS</small></div><div className="timeline">{latestEvents.length?latestEvents.map(e=><div className="timeline-item" key={e.id}><i/><div><b>{eventLabel(e)}</b><small>{e.actorId||"world"} · {formatTime(e.createdAt)}</small></div></div>):<div className="empty"><b>No world events yet.</b><span>The event stream will appear here as the simulation moves.</span></div>}</div></article>
     <article className="panel"><div className="panel-head"><span>OPPORTUNITIES</span><small>YOUR STATE</small></div><div className="opportunity"><div><b>{employment?employment.name:"Find employment"}</b><span>{employment?employment.baseSalary.toLocaleString()+" WBX / cycle":"Career progression is available."}</span></div><button onClick={()=>setView("Jobs")}>{employment?"VIEW":"EXPLORE"}</button></div><div className="opportunity"><div><b>{training?training.courseName:"Continue education"}</b><span>{training?"Training in progress":"Build education and unlock roles."}</span></div><button onClick={()=>setView("Education")}>{training?"VIEW":"STUDY"}</button></div><div className="opportunity"><div><b>{inbox.length?inbox.length+" messages":"Social network"}</b><span>Relationships and reputation persist.</span></div><button onClick={()=>setView("Social")}>OPEN</button></div></article>
    </section></>}
    {view==="Missions"&&<Module title="MISSIONS" eyebrow="SERVER-AUTHORITATIVE OBJECTIVES">{missions.length?<div className="job-list">{missions.map(m=><div className="job" key={m.id}><div><b>{m.title}</b><small>{m.description} · {m.missionType} · Reward {m.rewardXP} XP / {m.rewardCash} WBX</small></div><button onClick={()=>acceptMission(m.id)}>ACCEPT</button></div>)}</div>:<Empty title="No active missions" text="The world has no missions available for your current level."/>}</Module>}
    {view==="Market"&&<Module title="MARKET" eyebrow="LIVE ECONOMY">{assets.length?<div className="job-list">{assets.map(a=><div className="job" key={a.id}><div><b>{a.name}</b><small>{a.symbol??a.id} · Server-priced asset</small></div><strong>{a.lastPrice==null?"—":a.lastPrice.toLocaleString()} WBX</strong></div>)}</div>:<Empty title="Market unavailable" text="No market assets were returned by the authoritative server."/>}</Module>}
    {view==="Achievements"&&<Module title="ACHIEVEMENTS" eyebrow="PERSISTENT PROGRESSION">{achievements.length?<div className="job-list">{achievements.map(a=><div className="job" key={a.id}><div><b>{a.title??a.name??a.code??a.id}</b><small>{a.description??"Persistent achievement"}</small></div><span>{a.unlocked||a.completed?"UNLOCKED":"LOCKED"}</span></div>)}</div>:<Empty title="No achievements yet" text="Achievements will appear as the world records your progress."/>}</Module>}
    {view==="News"&&<Module title="WORLD NEWS" eyebrow="PERSISTENT WORLD FEED">{news.length?<div className="job-list">{news.map(n=><article className="job" key={n.id}><div><b>{n.title}</b><small>{n.summary??n.body??"World report"}{n.publishedAt?" · "+formatTime(n.publishedAt):""}</small></div></article>)}</div>:<Empty title="No news" text="The world press has no reports available yet."/>}{worldEvents.length>0&&<><div className="panel-head" style={{marginTop:24}}><span>WORLD EVENTS</span><small>{worldEvents.length} RECENT</small></div><div className="timeline">{worldEvents.slice(0,8).map(e=><div className="timeline-item" key={e.id}><i/><div><b>{e.title??e.type}</b><small>{e.description??e.type}{e.createdAt?" · "+formatTime(e.createdAt):""}</small></div></div>)}</div></>}</Module>}
     {view==="Character"&&state&&<Module title="CHARACTER" eyebrow="AUTHORITATIVE PROFILE">
      <div className="character-head"><div className="avatar large">W</div><div><h2>{state.profile.displayName}</h2><p>Level {state.profile.level} · {state.profile.xp} XP</p></div></div>
      <div className="large-stats">{[["Strength",state.profile.strength],["Defense",state.profile.defense],["Speed",state.profile.speed],["Intelligence",state.profile.intelligence],["Endurance",state.profile.endurance]].map(([k,v])=><div className="large-stat" key={String(k)}><span>{k}</span><strong>{v}</strong><div className="stat-bar"><i style={{width:Math.min(100,Number(v))+"%"}} /></div></div>)}</div>
     </Module>}
    {view==="Inventory"&&state&&<Module title="INVENTORY" eyebrow="PERSISTENT ITEMS">{state.inventory.stacks.length?<div className="item-list">{state.inventory.stacks.map(x=><div className="item" key={x.itemId}><div><b>{x.itemId}</b><small>Persistent item stack</small></div><strong>x{x.quantity}</strong></div>)}</div>:<Empty title="Inventory is empty" text="Items acquired through the world will appear here."/>}</Module>}
    {view==="Economy"&&state&&<Module title="ECONOMY" eyebrow={state.economy.currency}><div className="economy-total">{state.economy.balance.toLocaleString()}<small>{state.economy.currency}</small></div><div className="economy-note">Balances are server-authoritative and persist across sessions.</div></Module>}
    {view==="World"&&<Module title="WORLD" eyebrow="LIVE SIMULATION"><div className="world-card"><div className="clock">{world.time}</div><div className="world-line"><span>DAY {world.day}</span><i/><span>TICK {world.tick}</span><i/><span>{world.onlineCount} ACTIVE</span></div></div><div className="map-grid">{locations.length?locations.map(l=><div className="location-card" key={l.id}><b>{l.name||l.id}</b><small>{l.description||"World location"}</small><span>LOCATION</span></div>):<Empty title="No locations returned" text="The authoritative world has not exposed location data yet."/>}</div><div className="section-note">World navigation and time-based travel are server-authoritative systems.</div></Module>}
    {view==="Jobs"&&<Module title="CAREER MARKET" eyebrow="PERSISTENT WORLD">{employment&&<div className="employment-banner"><span>EMPLOYED</span><b>{employment.name}</b><small>{employment.department} · {employment.baseSalary.toLocaleString()} WBX / cycle</small></div>}{jobs.length?<div className="job-list">{jobs.map(j=><div className="job" key={j.id}><div><b>{j.name}</b><small>{j.department} · Level {j.requiredLevel}+ · {j.requiredStat} total stats · Education {j.requiredEducation}+</small></div><div className="job-action"><strong>{j.baseSalary.toLocaleString()} WBX</strong><button disabled={!!employment||jobBusy||((state?.profile.level??0)<j.requiredLevel)} onClick={()=>hire(j.id)}>{employment?"EMPLOYED":((state?.profile.level??0)<j.requiredLevel?"LOCKED":"HIRE")}</button></div></div>)}</div>:<Empty title="No jobs available" text="The career market has no currently listed positions."/>}</Module>}
    {view==="Education"&&<Module title="ACADEMY" eyebrow={"EDUCATION "+education}>{training&&<div className="employment-banner"><span>IN TRAINING</span><b>{training.courseName}</b><small>Completes {new Date(training.completesAt).toLocaleString()}</small></div>}<div className="job-list">{courses.map(c=><div className="job" key={c.id}><div><b>{c.name}</b><small>{c.durationMinutes} min · +{c.educationGain} education · Level {c.requiredLevel}+</small></div><button disabled={!!training||courseBusy||((state?.profile.level??0)<c.requiredLevel)} onClick={()=>enroll(c.id)}>{training?"TRAINING":((state?.profile.level??0)<c.requiredLevel?"LOCKED":"ENROLL")}</button></div>)}</div></Module>}
    {view==="Social"&&<Module title="SOCIAL NETWORK" eyebrow="PERSISTENT"><div className="item-list">{inbox.length?inbox.map(m=><div className="item" key={m.id}><div><b>Message from {m.fromId}</b><small>{new Date(m.createdAt).toLocaleString()}</small></div><span>OPEN</span></div>):<Empty title="No messages" text="Social events and messages will appear here."/>}</div></Module>}
    {view==="Reputation"&&<Module title="REPUTATION" eyebrow="CONSEQUENCES">{reputation?<div className="large-stats compact">{[["Public",reputation.publicScore],["Trust",reputation.trustScore],["Notoriety",reputation.notorietyScore]].map(([k,v])=><div className="large-stat" key={String(k)}><span>{k}</span><strong>{v}</strong></div>)}</div>:<Empty title="Reputation unavailable" text="No reputation data was returned by the server."/>}</Module>}
    {view==="Organizations"&&<Module title="ORGANIZATIONS" eyebrow="FACTIONS">{organizations.length?<div className="job-list">{organizations.map(o=><div className="job" key={o.id}><div><b>{o.name}</b><small>{o.type} · Level {o.level} · Reputation {o.reputation} · {o.maxMembers} member capacity</small></div><strong>{o.treasury.toLocaleString()} WBX</strong></div>)}</div>:<Empty title="No organizations" text="There are no organizations available in the current world state."/>}</Module>}
   </section>
  </main>
 </div>;
}
function Module({title,eyebrow,children}:{title:string;eyebrow:string;children:React.ReactNode}){return <section className="panel module"><div className="panel-head"><span>{title}</span><small>{eyebrow}</small></div>{children}</section>}
function Empty({title,text}:{title:string;text:string}){return <div className="empty"><b>{title}</b><span>{text}</span></div>}
createRoot(document.getElementById("root")!).render(<React.StrictMode><App/></React.StrictMode>);

