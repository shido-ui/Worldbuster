import React,{useEffect,useRef,useState} from "react";
import {createRoot} from "react-dom/client";
import {AuthPanel} from "./auth";
import "./styles.css";

type Account={id:string;username:string};
type State={profile:{displayName:string;level:number;xp:number;cash:number;energy:number;strength:number;defense:number;speed:number;intelligence:number;endurance:number};economy:{balance:number;currency:string};inventory:{stacks:{itemId:string;quantity:number}[]}};
type World={tick:number;onlineCount:number;day:number;time:string;status:string};
type Job={id:string;name:string;department:string;baseSalary:number;requiredLevel:number;requiredStat:number;requiredEducation:number};
type Course={id:string;name:string;durationMinutes:number;educationGain:number;requiredLevel:number};
type ConnectionStatus="LOADING"|"CONNECTED"|"AUTHENTICATION_REQUIRED"|"SERVER_ERROR"|"OFFLINE";
const worldFallback:World={tick:0,onlineCount:0,day:1,time:"00:00",status:"CONNECTING"};

async function fetchJSON(url:string,init:RequestInit={},signal?:AbortSignal){
 const response=await fetch(url,{...init,credentials:"include",signal});
 const data=await response.json().catch(()=>null);
 if(!response.ok){const error=new Error(data?.error??`HTTP ${response.status}`);(error as any).status=response.status;throw error}
 return data;
}

