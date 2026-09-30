package com.worldbuster.mobile

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowForward
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.CloudOff
import androidx.compose.material.icons.filled.Inventory
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.material.icons.filled.TrendingDown
import androidx.compose.material.icons.filled.TrendingUp
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
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import java.text.NumberFormat
import java.util.Locale

private enum class Tab(val title: String) { WORLD("World"), MARKET("Market"), MISSIONS("Missions"), PROFILE("Profile") }
private enum class SyncState(val label: String) { CONNECTED("LIVE"), DEGRADED("SYNCED"), OFFLINE("OFFLINE"), AUTH("SIGN IN") }

class MainActivity : ComponentActivity() {
    override fun onCreate(state: Bundle?) {
        super.onCreate(state)
        setContent { WorldbusterApp() }
    }
}

private data class Endpoint(val key: String, val path: String)

private val Bg = Color(0xFF070910)
private val Surface = Color(0xFF0E111A)
private val SurfaceRaised = Color(0xFF141926)
private val SurfaceSoft = Color(0xFF191E2C)
private val Border = Color(0xFF252C3D)
private val Accent = Color(0xFF8B6CFF)
private val AccentBright = Color(0xFFB7A6FF)
private val PrimaryText = Color(0xFFF5F3FF)
private val SecondaryText = Color(0xFF969BAE)
private val Success = Color(0xFF62E6A8)
private val Warning = Color(0xFFFFC46B)
private val Danger = Color(0xFFFF7188)
private val Blue = Color(0xFF68B9FF)

@Composable
fun WorldbusterApp(viewModel: WorldbusterViewModel = viewModel()) {
    val state by viewModel.uiState.collectAsStateWithLifecycle()
    var tabName by remember { mutableStateOf(Tab.WORLD.name) }
    val tab = runCatching { Tab.valueOf(tabName) }.getOrDefault(Tab.WORLD)

    MaterialTheme(
        colorScheme = darkColorScheme(
            background = Bg, surface = Surface, surfaceVariant = SurfaceRaised,
            primary = Accent, onPrimary = Color.White, onBackground = PrimaryText,
            onSurface = PrimaryText, onSurfaceVariant = SecondaryText, outline = Border
        )
    ) {
        if (!state.authenticated) {
            LoginScreen(
                viewModel = viewModel,
                error = state.error
            )
        } else {
            Scaffold(
                containerColor = Bg,
                topBar = { TopBar(state.syncState, state.dashboard.world) },
                bottomBar = {
                    NavigationBar(containerColor = Surface, tonalElevation = 0.dp) {
                        Tab.values().forEach { item ->
                            NavigationBarItem(
                                selected = tab == item,
                                onClick = { tabName = item.name },
                                icon = {
                                    Icon(
                                        when (item) {
                                            Tab.WORLD -> Icons.Default.Wifi
                                            Tab.MARKET -> Icons.Default.TrendingUp
                                            Tab.MISSIONS -> Icons.Default.CheckCircle
                                            Tab.PROFILE -> Icons.Default.Person
                                        }, item.title
                                    )
                                },
                                label = { Text(item.title) }
                            )
                        }
                    }
                }
            ) { pad ->
                Column(Modifier.fillMaxSize().padding(pad)) {
                    AnimatedVisibility(state.error.isNotBlank()) {
                        ConnectionBanner(state.error, state.syncState)
                    }
                    LazyColumn(
                        Modifier.fillMaxSize(),
                        verticalArrangement = Arrangement.spacedBy(12.dp),
                        contentPadding = PaddingValues(start = 16.dp, end = 16.dp, top = 14.dp, bottom = 28.dp)
                    ) {
                        item {
                            SyncToolbar(
                                state.syncState,
                                state.loading,
                                state.dashboard.lastSync,
                                viewModel::refresh,
                                viewModel::logout
                            )
                        }
                        when (tab) {
                            Tab.WORLD -> {
                                item { WorldHero(state.dashboard.world, state.dashboard.player) }
                                item { MarketPulse(state.dashboard.market.take(3)) }
                                item { SectionTitle("WORLD FEED", "Live consequences from the server") }
                                items(state.dashboard.news.take(4)) { NewsCard(it) }
                                items(state.dashboard.events.take(3)) { EventCard(it) }
                                if (state.dashboard.news.isEmpty() && state.dashboard.events.isEmpty()) {
                                    item { EmptyCard("The world is quiet.", "New events and news will appear here as the simulation produces them.") }
                                }
                            }
                            Tab.MARKET -> {
                                item { SectionTitle("MARKET", "Authoritative live assets") }
                                if (state.dashboard.market.isEmpty()) {
                                    item { EmptyCard("No market assets available.", "The server has not published any assets yet.") }
                                } else {
                                    items(state.dashboard.market) { MarketCard(it) }
                                }
                            }
                            Tab.MISSIONS -> {
                                item { SectionTitle("MISSIONS", "Contracts issued by the world") }
                                if (state.dashboard.missions.isEmpty()) {
                                    item { EmptyCard("No missions available.", "New contracts will appear as the world unlocks them.") }
                                } else {
                                    items(state.dashboard.missions) { MissionCard(it) }
                                }
                            }
                            Tab.PROFILE -> {
                                item { ProfileHero(state.dashboard.player) }
                                item { InventoryCard(state.dashboard.player) }
                                item { SectionTitle("ACHIEVEMENTS", "Progress recorded by the server") }
                                items(state.dashboard.achievements.take(6)) { AchievementCard(it) }
                                item { SectionTitle("ORGANIZATIONS", "Groups shaping the world") }
                                items(state.dashboard.organizations.take(5)) { OrganizationCard(it) }
                            }
                        }
                    }
                }
            }
        }
    }
}

