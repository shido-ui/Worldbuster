import React, { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import "./styles.css";

type WorldState = { tick:number; onlineCount:number; day:number; time:string; status:string };
const fallback: WorldState = {tick:0,onlineCount:0,day:1,time:"00:00",status:"CONNECTING"};

function App() {
  const [world,setWorld] = useState(fallback);
  useEffect(() => {
    const load = () => fetch("/api/v1/world").then(r=>r.json()).then(setWorld).catch(()=>{});
    load(); const id=setInterval(load,1000); return ()=>clearInterval(id);
  }, []);
  return <div className="app">
    <header className="topbar">
      <div className="brand"><span className="brand-mark">W</span><span>WORLDBUSTER</span></div>
      <div className="world-status"><i/> {world.status} <b>DAY {world.day}</b></div>
    </header>
    <section className="hero">
      <div className="kicker">A WORLD THAT REMEMBERS</div>
      <h1>ENTER THE<br/><em>LIVING WORLD.</em></h1>
      <p>Every action changes something. Every system feeds another system.</p>
      <div className="actions"><button>ENTER WORLD</button><button className="ghost">EXPLORE SYSTEMS</button></div>
    </section>
    <section className="dashboard">
      <article><span>WORLD CLOCK</span><strong>{world.time}</strong><small>TICK {world.tick}</small></article>
      <article><span>WORLD STATE</span><strong>STABLE</strong><small>SERVER AUTHORITATIVE</small></article>
      <article><span>SIMULATION</span><strong>READY</strong><small>EVENT ENGINE ONLINE</small></article>
    </section>
  </div>;
}
createRoot(document.getElementById("root")!).render(<React.StrictMode><App/></React.StrictMode>);
