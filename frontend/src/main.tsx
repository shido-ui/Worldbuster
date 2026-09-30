import React from "react";
import { createRoot } from "react-dom/client";
import "./styles.css";

function App() {
  return <main className="shell">
    <div className="eyebrow">WORLD // ONLINE</div>
    <h1>WORLDBUSTER</h1>
    <p>A living world is initializing.</p>
    <div className="pulse">SERVER ONLINE</div>
  </main>;
}
createRoot(document.getElementById("root")!).render(<React.StrictMode><App/></React.StrictMode>);