private fun applyResponse(data: DashboardView, key: String, body: String): DashboardView = runCatching {
    when (key) {
        "world" -> data.copy(world = parseWorld(body))
        "dashboard" -> data.copy(player = parsePlayerDashboard(body))
        "market" -> data.copy(market = parseMarket(body))
        "missions" -> data.copy(missions = parseMissions(body))
        "achievements" -> data.copy(achievements = parseAchievements(body))
        "organizations" -> data.copy(organizations = parseOrganizations(body))
        "news" -> data.copy(news = parseNews(body))
        "events" -> data.copy(events = parseEvents(body))
        else -> data
    }
}.getOrElse { data }

private fun nowLabel(): String = SimpleDateFormat("HH:mm:ss", Locale.getDefault()).format(Date())
private fun money(value: Long): String = NumberFormat.getIntegerInstance(Locale.getDefault()).format(value) + " WBX"

@Composable
private fun TopBar(state: SyncState, world: WorldView) {
    Surface(color = Bg, tonalElevation = 0.dp) {
        Row(Modifier.fillMaxWidth().padding(horizontal = 16.dp, vertical = 12.dp), verticalAlignment = Alignment.CenterVertically) {
            WorldMark(36)
            Spacer(Modifier.width(10.dp))
            Column(Modifier.weight(1f)) {
                Text("WORLDBUSTER", fontSize = 14.sp, fontWeight = FontWeight.Black, letterSpacing = 2.sp)
                Text("DAY ${world.day}  •  ${world.time}", fontSize = 10.sp, color = SecondaryText, letterSpacing = 1.sp)
            }
            StatusPill(state)
        }
    }
}

@Composable
private fun WorldMark(size: Int) {
    Box(
        Modifier.size(size.dp).clip(RoundedCornerShape((size / 3).dp))
            .background(Brush.linearGradient(listOf(Accent, Color(0xFF4B36A6)))),
        contentAlignment = Alignment.Center
    ) {
        Text("W", color = Color.White, fontSize = (size * .48f).sp, fontWeight = FontWeight.Black)
    }
}