function App(){
 const [view,setView]=useState("Overview"),[account,setAccount]=useState<Account|null>(null),[checking,setChecking]=useState(true),[world,setWorld]=useState(worldFallback),[state,setState]=useState<State|null>(null),[jobs,setJobs]=useState<Job[]>([]),[employment,setEmployment]=useState<Job|null>(null),[courses,setCourses]=useState<Course[]>([]),[training,setTraining]=useState<any>(null),[education,setEducation]=useState(0),[courseBusy,setCourseBusy]=useState(false),[jobBusy,setJobBusy]=useState(false),[jobMessage,setJobMessage]=useState(""),[inbox,setInbox]=useState<any[]>([]),[reputation,setReputation]=useState<any>(null),[organizations,setOrganizations]=useState<any[]>([]),[connection,setConnection]=useState<ConnectionStatus>("LOADING"),[connectionError,setConnectionError]=useState("");
 const refreshSequence=useRef(0);

 const refresh=async(signal?:AbortSignal)=>{
  const sequence=++refreshSequence.current;
  try{
   const results=await Promise.all([
    fetchJSON("/api/v1/world",{},signal),
    fetchJSON("/api/v1/player/state",{},signal),
    fetchJSON("/api/v1/jobs",{},signal),
    fetchJSON("/api/v1/jobs/status",{},signal),
    fetchJSON("/api/v1/education/courses",{},signal),
    fetchJSON("/api/v1/education/status",{},signal),
    fetchJSON("/api/v1/social/inbox",{},signal),
    fetchJSON("/api/v1/reputation",{},signal),
    fetchJSON("/api/v1/organizations",{},signal)
   ]);
   if(sequence!==refreshSequence.current)return;
   const [w,s,j,e,coursesData,trainingData,si,rep,orgs]=results as any[];
   setWorld(w);setState(s);setJobs(j.jobs??j);setEmployment(e.job??null);
   setCourses(coursesData.courses??coursesData);
   setEducation(trainingData.education??0);setTraining(trainingData.training??null);
   setInbox(si.messages??[]);setReputation(rep);setOrganizations(orgs.organizations??orgs);
   setConnection("CONNECTED");setConnectionError("");
  }catch(error){
   if((error as any)?.name==="AbortError")return;
   if(sequence!==refreshSequence.current)return;
   const status=(error as any)?.status;
   if(status===401){setAccount(null);setState(null);setConnection("AUTHENTICATION_REQUIRED");setConnectionError("Session expired. Please sign in again.");}
   else if(status>=500){setConnection("SERVER_ERROR");setConnectionError("World server is unavailable.");}
   else if(error instanceof TypeError){setConnection("OFFLINE");setConnectionError("Network connection unavailable.");}
   else {setConnection("SERVER_ERROR");setConnectionError((error as Error).message);}
  }
 };

 useEffect(()=>{
  let cancelled=false;
  fetchJSON("/api/v1/auth/me").then(data=>{if(!cancelled)setAccount({id:data.accountId,username:"PLAYER"})}).catch(error=>{if(!cancelled&&((error as any)?.status===401)){setConnection("AUTHENTICATION_REQUIRED")}}).finally(()=>{if(!cancelled)setChecking(false)});
  return()=>{cancelled=true};
 },[]);

 useEffect(()=>{
  if(!account)return;
  const controller=new AbortController();
  void refresh(controller.signal);
  const id=setInterval(()=>{void refresh(controller.signal)},5000);
  return()=>{controller.abort();clearInterval(id)};
 },[account]);

 const enroll=async(courseId:string)=>{setCourseBusy(true);try{const d=await fetchJSON("/api/v1/education/enroll",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({courseId})});setTraining(d.training);setJobMessage("Training started.")}catch(error){setJobMessage((error as Error).message)}finally{setCourseBusy(false)}};
 const hire=async(jobId:string)=>{setJobBusy(true);setJobMessage("");try{const d=await fetchJSON("/api/v1/jobs/employ",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({jobId})});setEmployment(d.job);setJobMessage("Employment confirmed.");await refresh()}catch(error){setJobMessage((error as Error).message)}finally{setJobBusy(false)}};
 if(checking)return <div className="auth-stage"><div className="kicker">WORLDBUSTER</div><h1>SYNCING IDENTITY.</h1></div>;
 if(!account)return <div className="auth-stage"><AuthPanel onAuthenticated={setAccount}/></div>;
 return <div className="app">
  <header className="topbar"><div className="brand"><span className="brand-mark">W</span>WORLDBUSTER</div><div className="world-status"><i/>{connection}<b>{connection==="CONNECTED"?`DAY ${world.day}`:"NO LIVE WORLD"}</b></div></header>
  <main className="game-shell">{connection!=="CONNECTED"&&<div className="job-message">{connectionError||"Synchronizing with the authoritative server..."}</div>}
   <aside className="side"><div className="side-title">COMMAND</div><button className={"nav "+(view==="Overview"?"active":"")} onClick={()=>setView("Overview")}>Overview</button><button className={"nav "+(view==="Character"?"active":"")} onClick={()=>setView("Character")}>Character</button><button className={"nav "+(view==="Inventory"?"active":"")} onClick={()=>setView("Inventory")}>Inventory</button><button className={"nav "+(view==="Economy"?"active":"")} onClick={()=>setView("Economy")}>Economy</button><button className={"nav "+(view==="World"?"active":"")} onClick={()=>setView("World")}>World</button><button className={"nav "+(view==="Jobs"?"active":"")} onClick={()=>setView("Jobs")}>Jobs</button><button className={"nav "+(view==="Education"?"active":"")} onClick={()=>setView("Education")}>Education</button><button className={"nav "+(view==="Social"?"active":"")} onClick={()=>setView("Social")}>Social</button><button className={"nav "+(view==="Reputation"?"active":"")} onClick={()=>setView("Reputation")}>Reputation</button><button className={"nav "+(view==="Organizations"?"active":"")} onClick={()=>setView("Organizations")}>Organizations</button></aside>
   <section className="content">
    <div className="kicker">LIVE PLAYER STATE · {state?.profile.displayName??account.username}</div>
    <div className="title-row"><div><h2>COMMAND CENTER</h2><p>One persistent identity. One living world.</p></div><button onClick={refresh}>SYNC</button></div>
    {state&&<>{view==="Inventory"&&state&&<section className="panel module"><div className="panel-head"><span>INVENTORY</span><small>PERSISTENT</small></div>{state.inventory.stacks.length?<div className="item-list">{state.inventory.stacks.map(x=><div className="item" key={x.itemId}><b>{x.itemId}</b><span>x{x.quantity}</span></div>)}</div>:<p>No items yet.</p>}</section>}
{view==="Economy"&&state&&<section className="panel module"><div className="panel-head"><span>ECONOMY</span><small>{state.economy.currency}</small></div><div className="clock">{state.economy.balance.toLocaleString()}</div><p>Persistent player balance.</p></section>}
{view==="Character"&&state&&<section className="panel module"><div className="panel-head"><span>CHARACTER</span><small>PROFILE</small></div><div className="stats">{Object.entries({Strength:state.profile.strength,Defense:state.profile.defense,Speed:state.profile.speed,Intelligence:state.profile.intelligence,Endurance:state.profile.endurance}).map(([k,v])=><div className="stat" key={k}><span>{k}</span><b>{v}</b></div>)}</div></section>}
{view==="World"&&<section className="panel module"><div className="panel-head"><span>WORLD</span><small>LIVE</small></div><div className="clock">{world.time}</div><p>Day {world.day} · Tick {world.tick} · {world.onlineCount} active connections.</p></section>}
{view==="Education"&&<section className="panel module"><div className="panel-head"><span>ACADEMY</span><small>EDUCATION {education}</small></div>{training&&<div className="employment-banner"><span>IN TRAINING</span><b>{training.courseName}</b><small>Completes {new Date(training.completesAt).toLocaleString()}</small></div>}<div className="job-list">{courses.map(c=><div className="job" key={c.id}><div><b>{c.name}</b><small>{c.durationMinutes} min · +{c.educationGain} education · Level {c.requiredLevel}+</small></div><button disabled={!!training||courseBusy||((state?.profile.level??0)<c.requiredLevel)} onClick={()=>enroll(c.id)}>{training?"TRAINING":((state?.profile.level??0)<c.requiredLevel?"LOCKED":"ENROLL")}</button></div>)}</div></section>}{view==="Jobs"&&<section className="panel module"><div className="panel-head"><span>CAREER MARKET</span><small>PERSISTENT WORLD</small></div>{employment&&<div className="employment-banner"><span>EMPLOYED</span><b>{employment.name}</b><small>{employment.department} · {employment.baseSalary.toLocaleString()} WBX / cycle</small></div>}{jobMessage&&<div className="job-message">{jobMessage}</div>}<div className="job-list">{jobs.map(j=><div className="job" key={j.id}><div><b>{j.name}</b><small>{j.department} · Level {j.requiredLevel}+ · {j.requiredStat} total stats · Education {j.requiredEducation}+</small></div><div className="job-action"><strong>{j.baseSalary.toLocaleString()} WBX</strong><button disabled={!!employment||jobBusy||((state?.profile.level??0)<j.requiredLevel)} onClick={()=>hire(j.id)}>{employment?"EMPLOYED":((state?.profile.level??0)<j.requiredLevel?"LOCKED":"HIRE")}</button></div></div>)}{!jobs.length&&<p>No jobs are currently available.</p>}</div></section>}{view==="Organizations"&&<section className="panel module"><div className="panel-head"><span>ORGANIZATIONS</span><small>FACTIONS</small></div><div className="job-list">{organizations.length?organizations.map(o=><div className="job" key={o.id}><div><b>{o.name}</b><small>{o.type} · Level {o.level} · Reputation {o.reputation} · {o.maxMembers} member capacity</small></div><strong>{o.treasury.toLocaleString()} WBX</strong></div>):<p>No organizations are available.</p>}</div></section>}{view==="Reputation"&&<section className="panel module"><div className="panel-head"><span>REPUTATION</span><small>CONSEQUENCES</small></div>{reputation?<div className="stats"><div className="stat"><span>Public</span><b>{reputation.publicScore}</b></div><div className="stat"><span>Trust</span><b>{reputation.trustScore}</b></div><div className="stat"><span>Notoriety</span><b>{reputation.notorietyScore}</b></div></div>:<p>Reputation data unavailable.</p>}</section>}{view==="Social"&&<section className="panel module"><div className="panel-head"><span>SOCIAL NETWORK</span><small>PERSISTENT</small></div><div className="item-list">{inbox.length?inbox.map(m=><div className="item" key={m.id}><div><b>Message</b><small>{m.fromId}</small></div><span>{new Date(m.createdAt).toLocaleString()}</span></div>):<p>No messages yet.</p>}</div></section>}
<section className="cards">
      <article><span>LEVEL</span><strong>{state.profile.level}</strong><small>{state.profile.xp} XP</small></article>
      <article><span>ENERGY</span><strong>{state.profile.energy}</strong><small>AVAILABLE</small></article>
      <article><span>BALANCE</span><strong>{state.economy.balance.toLocaleString()}</strong><small>{state.economy.currency}</small></article>
      <article><span>INVENTORY</span><strong>{state.inventory.stacks.length}</strong><small>ITEM STACKS</small></article>
    </section>
    <section className="grid">
      <article className="panel"><div className="panel-head"><span>CHARACTER STATS</span><small>AUTHORITATIVE</small></div><div className="stats">{Object.entries({Strength:state.profile.strength,Defense:state.profile.defense,Speed:state.profile.speed,Intelligence:state.profile.intelligence,Endurance:state.profile.endurance}).map(([k,v])=><div className="stat" key={k}><span>{k}</span><b>{v}</b></div>)}</div></article>
      <article className="panel"><div className="panel-head"><span>WORLD CLOCK</span><small>SERVER</small></div><div className="clock">{world.time}</div><p>Tick {world.tick} · Day {world.day} · {world.onlineCount} active connections</p></article>
    </section></>}
   </section>
  </main>
 </div>
}
createRoot(document.getElementById("root")!).render(<React.StrictMode><App/></React.StrictMode>);
