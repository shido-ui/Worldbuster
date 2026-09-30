import React,{FormEvent,useState} from "react";

type Props={onAuthenticated:(account:{id:string;username:string})=>void};

export function AuthPanel({onAuthenticated}:Props){
 const [mode,setMode]=useState<"login"|"register">("login");
 const [username,setUsername]=useState(""); const [password,setPassword]=useState(""); const [error,setError]=useState("");
 async function submit(e:FormEvent){e.preventDefault();setError("");try{
  const endpoint=mode==="login"?"/api/v1/auth/login":"/api/v1/auth/register";
  const response=await fetch(endpoint,{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({username,password}),credentials:"include"});
  const data=await response.json().catch(()=>({}));
  if(!response.ok){setError(data.error??"Request failed");return}
  if(mode==="register"){setMode("login");setPassword("");setError("Account created. Sign in to enter.");return}
  onAuthenticated(data);
 }catch(error){setError(error instanceof TypeError?"Network connection unavailable.":(error as Error).message)} }
 return <form className="auth-card" onSubmit={submit}>
  <div className="kicker">IDENTITY GATE</div><h2>{mode==="login"?"Return to the world.":"Create your identity."}</h2>
  <input value={username} onChange={e=>setUsername(e.target.value)} placeholder="Username" minLength={3} maxLength={24} autoComplete="username" required/>
  <input value={password} onChange={e=>setPassword(e.target.value)} placeholder="Password" type="password" minLength={8} autoComplete={mode==="login"?"current-password":"new-password"} required/>
  {error&&<small className="auth-error">{error}</small>}
  <button type="submit">{mode==="login"?"ENTER":"CREATE ACCOUNT"}</button>
  <button type="button" className="ghost" onClick={()=>{setMode(mode==="login"?"register":"login");setError("")}}>{mode==="login"?"NEW PLAYER":"BACK TO LOGIN"}</button>
 </form>
}