@Composable
private fun StatusPill(state: SyncState) {
    val color by animateColorAsState(
        when (state) {
            SyncState.CONNECTED -> Success
            SyncState.DEGRADED -> Warning
            SyncState.OFFLINE -> Danger
            SyncState.AUTH -> SecondaryText
        }, label = "status"
    )
    Surface(shape = RoundedCornerShape(50), color = color.copy(alpha = .09f), border = BorderStroke(1.dp, color.copy(alpha = .22f))) {
        Row(Modifier.padding(horizontal = 9.dp, vertical = 6.dp), verticalAlignment = Alignment.CenterVertically) {
            Box(Modifier.size(6.dp).clip(RoundedCornerShape(50)).background(color))
            Spacer(Modifier.width(6.dp))
            Text(state.label, fontSize = 9.sp, color = color, letterSpacing = .8.sp, fontWeight = FontWeight.Bold)
        }
    }
}

@Composable
private fun SyncToolbar(state: SyncState, loading: Boolean, lastSync: String, onSync: () -> Unit, onLogout: () -> Unit) {
    Row(Modifier.fillMaxWidth(), verticalAlignment = Alignment.CenterVertically) {
        Column(Modifier.weight(1f)) {
            Text("WORLD SYNC", fontSize = 9.sp, letterSpacing = 1.4.sp, color = SecondaryText)
            Text(if (lastSync.isBlank()) "Waiting for first sync" else "Last synced ${lastSync}", fontSize = 12.sp, color = PrimaryText)
        }
        IconButton(onClick = onSync, enabled = !loading) { Icon(Icons.Default.Refresh, "Sync") }
        TextButton(onClick = onLogout) { Text("EXIT", color = SecondaryText) }
    }
}

@Composable
private fun WorldHero(world: WorldView, player: PlayerView) {
    Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(26.dp), color = SurfaceRaised, border = BorderStroke(1.dp, Border)) {
        Box(Modifier.background(Brush.linearGradient(listOf(Color(0xFF17152A), SurfaceRaised)))) {
            Column(Modifier.padding(22.dp)) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Column(Modifier.weight(1f)) {
                        Text("THE WORLD IS MOVING", fontSize = 10.sp, color = AccentBright, letterSpacing = 1.8.sp, fontWeight = FontWeight.Bold)
                        Text("Welcome back, ${player.name}.", fontSize = 25.sp, fontWeight = FontWeight.Bold)
                        Text("Everything here is server-authoritative.", fontSize = 12.sp, color = SecondaryText)
                    }
                    WorldMark(48)
                }
                Spacer(Modifier.height(22.dp))
                Row(horizontalArrangement = Arrangement.spacedBy(10.dp)) {
                    Metric("TICK", world.tick.toString(), Modifier.weight(1f))
                    Metric("ONLINE", world.online.toString(), Modifier.weight(1f))
                    Metric("LEVEL", player.level.toString(), Modifier.weight(1f))
                }
            }
        }
    }
}

@Composable
private fun ProfileHero(player: PlayerView) {
    Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(26.dp), color = SurfaceRaised, border = BorderStroke(1.dp, Border)) {
        Column(Modifier.padding(22.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                WorldMark(52)
                Spacer(Modifier.width(14.dp))
                Column(Modifier.weight(1f)) {
                    Text(player.name, fontSize = 24.sp, fontWeight = FontWeight.Bold)
                    Text("LEVEL ${player.level}", fontSize = 10.sp, color = AccentBright, letterSpacing = 1.5.sp, fontWeight = FontWeight.Bold)
                }
                Column(horizontalAlignment = Alignment.End) {
                    Text(money(player.cash), fontSize = 13.sp, fontWeight = FontWeight.Bold)
                    Text("BALANCE", fontSize = 9.sp, color = SecondaryText, letterSpacing = 1.sp)
                }
            }
            Spacer(Modifier.height(20.dp))
            StatBar("ENERGY", player.energy / 100f, "${player.energy}%")
            Spacer(Modifier.height(10.dp))
            StatBar("XP", ((player.xp % 1000L) / 1000f).toFloat().coerceIn(0f, 1f), player.xp.toString())
        }
    }
}

