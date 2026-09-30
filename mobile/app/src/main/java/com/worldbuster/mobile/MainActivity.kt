package com.worldbuster.mobile

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowForward
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.CloudOff
import androidx.compose.material.icons.filled.Lock
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Visibility
import androidx.compose.material.icons.filled.VisibilityOff
import androidx.compose.material.icons.filled.Wifi
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

private enum class Tab(val title: String) { WORLD("World"), MARKET("Market"), MISSIONS("Missions"), PROFILE("Profile") }

class MainActivity : ComponentActivity() {
    override fun onCreate(state: Bundle?) {
        super.onCreate(state)
        setContent { WorldbusterApp() }
    }
}

private data class LiveData(
    val world: String = "", val market: String = "", val missions: String = "", val profile: String = "",
    val achievements: String = "", val organizations: String = "", val news: String = "", val events: String = "",
    val status: String = "DISCONNECTED"
)

private val Bg = Color(0xFF080A10)
private val Surface = Color(0xFF10131C)
private val Surface2 = Color(0xFF151927)
private val Border = Color(0xFF292F42)
private val Accent = Color(0xFF8B6CFF)
private val AccentBright = Color(0xFFA994FF)
private val PrimaryText = Color(0xFFF4F2FF)
private val SecondaryText = Color(0xFFA7A9B8)
private val Success = Color(0xFF61E6A8)
private val Danger = Color(0xFFFF7D91)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun WorldbusterApp() {
    val context = androidx.compose.ui.platform.LocalContext.current
    val scope = rememberCoroutineScope()
    val api = remember { ApiClient(context) }
    var authenticated by remember { mutableStateOf(api.hasSession()) }
    var tab by remember { mutableStateOf(Tab.WORLD) }
    var live by remember { mutableStateOf(LiveData()) }
    var loading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf("") }

    suspend fun load() {
        loading = true
        val result = withContext(Dispatchers.IO) {
            listOf(
                api.get("/api/v1/world"), api.get("/api/v1/player/state"), api.get("/api/v1/market/assets"),
                api.get("/api/v1/missions"), api.get("/api/v1/achievements"), api.get("/api/v1/organizations"),
                api.get("/api/v1/news"), api.get("/api/v1/world-events")
            )
        }
        if (result.any { it.getOrNull()?.code == 401 }) {
            api.clearSession()
            authenticated = false
            live = LiveData(status = "AUTHENTICATION REQUIRED")
            error = "Your session expired. Sign in again."
        } else if (result.any { it.isFailure || (it.getOrNull()?.code ?: 200) >= 500 }) {
            live = live.copy(status = "SERVER ERROR")
            error = api.describeFailure(result.firstOrNull { it.isFailure } ?: Result.failure(Exception("Server error")))
        } else {
            val values = result.map { it.getOrNull()?.body.orEmpty() }
            live = LiveData(values[0], values[2], values[3], values[1], values[4], values[5], values[6], values[7], "CONNECTED")
            error = ""
        }
        loading = false
    }

    LaunchedEffect(authenticated) { if (authenticated) load() }

    MaterialTheme(colorScheme = darkColorScheme(
        background = Bg, surface = Surface, primary = Accent, onPrimary = Color.White,
        onBackground = PrimaryText, onSurface = PrimaryText
    )) {
        if (!authenticated) {
            LoginScreen(api, { authenticated = true }, { error = it })
        } else {
            Scaffold(
                containerColor = Bg,
                topBar = {
                    TopAppBar(
                        title = {
                            Column {
                                Text("WORLDBUSTER", fontWeight = FontWeight.Bold, letterSpacing = 2.sp)
                                Text("PERSISTENT WORLD", style = MaterialTheme.typography.labelSmall, color = SecondaryText)
                            }
                        },
                        actions = { StatusPill(live.status == "CONNECTED", live.status); Spacer(Modifier.width(12.dp)) },
                        colors = TopAppBarDefaults.topAppBarColors(containerColor = Bg)
                    )
                },
                bottomBar = {
                    NavigationBar(containerColor = Surface) {
                        Tab.values().forEach { t ->
                            NavigationBarItem(
                                selected = tab == t, onClick = { tab = t },
                                icon = {
                                    Icon(
                                        when (t) {
                                            Tab.WORLD -> Icons.Default.Wifi
                                            Tab.MARKET -> Icons.Default.ArrowForward
                                            Tab.MISSIONS -> Icons.Default.CheckCircle
                                            Tab.PROFILE -> Icons.Default.Person
                                        }, t.title
                                    )
                                },
                                label = { Text(t.title) }
                            )
                        }
                    }
                }
            ) { pad ->
                Column(Modifier.fillMaxSize().padding(pad).padding(horizontal = 16.dp)) {
                    Spacer(Modifier.height(12.dp))
                    Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(10.dp)) {
                        Button(Modifier.weight(1f), enabled = !loading, onClick = { scope.launch { load() } }) {
                            Text(if (loading) "SYNCING…" else "SYNC WORLD")
                        }
                        OutlinedButton(onClick = {
                            scope.launch {
                                withContext(Dispatchers.IO) { api.logout() }
                                authenticated = false
                                live = LiveData()
                            }
                        }) { Text("LOG OUT") }
                    }
                    if (error.isNotEmpty()) {
                        Spacer(Modifier.height(12.dp))
                        ConnectionBanner(error)
                    }
                    LazyColumn(
                        Modifier.fillMaxSize(),
                        verticalArrangement = Arrangement.spacedBy(12.dp),
                        contentPadding = PaddingValues(top = 12.dp, bottom = 24.dp)
                    ) {
                        when (tab) {
                            Tab.WORLD -> {
                                item { StatCard("WORLD", "AUTHORITATIVE", live.world) }
                                item { StatCard("NEWS", "SERVER FEED", live.news) }
                                item { StatCard("EVENTS", "WORLD EVENTS", live.events) }
                            }
                            Tab.MARKET -> item { StatCard("MARKET", "AUTHORITATIVE ASSETS", live.market) }
                            Tab.MISSIONS -> item { StatCard("MISSIONS", "SERVER CONTRACTS", live.missions) }
                            Tab.PROFILE -> {
                                item { StatCard("PLAYER STATE", "AUTHENTICATED", live.profile) }
                                item { StatCard("ACHIEVEMENTS", "SERVER PROGRESS", live.achievements) }
                                item { StatCard("ORGANIZATIONS", "FACTIONS", live.organizations) }
                            }
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun LoginScreen(api: ApiClient, onSuccess: () -> Unit, onError: (String) -> Unit) {
    val scope = rememberCoroutineScope()
    var username by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var serverUrl by remember { mutableStateOf(api.baseUrl()) }
    var registering by remember { mutableStateOf(false) }
    var busy by remember { mutableStateOf(false) }
    var checking by remember { mutableStateOf(false) }
    var showPassword by remember { mutableStateOf(false) }
    var connectionMessage by remember { mutableStateOf("") }
    var connected by remember { mutableStateOf(false) }
    var showServerSettings by remember { mutableStateOf(api.isPhysicalDevice()) }

    LaunchedEffect(Unit) {
        if (serverUrl.isNotBlank()) {
            checking = true
            val result = withContext(Dispatchers.IO) { api.health() }
            checking = false
            connected = result.isSuccess && result.getOrNull()?.code in 200..299
            connectionMessage = if (connected) "Server reachable" else api.describeFailure(result)
        }
    }

    Box(Modifier.fillMaxSize().background(Brush.verticalGradient(listOf(Color(0xFF17132A), Bg, Bg)))) {
        Column(
            Modifier.fillMaxSize().padding(horizontal = 24.dp).padding(top = 56.dp, bottom = 28.dp),
            verticalArrangement = Arrangement.SpaceBetween
        ) {
            Column(Modifier.fillMaxWidth()) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Box(
                        Modifier.size(44.dp).clip(RoundedCornerShape(14.dp))
                            .background(Brush.linearGradient(listOf(Accent, Color(0xFF5941B8)))),
                        contentAlignment = Alignment.Center
                    ) { Text("W", color = Color.White, fontSize = 22.sp, fontWeight = FontWeight.Black) }
                    Spacer(Modifier.width(12.dp))
                    Column {
                        Text("WORLDBUSTER", color = PrimaryText, fontSize = 17.sp, fontWeight = FontWeight.Bold, letterSpacing = 2.5.sp)
                        Text("PERSISTENT WORLD", color = SecondaryText, fontSize = 10.sp, letterSpacing = 1.7.sp)
                    }
                }

                Spacer(Modifier.height(52.dp))
                Text(if (registering) "Create your identity." else "Return to the world.",
                    color = PrimaryText, fontSize = 34.sp, lineHeight = 39.sp, fontWeight = FontWeight.Bold)
                Spacer(Modifier.height(10.dp))
                Text(
                    if (registering) "Your character lives in a world that keeps moving."
                    else "Sign in to continue your persistent-world session.",
                    color = SecondaryText, fontSize = 15.sp, lineHeight = 22.sp
                )

                Spacer(Modifier.height(28.dp))
                AuthField(username, { username = it }, "USERNAME", "Enter username",
                    { Icon(Icons.Default.Person, null) }, keyboardType = KeyboardType.Text)
                Spacer(Modifier.height(12.dp))
                AuthField(password, { password = it }, "PASSWORD", "Enter password",
                    { Icon(Icons.Default.Lock, null) }, {
                        IconButton(onClick = { showPassword = !showPassword }) {
                            Icon(if (showPassword) Icons.Default.VisibilityOff else Icons.Default.Visibility, "Toggle password visibility")
                        }
                    }, KeyboardType.Password,
                    if (showPassword) VisualTransformation.None else PasswordVisualTransformation())

                Spacer(Modifier.height(16.dp))
                Surface(
                    Modifier.fillMaxWidth(), shape = RoundedCornerShape(16.dp), color = Surface2,
                    border = BorderStroke(1.dp, Border)
                ) {
                    Row(Modifier.padding(14.dp), verticalAlignment = Alignment.CenterVertically) {
                        Icon(if (connected) Icons.Default.Wifi else Icons.Default.CloudOff, null,
                            tint = if (connected) Success else SecondaryText)
                        Spacer(Modifier.width(10.dp))
                        Column(Modifier.weight(1f)) {
                            Text(
                                if (checking) "CHECKING SERVER…" else if (connected) "SERVER ONLINE" else "SERVER NOT CONNECTED",
                                fontSize = 12.sp, fontWeight = FontWeight.Bold,
                                color = if (connected) Success else PrimaryText
                            )
                            Text(connectionMessage.ifBlank { "Set the server address below." },
                                fontSize = 11.sp, color = SecondaryText, maxLines = 2)
                        }
                        TextButton(onClick = { showServerSettings = !showServerSettings }) {
                            Text(if (showServerSettings) "HIDE" else "SERVER")
                        }
                    }
                }

                if (showServerSettings) {
                    Spacer(Modifier.height(10.dp))
                    AuthField(
                        serverUrl, {
                            serverUrl = it
                            connected = false
                            connectionMessage = ""
                        },
                        "SERVER ADDRESS",
                        if (api.isPhysicalDevice()) "http://192.168.x.x:8080" else "http://10.0.2.2:8080",
                        keyboardType = KeyboardType.Uri
                    )
                    Spacer(Modifier.height(8.dp))
                    Text(
                        if (api.isPhysicalDevice())
                            "On a physical phone, use the computer's LAN address running Worldbuster."
                        else
                            "Emulator uses 10.0.2.2 to reach the development machine.",
                        color = SecondaryText, fontSize = 11.sp, lineHeight = 16.sp
                    )
                    TextButton(
                        enabled = !checking && serverUrl.isNotBlank(),
                        onClick = {
                            api.setBaseUrl(serverUrl)
                            scope.launch {
                                checking = true
                                val result = withContext(Dispatchers.IO) { api.health() }
                                checking = false
                                connected = result.isSuccess && result.getOrNull()?.code in 200..299
                                connectionMessage = if (connected) "Server reachable" else api.describeFailure(result)
                            }
                        }
                    ) { Text("TEST CONNECTION") }
                }
            }

            Column {
                Button(
                    Modifier.fillMaxWidth().height(56.dp), enabled = !busy && username.length >= 3 &&
                        password.length >= 8 && serverUrl.isNotBlank(), shape = RoundedCornerShape(17.dp),
                    onClick = {
                        api.setBaseUrl(serverUrl)
                        busy = true
                        onError("")
                        scope.launch {
                            val result = withContext(Dispatchers.IO) {
                                if (registering) api.register(username, password) else api.login(username, password)
                            }
                            busy = false
                            val response = result.getOrNull()
                            if (result.isFailure) {
                                connected = false
                                connectionMessage = api.describeFailure(result)
                                onError(connectionMessage)
                                return@launch
                            }
                            if (response!!.code !in 200..299) {
                                onError(api.describeHttpError(response))
                                return@launch
                            }
                            connected = true
                            connectionMessage = "Server reachable"
                            if (registering) {
                                registering = false
                                password = ""
                                onError("Account created. Sign in to continue.")
                            } else onSuccess()
                        }
                    }
                ) {
                    Text(if (busy) "CONNECTING…" else if (registering) "CREATE ACCOUNT" else "ENTER WORLD",
                        fontWeight = FontWeight.Bold, letterSpacing = 1.sp)
                    if (!busy) { Spacer(Modifier.width(8.dp)); Icon(Icons.Default.ArrowForward, null) }
                }
                Spacer(Modifier.height(8.dp))
                TextButton(Modifier.fillMaxWidth(), onClick = {
                    registering = !registering
                    onError("")
                }) {
                    Text(if (registering) "BACK TO SIGN IN" else "CREATE NEW IDENTITY",
                        color = AccentBright, fontWeight = FontWeight.SemiBold)
                }
            }
        }
    }
}

@Composable
private fun AuthField(
    value: String, onValueChange: (String) -> Unit, label: String, placeholder: String,
    leadingIcon: @Composable (() -> Unit)? = null, trailingIcon: @Composable (() -> Unit)? = null,
    keyboardType: KeyboardType = KeyboardType.Text,
    visualTransformation: VisualTransformation = VisualTransformation.None
) {
    OutlinedTextField(
        Modifier.fillMaxWidth(), value, onValueChange,
        label = { Text(label, letterSpacing = 1.1.sp, fontSize = 10.sp) },
        placeholder = { Text(placeholder, color = SecondaryText) },
        leadingIcon = leadingIcon, trailingIcon = trailingIcon, singleLine = true,
        keyboardOptions = KeyboardOptions(keyboardType = keyboardType),
        visualTransformation = visualTransformation, shape = RoundedCornerShape(16.dp),
        colors = OutlinedTextFieldDefaults.colors(
            focusedBorderColor = Accent, unfocusedBorderColor = Border,
            focusedLabelColor = AccentBright, unfocusedLabelColor = SecondaryText,
            focusedLeadingIconColor = AccentBright, unfocusedLeadingIconColor = SecondaryText,
            focusedTextColor = PrimaryText, unfocusedTextColor = PrimaryText, cursorColor = AccentBright
        )
    )
}

@Composable
private fun StatusPill(connected: Boolean, label: String) {
    Surface(
        shape = RoundedCornerShape(50), color = if (connected) Color(0x142EEB86) else Color(0x14FFFFFF),
        border = BorderStroke(1.dp, if (connected) Color(0x405EEFA3) else Border)
    ) {
        Row(Modifier.padding(horizontal = 10.dp, vertical = 6.dp), verticalAlignment = Alignment.CenterVertically) {
            Box(Modifier.size(6.dp).clip(RoundedCornerShape(50)).background(if (connected) Success else SecondaryText))
            Spacer(Modifier.width(6.dp))
            Text(label, fontSize = 9.sp, letterSpacing = .8.sp, color = if (connected) Success else SecondaryText)
        }
    }
}

@Composable
private fun ConnectionBanner(message: String) {
    Surface(
        Modifier.fillMaxWidth(), shape = RoundedCornerShape(14.dp), color = Color(0x14FF7D91),
        border = BorderStroke(1.dp, Color(0x35FF7D91))
    ) {
        Row(Modifier.padding(14.dp), verticalAlignment = Alignment.CenterVertically) {
            Icon(Icons.Default.CloudOff, null, tint = Danger)
            Spacer(Modifier.width(10.dp))
            Text(message, color = PrimaryText, fontSize = 12.sp, lineHeight = 17.sp)
        }
    }
}

@Composable
private fun StatCard(title: String, value: String, detail: String) {
    Card(
        Modifier.fillMaxWidth(), shape = RoundedCornerShape(18.dp),
        colors = CardDefaults.cardColors(containerColor = Surface),
        border = BorderStroke(1.dp, Border)
    ) {
        Column(Modifier.padding(18.dp), verticalArrangement = Arrangement.spacedBy(6.dp)) {
            Text(title, fontSize = 10.sp, letterSpacing = 1.5.sp, color = SecondaryText)
            Text(value, fontSize = 20.sp, fontWeight = FontWeight.Bold)
            Text(detail.take(6000), fontSize = 12.sp, lineHeight = 17.sp, color = SecondaryText)
        }
    }
}
