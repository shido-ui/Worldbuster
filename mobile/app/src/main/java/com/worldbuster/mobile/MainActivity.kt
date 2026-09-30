package com.worldbuster.mobile
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

class MainActivity:ComponentActivity(){
 override fun onCreate(state:Bundle?){super.onCreate(state);setContent{WorldbusterApp()}}
}
@Composable fun WorldbusterApp(){
 var tab by remember{mutableStateOf("World")}
 MaterialTheme{
  Scaffold(bottomBar={NavigationBar{listOf("World","Market","Missions","Profile").forEach{t->NavigationBarItem(selected=tab==t,onClick={tab=t},icon={},label={Text(t)})}}}){pad->
   Column(Modifier.padding(pad).padding(20.dp)){Text("WORLDBUSTER",style=MaterialTheme.typography.headlineMedium);Spacer(Modifier.height(12.dp));Text(tab,style=MaterialTheme.typography.titleLarge);Spacer(Modifier.height(16.dp));Card(Modifier.fillMaxWidth()){Column(Modifier.padding(16.dp)){Text("The world remembers.");Text("Live simulation, economy, events and progression will appear here.")}}}
  }
 }
}