@Composable
private fun Metric(label: String, value: String, modifier: Modifier = Modifier) {
    Surface(modifier, shape = RoundedCornerShape(15.dp), color = SurfaceSoft) {
        Column(Modifier.padding(12.dp)) {
            Text(label, fontSize = 8.sp, color = SecondaryText, letterSpacing = 1.2.sp)
            Text(value, fontSize = 17.sp, fontWeight = FontWeight.Bold)
        }
    }
}

@Composable
private fun StatBar(label: String, value: Float, detail: String) {
    Row(verticalAlignment = Alignment.CenterVertically) {
        Text(label, Modifier.width(54.dp), fontSize = 8.sp, color = SecondaryText, letterSpacing = .8.sp)
        LinearProgressIndicator(progress = value.coerceIn(0f, 1f), Modifier.weight(1f).height(6.dp).clip(RoundedCornerShape(50)), color = Accent, trackColor = Border)
        Spacer(Modifier.width(10.dp))
        Text(detail, fontSize = 9.sp, color = PrimaryText)
    }
}

@Composable
private fun MarketPulse(assets: List<MarketAssetView>) {
    if (assets.isEmpty()) return
    Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
        SectionTitle("MARKET PULSE", "Supply and demand pressure")
        assets.forEach { asset ->
            Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(16.dp), color = Surface, border = BorderStroke(1.dp, Border)) {
                Row(Modifier.padding(13.dp), verticalAlignment = Alignment.CenterVertically) {
                    Column(Modifier.weight(1f)) {
                        Text(asset.name, fontSize = 13.sp, fontWeight = FontWeight.SemiBold)
                        Text(asset.symbol.uppercase(), fontSize = 8.sp, color = SecondaryText, letterSpacing = 1.sp)
                    }
                    Text(money(asset.price), fontSize = 12.sp, fontWeight = FontWeight.Bold)
                    Spacer(Modifier.width(10.dp))
                    Icon(if (asset.pressure >= .5f) Icons.Default.TrendingUp else Icons.Default.TrendingDown, null, tint = if (asset.pressure >= .5f) Warning else Blue, modifier = Modifier.size(18.dp))
                }
            }
        }
    }
}

@Composable
private fun MarketCard(asset: MarketAssetView) {
    val pressure by animateFloatAsState(asset.pressure.coerceIn(0f, 1f), label = asset.symbol)
    val hot = pressure >= .6f
    Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(20.dp), color = Surface, border = BorderStroke(1.dp, if (hot) Color(0x445E4FAE) else Border)) {
        Column(Modifier.padding(18.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Box(Modifier.size(42.dp).clip(RoundedCornerShape(13.dp)).background(if (hot) Color(0x202F254E) else SurfaceSoft), contentAlignment = Alignment.Center) {
                    Icon(if (hot) Icons.Default.TrendingUp else Icons.Default.Inventory, null, tint = if (hot) AccentBright else SecondaryText)
                }
                Spacer(Modifier.width(12.dp))
                Column(Modifier.weight(1f)) {
                    Text(asset.name, fontSize = 15.sp, fontWeight = FontWeight.Bold)
                    Text(asset.category.uppercase() + "  •  " + asset.symbol.uppercase(), fontSize = 8.sp, color = SecondaryText, letterSpacing = .8.sp)
                }
                Text(money(asset.price), fontSize = 14.sp, fontWeight = FontWeight.Black)
            }
            Spacer(Modifier.height(16.dp))
            Row {
                MarketStat("SUPPLY", asset.supply)
                MarketStat("DEMAND", asset.demand)
            }
            Spacer(Modifier.height(10.dp))
            LinearProgressIndicator(progress = pressure, Modifier.fillMaxWidth().height(5.dp).clip(RoundedCornerShape(50)), color = if (hot) Warning else Blue, trackColor = Border)
        }
    }
}

