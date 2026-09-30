package com.worldbuster.mobile

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Assignment
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Public
import androidx.compose.material.icons.filled.TrendingUp
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

private enum class Tab(val title:String){ WORLD("World"), MARKET("Market"), MISSIONS("Missions"), PROFILE("Profile") }

class MainActivity:ComponentActivity(){
 override fun onCreate(state:Bundle?){super.onCreate(state);setContent{WorldbusterApp()}}
}

private data class LiveData(
 val world:String="",
 val market:String="",
 val missions:String="",
 val profile:String="",
 val achievements:String="",
 val organizations:String="",
 val news:String="",
 val events:String="",
 val status:String="DISCONNECTED"
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun WorldbusterApp(){
 val context=LocalContext.current
 val api=remember{ApiClient(context,BuildConfig.WORLDBUSTER_BASE_URL)}
 var authenticated by remember{mutableStateOf(api.hasSession())}
 var tab by remember{mutableStateOf(Tab.WORLD)}
 var live by remember{mutableStateOf(LiveData())}
 var loading by remember{mutableStateOf(false)}
 var error by remember{mutableStateOf("")}

 suspend fun load(){
  loading=true
  val result=withContext(Dispatchers.IO){
   val world=api.get("/api/v1/world")
   val profile=api.get("/api/v1/player/state")
   val market=api.get("/api/v1/market/assets")
   val missions=api.get("/api/v1/missions")
   val achievements=api.get("/api/v1/achievements")
   val organizations=api.get("/api/v1/organizations")
   val news=api.get("/api/v1/news")
   val events=api.get("/api/v1/world-events")
   listOf(world,profile,market,missions,achievements,organizations,news,events)
  }
  val unauthorized=result.firstOrNull{it.getOrNull()?.code==401}!=null
  if(unauthorized){api.clearSession();authenticated=false;live=LiveData(status="AUTHENTICATION_REQUIRED");error="Session expired. Sign in again."}
  else if(result.any{it.isFailure||it.getOrNull()?.code?:200>=500}){live=live.copy(status="SERVER_ERROR");error="One or more live services are unavailable."}
  else {
   val values=result.map{it.getOrNull()?.body.orEmpty()}
   live=LiveData(values[0],values[2],values[3],values[1],values[4],values[5],values[6],values[7],"CONNECTED")
   error=""
  }
  loading=false
 }

 LaunchedEffect(authenticated){
  if(authenticated)load()
 }

 MaterialTheme{
  if(!authenticated){
   LoginScreen(api,{authenticated=true},{error=it})
  }else{
   Scaffold(
    topBar={TopAppBar(title={Text("WORLDBUSTER")},actions={Text(live.status,style=MaterialTheme.typography.labelSmall);Spacer(Modifier.width(8.dp))})},
    bottomBar={
     NavigationBar{Tab.values().forEach{t->
      NavigationBarItem(selected=tab==t,onClick={tab=t},icon={Icon(if(t==Tab.WORLD)Icons.Default.Public else if(t==Tab.MARKET)Icons.Default.TrendingUp else if(t==Tab.MISSIONS)Icons.Default.Assignment else Icons.Default.Person,t.title)},label={Text(t.title)})
     }}
    }
   ){pad->
    LazyColumn(Modifier.fillMaxSize().padding(pad).padding(16.dp),verticalArrangement=Arrangement.spacedBy(12.dp)){
     item{Row(Modifier.fillMaxWidth(),horizontalArrangement=Arrangement.spacedBy(8.dp)){Button(onClick={LaunchedEffectKey.launch{load()}}){Text(if(loading)"SYNCING" else "SYNC")};OutlinedButton(onClick={LaunchedEffectKey.logout{api.logout();authenticated=false;live=LiveData(status="DISCONNECTED")}}){Text("LOG OUT")}}}
     if(error.isNotEmpty())item{StatCard("CONNECTION",live.status,error)}
     when(tab){
      Tab.WORLD->{item{StatCard("WORLD","AUTHORITATIVE",live.world)};item{StatCard("NEWS","SERVER FEED",live.news)};item{StatCard("EVENTS","WORLD EVENTS",live.events)}}
      Tab.MARKET->{item{StatCard("MARKET","AUTHORITATIVE ASSETS",live.market)}}
      Tab.MISSIONS->{item{StatCard("MISSIONS","SERVER CONTRACTS",live.missions)}}
      Tab.PROFILE->{item{StatCard("PLAYER STATE","AUTHENTICATED",live.profile)};item{StatCard("ACHIEVEMENTS","SERVER PROGRESS",live.achievements)};item{StatCard("ORGANIZATIONS","FACTIONS",live.organizations)}}
     }
    }
   }
  }
 }
}

private object LaunchedEffectKey {
 var launch:(suspend()->Unit)?=null
 var logout:(()->Unit)?=null
}

@Composable
private fun LoginScreen(api:ApiClient,onSuccess:()->Unit,onError:(String)->Unit){
 var username by remember{mutableStateOf("")}
 var password by remember{mutableStateOf("")}
 var registering by remember{mutableStateOf(false)}
 var busy by remember{mutableStateOf(false)}
 var message by remember{mutableStateOf("")}
 Column(Modifier.fillMaxSize().padding(24.dp),verticalArrangement=Arrangement.spacedBy(12.dp)){
  Text("WORLDBUSTER",style=MaterialTheme.typography.headlineLarge)
  Text(if(registering)"Create your persistent identity." else "Sign in to the persistent world.")
  OutlinedTextField(username,{username=it},label={Text("Username")},singleLine=true)
  OutlinedTextField(password,{password=it},label={Text("Password")},singleLine=true)
  Button(enabled=!busy&&username.length>=3&&password.length>=8,onClick={
   busy=true
   message=""
   // Network work is moved off the UI thread; no credentials are persisted by the client.
   kotlinx.coroutines.GlobalScope.launch(Dispatchers.IO){
    val result=if(registering)api.register(username,password) else api.login(username,password)
    withContext(Dispatchers.Main){
     busy=false
     val r=result.getOrNull()
     if(result.isFailure){message="Network connection unavailable.";return@withContext}
     if(r!!.code !in 200..299){message="Server rejected request: "+r.body.take(300);return@withContext}
     if(registering){registering=false;password="";message="Account created. Sign in to continue."}else{onSuccess()}
    }
   }
  }){Text(if(busy)"CONNECTING..." else if(registering)"CREATE ACCOUNT" else "SIGN IN")}
  TextButton(onClick={registering=!registering;message=""}){Text(if(registering)"BACK TO SIGN IN" else "CREATE ACCOUNT")}
  if(message.isNotEmpty())Text(message)
 }
}

@Composable
private fun StatCard(title:String,value:String,detail:String){
 Card(Modifier.fillMaxWidth()){
  Column(Modifier.padding(18.dp),verticalArrangement=Arrangement.spacedBy(6.dp)){
   Text(title,style=MaterialTheme.typography.labelMedium)
   Text(value,style=MaterialTheme.typography.titleLarge)
   Text(detail.take(6000),style=MaterialTheme.typography.bodySmall)
  }
 }
}
