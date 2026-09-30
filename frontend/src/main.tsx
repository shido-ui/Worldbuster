import React,{useEffect,useState} from "react";
import {createRoot} from "react-dom/client";
import {AuthPanel} from "./auth";
import "./styles.css";

type Account={id:string;username:string};
type State={profile:{displayName:string;level:number;xp:number;cash:number;energy:number;strength:number;defense:number;speed:number;intelligence:number;endurance:number};economy:{balance:number;currency:string};inventory:{stacks:{itemId:string;quantity:number}[]}};
type World={tick:number;onlineCount:number;day:number;time:string;status:string};
const worldFallback:World={tick:0,onlineCount:0,day:1,time:"00:00",status:"CONNECTING"};

function App(){
 const [account,setAccount]=useState<Account|null>(null),[checking,setChecking]=useState(true),[world,setWorld]=useState(worldFallback),[state,setState]=useState<State|null>(null);
 const refresh=()=>Promise.all([
  fetch("/api/v1/world",{credentials:"include"}).then(r=>r.ok?r.json():null),
  fetch("/api/v1/player/state",{credentials:"include"}).then(r=>r.ok?r.json():null)
 ]).then(([w,s])=>{if(w)setWorld(w);if(s)setState(s)}).catch(()=>{});
 useEffect(()=>{fetch("/api/v1/auth/me",{credentials:"include"}).then(async r=>r.ok?setAccount({id:(await r.json()).accountId,username:"PLAYER"}):null).finally(()=>setChecking(false))},[]);
 useEffect(()=>{if(!account)return;refresh();const id=setInterval(refresh,5000);return()=>clearInterval(id)},[account]);
 if(checking)return <div className="auth-stage"><div className="kicker">WORLDBUSTER</div><h1>SYNCING IDENTITY.</h1></div>;
 if(!account)return <div className="auth-stage"><AuthPanel onAuthenticated={setAccount}/></div>;
 return <div className="app">
  <header className="topbar"><div className="brand"><span className="brand-mark">W</span>WORLDBUSTER</div><div className="world-status"><i/>{world.status}<b>DAY {world.day}</b></div></header>
  <main className="game-shell">
   <aside className="side"><div className="side-title">COMMAND</div><button className="nav active">Overview</button><button className="nav">Character</button><button className="nav">Inventory</button><button className="nav">Economy</button><button className="nav">World</button></aside>
   <section className="content">
    <div className="kicker">LIVE PLAYER STATE · {state?.profile.displayName??account.username}</div>
    <div className="title-row"><div><h2>COMMAND CENTER</h2><p>One persistent identity. One living world.</p></div><button onClick={refresh}>SYNC</button></div>
    {state&&<><section className="cards">
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
