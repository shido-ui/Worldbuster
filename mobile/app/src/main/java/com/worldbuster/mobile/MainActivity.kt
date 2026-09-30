package com.worldbuster.mobile

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

private enum class Tab(val title:String){ WORLD("World"), MARKET("Market"), MISSIONS("Missions"), PROFILE("Profile") }

class MainActivity:ComponentActivity(){
 override fun onCreate(state:Bundle?){super.onCreate(state);setContent{WorldbusterApp()}}
}

@Composable
fun WorldbusterApp(){
 var tab by remember{mutableStateOf(Tab.WORLD)}
 MaterialTheme{
  Scaffold(
   topBar={TopAppBar(title={Text("WORLDBUSTER")})},
   bottomBar={
    NavigationBar{
     Tab.values().forEach{t->
      NavigationBarItem(
       selected=tab==t,onClick={tab=t},
       icon={Icon(if(t==Tab.WORLD)Icons.Default.Public else if(t==Tab.MARKET)Icons.Default.TrendingUp else if(t==Tab.MISSIONS)Icons.Default.Assignment else Icons.Default.Person,t.title)},
       label={Text(t.title)})
     }
    }
   }
  ){pad-> Box(Modifier.padding(pad)){Screen(tab)}}
 }
}

@Composable
private fun Screen(tab:Tab){
 LazyColumn(Modifier.fillMaxSize().padding(16.dp),verticalArrangement=Arrangement.spacedBy(12.dp)){
  item{Text(tab.title,style=MaterialTheme.typography.headlineMedium)}
  when(tab){
   Tab.WORLD->{
    item{StatCard("WORLD STATUS","Central District","Live simulation active")}
    item{StatCard("LATEST","Market conditions","World events and news will appear here")}
    item{StatCard("SIMULATION","Population","NPC activity is processed by the server")}
   }
   Tab.MARKET->{
    item{StatCard("MARKET","Tracked assets","Prices are authoritative on the server")}
    items(listOf("Bottled Water","Fabric","Basic Tools")){StatCard(it,"Market data","Awaiting authenticated API connection")}
   }
   Tab.MISSIONS->{
    item{StatCard("MISSIONS","Available contracts","Deterministic mission rules will populate this screen")}
    items(listOf("Training","Exploration","Community")){StatCard(it,"Mission category","Rewards and completion are server-authoritative")}
   }
   Tab.PROFILE->{
    item{StatCard("PLAYER","Guest","Authentication integration is next")}
    item{StatCard("PROGRESSION","Level 1","XP, education, skills and unlocks")}
    item{StatCard("ACHIEVEMENTS","0 unlocked","Achievements and rankings")}
   }
  }
 }
}

@Composable
private fun StatCard(title:String,value:String,detail:String){
 Card(Modifier.fillMaxWidth()){
  Column(Modifier.padding(18.dp),verticalArrangement=Arrangement.spacedBy(6.dp)){
   Text(title,style=MaterialTheme.typography.labelMedium)
   Text(value,style=MaterialTheme.typography.titleLarge)
   Text(detail,style=MaterialTheme.typography.bodyMedium)
  }
 }
}