@Composable
private fun MarketStat(label: String, value: Long) {
    Column(Modifier.weight(1f)) {
        Text(label, fontSize = 8.sp, color = SecondaryText, letterSpacing = 1.sp)
        Text(value.toString(), fontSize = 12.sp, fontWeight = FontWeight.SemiBold)
    }
}

@Composable
private fun MissionCard(mission: MissionView) {
    Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(20.dp), color = Surface, border = BorderStroke(1.dp, Border)) {
        Column(Modifier.padding(18.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Column(Modifier.weight(1f)) {
                    Text(mission.title, fontSize = 16.sp, fontWeight = FontWeight.Bold)
                    Text(mission.type, fontSize = 8.sp, color = AccentBright, letterSpacing = 1.1.sp)
                }
                Icon(Icons.Default.ArrowForward, null, tint = SecondaryText)
            }
            Spacer(Modifier.height(8.dp))
            Text(mission.description, fontSize = 12.sp, lineHeight = 18.sp, color = SecondaryText)
            Spacer(Modifier.height(14.dp))
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                Chip(money(mission.rewardCash) + " reward")
                Chip(mission.rewardXp.toString() + " XP")
                Chip("Target " + mission.target)
            }
        }
    }
}

@Composable
private fun AchievementCard(item: AchievementView) {
    val progress = if (item.target <= 0) 0f else (item.progress.toFloat() / item.target).coerceIn(0f, 1f)
    Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(18.dp), color = Surface, border = BorderStroke(1.dp, Border)) {
        Column(Modifier.padding(16.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Icon(Icons.Default.CheckCircle, null, tint = if (progress >= 1f) Success else SecondaryText)
                Spacer(Modifier.width(10.dp))
                Column(Modifier.weight(1f)) {
                    Text(item.title, fontSize = 14.sp, fontWeight = FontWeight.SemiBold)
                    Text(item.description, fontSize = 11.sp, color = SecondaryText)
                }
                Text(item.points.toString() + " pts", fontSize = 10.sp, color = AccentBright)
            }
            Spacer(Modifier.height(10.dp))
            StatBar("PROGRESS", progress, item.progress.toString() + "/" + item.target)
        }
    }
}

@Composable
private fun OrganizationCard(item: OrganizationView) {
    Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(18.dp), color = Surface, border = BorderStroke(1.dp, Border)) {
        Row(Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
            Box(Modifier.size(40.dp).clip(RoundedCornerShape(12.dp)).background(SurfaceSoft), contentAlignment = Alignment.Center) {
                Text(item.name.take(1).uppercase(), fontWeight = FontWeight.Black, color = AccentBright)
            }
            Spacer(Modifier.width(12.dp))
            Column(Modifier.weight(1f)) {
                Text(item.name, fontSize = 14.sp, fontWeight = FontWeight.SemiBold)
                Text(item.type + "  •  LEVEL " + item.level, fontSize = 9.sp, color = SecondaryText)
            }
            Text("REP " + item.reputation, fontSize = 9.sp, color = AccentBright)
        }
    }
}

@Composable
private fun InventoryCard(player: PlayerView) {
    Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(20.dp), color = Surface, border = BorderStroke(1.dp, Border)) {
        Row(Modifier.padding(18.dp), verticalAlignment = Alignment.CenterVertically) {
            Icon(Icons.Default.Inventory, null, tint = AccentBright)
            Spacer(Modifier.width(12.dp))
            Column(Modifier.weight(1f)) {
                Text("INVENTORY", fontSize = 10.sp, color = SecondaryText, letterSpacing = 1.2.sp)
                Text(player.inventoryUsed.toString() + " / " + player.inventoryCapacity + " occupied", fontSize = 14.sp, fontWeight = FontWeight.SemiBold)
            }
            Text("READY", fontSize = 9.sp, color = Success, fontWeight = FontWeight.Bold)
        }
    }
}

