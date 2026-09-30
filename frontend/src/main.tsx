import React,{useEffect,useState} from "react";
import {createRoot} from "react-dom/client";
import {AuthPanel} from "./auth";
import "./styles.css";

type WorldState={tick:number;onlineCount:number;day:number;time:string;status:string};
type Account={id:string;username:string};
const fallback:WorldState={tick:0,onlineCount:0,day:1,time:"00:00",status:"CONNECTING"};

function App(){
 const [world,setWorld]=useState(fallback);
 const [account,setAccount]=useState<Account|null>(null);
 const [checking,setChecking]=useState(true);

 useEffect(()=>{
  fetch("/api/v1/auth/me",{credentials:"include"}).then(async r=>r.ok?setAccount({id:(await r.json()).accountId,username:"PLAYER"}):null).finally(()=>setChecking(false));
 },[]);

 useEffect(()=>{
  if(!account)return;
  const load=()=>fetch("/api/v1/world").then(r=>r.json()).then(setWorld).catch(()=>{});
  load();const id=setInterval(load,1000);return()=>clearInterval(id);
 },[account]);

 if(checking)return <div className="auth-stage"><div className="kicker">WORLDBUSTER</div><h1>SYNCING IDENTITY.</h1></div>;
 if(!account)return <div className="auth-stage"><AuthPanel onAuthenticated={setAccount}/></div>;

 return <div className="app">
  <header className="topbar"><div className="brand"><span className="brand-mark">W</span><span>WORLDBUSTER</span></div><div className="world-status"><i/>{world.status}<b>DAY {world.day}</b></div></header>
  <section className="hero"><div className="kicker">IDENTITY CONFIRMED · {account.username}</div><h1>ENTER THE<br/><em>LIVING WORLD.</em></h1><p>Every action changes something. Every system feeds another system.</p><div className="actions"><button>ENTER WORLD</button><button className="ghost">EXPLORE SYSTEMS</button></div></section>
  <section className="dashboard">
   <article><span>WORLD CLOCK</span><strong>{world.time}</strong><small>TICK {world.tick}</small></article>
   <article><span>WORLD STATE</span><strong>STABLE</strong><small>SERVER AUTHORITATIVE</small></article>
   <article><span>SIMULATION</span><strong>READY</strong><small>EVENT ENGINE ONLINE</small></article>
  </section>
 </div>;
}
createRoot(document.getElementById("root")!).render(<React.StrictMode><App/></React.StrictMode>);
