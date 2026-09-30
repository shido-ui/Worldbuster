import React,{useEffect,useState} from "react";
import {createRoot} from "react-dom/client";
import {AuthPanel} from "./auth";
import "./styles.css";

type Account={id:string;username:string};
type State={profile:{displayName:string;level:number;xp:number;cash:number;energy:number;strength:number;defense:number;speed:number;intelligence:number;endurance:number};economy:{balance:number;currency:string};inventory:{stacks:{itemId:string;quantity:number}[]}};
type World={tick:number;onlineCount:number;day:number;time:string;status:string};
type Job={id:string;name:string;department:string;baseSalary:number;requiredLevel:number;requiredStat:number;requiredEducation:number};
type Course={id:string;name:string;durationMinutes:number;educationGain:number;requiredLevel:number};
const worldFallback:World={tick:0,onlineCount:0,day:1,time:"00:00",status:"CONNECTING"};

function App(){
 const [view,setView]=useState("Overview"),[account,setAccount]=useState<Account|null>(null),[checking,setChecking]=useState(true),[world,setWorld]=useState(worldFallback),[state,setState]=useState<State|null>(null),[jobs,setJobs]=useState<Job[]>([]),[employment,setEmployment]=useState<Job|null>(null),[courses,setCourses]=useState<Course[]>([]),[training,setTraining]=useState<any>(null),[education,setEducation]=useState(0),[courseBusy,setJobBusy]=useState(false),[jobMessage,setJobMessage]=useState("");
 const refresh=()=>Promise.all([
  fetch("/api/v1/world",{credentials:"include"}).then(r=>r.ok?r.json():null),
  fetch("/api/v1/player/state",{credentials:"include"}).then(r=>r.ok?r.json():null)
 ]).then(([w,s])=>{if(w)setWorld(w);if(s)setState(s)}).catch(()=>{});
 useEffect(()=>{fetch("/api/v1/auth/me",{credentials:"include"}).then(async r=>r.ok?setAccount({id:(await r.json()).accountId,username:"PLAYER"}):null).finally(()=>setChecking(false))},[]);
 useEffect(()=>{if(!account)return;refresh();const id=setInterval(refresh,5000);return()=>clearInterval(id)},[account]);
 const enroll=async(courseId:string)=>{setCourseBusy(true);try{const r=await fetch("/api/v1/education/enroll",{method:"POST",credentials:"include",headers:{"Content-Type":"application/json"},body:JSON.stringify({courseId})});const d=await r.json();if(!r.ok){setJobMessage(d.error??"Unable to enroll");return}setTraining(d.training);setJobMessage("Training started.")}catch{setJobMessage("Connection error.")}finally{setCourseBusy(false)}};
 const hire=async(jobId:string)=>{setJobBusy(true);setJobMessage("");try{const r=await fetch("/api/v1/jobs/employ",{method:"POST",credentials:"include",headers:{"Content-Type":"application/json"},body:JSON.stringify({jobId})});const d=await r.json();if(!r.ok){setJobMessage(d.error??"Unable to hire");return}setEmployment(d.job);setJobMessage("Employment confirmed.");await refresh()}catch{setJobMessage("Connection error.")}finally{setJobBusy(false)}};
 if(checking)return <div className="auth-stage"><div className="kicker">WORLDBUSTER</div><h1>SYNCING IDENTITY.</h1></div>;
 if(!account)return <div className="auth-stage"><AuthPanel onAuthenticated={setAccount}/></div>;
 return <div className="app">
  <header className="topbar"><div className="brand"><span className="brand-mark">W</span>WORLDBUSTER</div><div className="world-status"><i/>{world.status}<b>DAY {world.day}</b></div></header>
  <main className="game-shell">
   <aside className="side"><div className="side-title">COMMAND</div><button className={"nav "+(view==="Overview"?"active":"")} onClick={()=>setView("Overview")}>Overview</button><button className={"nav "+(view==="Character"?"active":"")} onClick={()=>setView("Character")}>Character</button><button className={"nav "+(view==="Inventory"?"active":"")} onClick={()=>setView("Inventory")}>Inventory</button><button className={"nav "+(view==="Economy"?"active":"")} onClick={()=>setView("Economy")}>Economy</button><button className={"nav "+(view==="World"?"active":"")} onClick={()=>setView("World")}>World</button><button className={"nav "+(view==="Jobs"?"active":"")} onClick={()=>setView("Jobs")}>Jobs</button><button className={"nav "+(view==="Education"?"active":"")} onClick={()=>setView("Education")}>Education</button><button className={"nav "+(view==="Social"?"active":"")} onClick={()=>setView("Social")}>Social</button></aside>
   <section className="content">
    <div className="kicker">LIVE PLAYER STATE · {state?.profile.displayName??account.username}</div>
    <div className="title-row"><div><h2>COMMAND CENTER</h2><p>One persistent identity. One living world.</p></div><button onClick={refresh}>SYNC</button></div>
    {state&&<>{view==="Inventory"&&state&&<section className="panel module"><div className="panel-head"><span>INVENTORY</span><small>PERSISTENT</small></div>{state.inventory.stacks.length?<div className="item-list">{state.inventory.stacks.map(x=><div className="item" key={x.itemId}><b>{x.itemId}</b><span>x{x.quantity}</span></div>)}</div>:<p>No items yet.</p>}</section>}
{view==="Economy"&&state&&<section className="panel module"><div className="panel-head"><span>ECONOMY</span><small>{state.economy.currency}</small></div><div className="clock">{state.economy.balance.toLocaleString()}</div><p>Persistent player balance.</p></section>}
{view==="Character"&&state&&<section className="panel module"><div className="panel-head"><span>CHARACTER</span><small>PROFILE</small></div><div className="stats">{Object.entries({Strength:state.profile.strength,Defense:state.profile.defense,Speed:state.profile.speed,Intelligence:state.profile.intelligence,Endurance:state.profile.endurance}).map(([k,v])=><div className="stat" key={k}><span>{k}</span><b>{v}</b></div>)}</div></section>}
{view==="World"&&<section className="panel module"><div className="panel-head"><span>WORLD</span><small>LIVE</small></div><div className="clock">{world.time}</div><p>Day {world.day} · Tick {world.tick} · {world.onlineCount} active connections.</p></section>}
{view==="Education"&&<section className="panel module"><div className="panel-head"><span>ACADEMY</span><small>EDUCATION {education}</small></div>{training&&<div className="employment-banner"><span>IN TRAINING</span><b>{training.courseName}</b><small>Completes {new Date(training.completesAt).toLocaleString()}</small></div>}<div className="job-list">{courses.map(c=><div className="job" key={c.id}><div><b>{c.name}</b><small>{c.durationMinutes} min · +{c.educationGain} education · Level {c.requiredLevel}+</small></div><button disabled={!!training||courseBusy||((state?.profile.level??0)<c.requiredLevel)} onClick={()=>enroll(c.id)}>{training?"TRAINING":((state?.profile.level??0)<c.requiredLevel?"LOCKED":"ENROLL")}</button></div>)}</div></section>}{view==="Jobs"&&<section className="panel module"><div className="panel-head"><span>CAREER MARKET</span><small>PERSISTENT WORLD</small></div>{employment&&<div className="employment-banner"><span>EMPLOYED</span><b>{employment.name}</b><small>{employment.department} · {employment.baseSalary.toLocaleString()} WBX / cycle</small></div>}{jobMessage&&<div className="job-message">{jobMessage}</div>}<div className="job-list">{jobs.map(j=><div className="job" key={j.id}><div><b>{j.name}</b><small>{j.department} · Level {j.requiredLevel}+ · {j.requiredStat} total stats · Education {j.requiredEducation}+</small></div><div className="job-action"><strong>{j.baseSalary.toLocaleString()} WBX</strong><button disabled={!!employment||jobBusy||((state?.profile.level??0)<j.requiredLevel)} onClick={()=>hire(j.id)}>{employment?"EMPLOYED":((state?.profile.level??0)<j.requiredLevel?"LOCKED":"HIRE")}</button></div></div>)}{!jobs.length&&<p>No jobs are currently available.</p>}</div></section>}{view==="Social"&&<section className="panel module"><div className="panel-head"><span>SOCIAL</span><small>WORLD NETWORK</small></div><p>Messages, relationships, groups and notifications will appear here as their persistent services come online.</p></section>}
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