@Composable
private fun NewsCard(item: NewsView) {
    Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(18.dp), color = Surface, border = BorderStroke(1.dp, Border)) {
        Column(Modifier.padding(16.dp)) {
            Text(item.category.uppercase(), fontSize = 8.sp, color = AccentBright, letterSpacing = 1.2.sp)
            Spacer(Modifier.height(5.dp))
            Text(item.headline, fontSize = 15.sp, fontWeight = FontWeight.Bold)
            if (item.body.isNotBlank()) {
                Spacer(Modifier.height(4.dp))
                Text(item.body, fontSize = 11.sp, color = SecondaryText, lineHeight = 17.sp)
            }
        }
    }
}

@Composable
private fun EventCard(item: EventView) {
    Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(18.dp), color = Surface, border = BorderStroke(1.dp, Border)) {
        Row(Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
            Box(Modifier.size(8.dp).clip(RoundedCornerShape(50)).background(if (item.severity >= 4) Danger else Warning))
            Spacer(Modifier.width(12.dp))
            Column(Modifier.weight(1f)) {
                Text(item.code.replace('-', ' ').uppercase(), fontSize = 13.sp, fontWeight = FontWeight.SemiBold)
                Text(item.status + "  •  " + item.location, fontSize = 9.sp, color = SecondaryText)
            }
            Text("S" + item.severity, fontSize = 9.sp, color = if (item.severity >= 4) Danger else Warning)
        }
    }
}

@Composable
private fun SectionTitle(title: String, subtitle: String) {
    Column(Modifier.padding(top = 4.dp, bottom = 2.dp)) {
        Text(title, fontSize = 10.sp, color = AccentBright, letterSpacing = 1.5.sp, fontWeight = FontWeight.Bold)
        Text(subtitle, fontSize = 12.sp, color = SecondaryText)
    }
}

@Composable
private fun Chip(text: String) {
    Surface(shape = RoundedCornerShape(50), color = SurfaceSoft) {
        Text(text, Modifier.padding(horizontal = 9.dp, vertical = 6.dp), fontSize = 9.sp, color = PrimaryText)
    }
}

@Composable
private fun EmptyCard(title: String, body: String) {
    Surface(Modifier.fillMaxWidth(), shape = RoundedCornerShape(20.dp), color = Surface, border = BorderStroke(1.dp, Border)) {
        Column(Modifier.padding(22.dp)) {
            Text(title, fontSize = 16.sp, fontWeight = FontWeight.Bold)
            Spacer(Modifier.height(6.dp))
            Text(body, fontSize = 12.sp, lineHeight = 18.sp, color = SecondaryText)
        }
    }
}

@Composable
private fun ConnectionBanner(message: String, state: SyncState) {
    val color = if (state == SyncState.DEGRADED) Warning else Danger
    Surface(Modifier.fillMaxWidth(), color = color.copy(alpha = .08f), border = BorderStroke(1.dp, color.copy(alpha = .2f))) {
        Row(Modifier.padding(horizontal = 16.dp, vertical = 11.dp), verticalAlignment = Alignment.CenterVertically) {
            Icon(if (state == SyncState.DEGRADED) Icons.Default.Wifi else Icons.Default.CloudOff, null, tint = color)
            Spacer(Modifier.width(10.dp))
            Text(message, Modifier.weight(1f), fontSize = 11.sp, color = PrimaryText, lineHeight = 16.sp)
        }
    }
}

@Composable
private fun LoginScreen(
    viewModel: WorldbusterViewModel,
    error: String
) {
    var username by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var serverUrl by remember(viewModel.serverUrl) { mutableStateOf(viewModel.serverUrl) }
    var registering by remember { mutableStateOf(false) }
    var showPassword by remember { mutableStateOf(false) }
    var showServerSettings by remember(viewModel.physicalDevice) { mutableStateOf(viewModel.physicalDevice) }
    val state by viewModel.uiState.collectAsStateWithLifecycle()

    LaunchedEffect(Unit) {
        if (serverUrl.isNotBlank()) viewModel.testConnection(serverUrl)
    }

    Box(
        Modifier.fillMaxSize().background(
            Brush.verticalGradient(listOf(Color(0xFF17132A), Bg, Bg))
        )
    ) {
        Column(
            Modifier.fillMaxSize().padding(horizontal = 24.dp).padding(top = 52.dp, bottom = 26.dp),
            verticalArrangement = Arrangement.SpaceBetween
        ) {
            Column(Modifier.fillMaxWidth()) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    WorldMark(46)
                    Spacer(Modifier.width(12.dp))
                    Column {
                        Text("WORLDBUSTER", color = PrimaryText, fontSize = 17.sp, fontWeight = FontWeight.Black, letterSpacing = 2.4.sp)
                        Text("PERSISTENT WORLD", color = SecondaryText, fontSize = 9.sp, letterSpacing = 1.7.sp)
                    }
                }
                Spacer(Modifier.height(50.dp))
                Text(
                    if (registering) "Create your identity." else "Return to the world.",
                    color = PrimaryText,
                    fontSize = 34.sp,
                    lineHeight = 39.sp,
                    fontWeight = FontWeight.Bold
                )
                Spacer(Modifier.height(10.dp))
                Text(
                    if (registering) "Your character exists in a world that keeps moving."
                    else "Sign in to continue your persistent-world session.",
                    color = SecondaryText,
                    fontSize = 15.sp,
                    lineHeight = 22.sp
                )
                Spacer(Modifier.height(28.dp))
                AuthField(
                    username,
                    { username = it },
                    "USERNAME",
                    "Enter username",
                    { Icon(Icons.Default.Person, null) }
                )
                Spacer(Modifier.height(12.dp))
                AuthField(
                    password,
                    { password = it },
                    "PASSWORD",
                    "Enter password",
                    { Icon(Icons.Default.Lock, null) },
                    {
                        IconButton(onClick = { showPassword = !showPassword }) {
                            Icon(
                                if (showPassword) Icons.Default.VisibilityOff else Icons.Default.Visibility,
                                "Toggle password visibility"
                            )
                        }
                    },
                    KeyboardType.Password,
                    if (showPassword) VisualTransformation.None else PasswordVisualTransformation()
                )
                Spacer(Modifier.height(16.dp))
                Surface(
                    Modifier.fillMaxWidth(),
                    shape = RoundedCornerShape(17.dp),
                    color = SurfaceRaised,
                    border = BorderStroke(
                        1.dp,
                        if (state.serverConnected) Color(0x4062E6A8) else Border
                    )
                ) {
                    Row(Modifier.padding(14.dp), verticalAlignment = Alignment.CenterVertically) {
                        Icon(
                            if (state.serverConnected) Icons.Default.Wifi else Icons.Default.CloudOff,
                            null,
                            tint = if (state.serverConnected) Success else SecondaryText
                        )
                        Spacer(Modifier.width(10.dp))
                        Column(Modifier.weight(1f)) {
                            Text(
                                when {
                                    state.checkingServer -> "CHECKING SERVER…"
                                    state.serverConnected -> "SERVER ONLINE"
                                    else -> "SERVER NOT CONNECTED"
                                },
                                fontSize = 11.sp,
                                fontWeight = FontWeight.Bold,
                                color = if (state.serverConnected) Success else PrimaryText
                            )
                            Text(
                                state.serverMessage.ifBlank { "Configure the server address to connect." },
                                fontSize = 10.sp,
                                color = SecondaryText,
                                maxLines = 2
                            )
                        }
                        TextButton(onClick = { showServerSettings = !showServerSettings }) {
                            Text(if (showServerSettings) "HIDE" else "SERVER")
                        }
                    }
                }
                AnimatedVisibility(showServerSettings) {
                    Column {
                        Spacer(Modifier.height(10.dp))
                        AuthField(
                            serverUrl,
                            {
                                serverUrl = it
                                viewModel.setServerUrl(it)
                            },
                            "SERVER ADDRESS",
                            if (viewModel.physicalDevice) "http://192.168.x.x:8080" else "http://10.0.2.2:8080",
                            keyboardType = KeyboardType.Uri
                        )
                        Spacer(Modifier.height(6.dp))
                        Text(
                            if (viewModel.physicalDevice)
                                "Physical phone: use the computer's LAN address running Worldbuster."
                            else
                                "Emulator: 10.0.2.2 reaches the development machine.",
                            color = SecondaryText,
                            fontSize = 10.sp,
                            lineHeight = 15.sp
                        )
                        TextButton(
                            enabled = !state.checkingServer && serverUrl.isNotBlank(),
                            onClick = { viewModel.testConnection(serverUrl) }
                        ) {
                            Text("TEST CONNECTION", color = AccentBright)
                        }
                    }
                }
                if (error.isNotBlank()) {
                    Spacer(Modifier.height(10.dp))
                    Text(error, color = if (error.startsWith("Account created")) Success else Danger, fontSize = 11.sp, lineHeight = 16.sp)
                }
            }
            Column {
                Button(
                    Modifier.fillMaxWidth().height(56.dp),
                    enabled = !state.authBusy && username.length >= 3 && password.length >= 8 && serverUrl.isNotBlank(),
                    shape = RoundedCornerShape(17.dp),
                    onClick = {
                        viewModel.setServerUrl(serverUrl)
                        viewModel.authenticate(username, password, registering)
                    }
                ) {
                    Text(
                        if (state.authBusy) "CONNECTING…"
                        else if (registering) "CREATE ACCOUNT"
                        else "ENTER WORLD",
                        fontWeight = FontWeight.Bold,
                        letterSpacing = 1.sp
                    )
                    if (!state.authBusy) {
                        Spacer(Modifier.width(8.dp))
                        Icon(Icons.Default.ArrowForward, null)
                    } else {
                        Spacer(Modifier.width(10.dp))
                        CircularProgressIndicator(
                            Modifier.size(18.dp),
                            strokeWidth = 2.dp,
                            color = Color.White
                        )
                    }
                }
                Spacer(Modifier.height(8.dp))
                TextButton(
                    Modifier.fillMaxWidth(),
                    onClick = {
                        registering = !registering
                        password = ""
                    }
                ) {
                    Text(
                        if (registering) "BACK TO SIGN IN" else "CREATE NEW IDENTITY",
                        color = AccentBright,
                        fontWeight = FontWeight.SemiBold
                    )
                }
            }
        }
    }
}

@Composable
private fun AuthField(
    value: String,
    onValueChange: (String) -> Unit,
    label: String,
    placeholder: String,
    leadingIcon: @Composable (() -> Unit)? = null,
    trailingIcon: @Composable (() -> Unit)? = null,
    keyboardType: KeyboardType = KeyboardType.Text,
    visualTransformation: VisualTransformation = VisualTransformation.None
) {
    OutlinedTextField(
        value = value,
        onValueChange = onValueChange,
        modifier = Modifier.fillMaxWidth(),
        label = { Text(label, letterSpacing = 1.1.sp, fontSize = 10.sp) },
        placeholder = { Text(placeholder, color = SecondaryText) },
        leadingIcon = leadingIcon,
        trailingIcon = trailingIcon,
        singleLine = true,
        keyboardOptions = KeyboardOptions(keyboardType = keyboardType),
        visualTransformation = visualTransformation,
        shape = RoundedCornerShape(16.dp),
        colors = OutlinedTextFieldDefaults.colors(
            focusedBorderColor = Accent,
            unfocusedBorderColor = Border,
            focusedLabelColor = AccentBright,
            unfocusedLabelColor = SecondaryText,
            focusedLeadingIconColor = AccentBright,
            unfocusedLeadingIconColor = SecondaryText,
            focusedTextColor = PrimaryText,
            unfocusedTextColor = PrimaryText,
            cursorColor = AccentBright
        )
    )
}
