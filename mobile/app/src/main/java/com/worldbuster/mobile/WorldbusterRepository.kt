package com.worldbuster.mobile

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.withContext
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

data class DashboardRefresh(
    val dashboard: DashboardView,
    val syncState: SyncState,
    val message: String,
    val unauthorized: Boolean
)

class WorldbusterRepository(context: android.content.Context) {
    private val api = ApiClient(context)

    val serverUrl: String get() = api.baseUrl()
    fun isPhysicalDevice(): Boolean = api.isPhysicalDevice()
    fun hasSession(): Boolean = api.hasSession()
    fun setServerUrl(value: String) { api.setBaseUrl(value) }

    suspend fun testConnection(): Result<ApiResponse> = withContext(Dispatchers.IO) { api.health() }
    suspend fun authenticate(username: String, password: String, register: Boolean): Result<ApiResponse> = withContext(Dispatchers.IO) {
        if (register) api.register(username, password) else api.login(username, password)
    }
    suspend fun logout(): Result<ApiResponse> = withContext(Dispatchers.IO) { api.logout() }

    suspend fun refresh(previous: DashboardView): DashboardRefresh = withContext(Dispatchers.IO) {
        if (!api.isNetworkAvailable() && api.baseUrl().isBlank()) {
            return@withContext DashboardRefresh(previous, SyncState.OFFLINE, "No network connection and no server address is configured.", false)
        }
        val endpoints = listOf(
            "world" to "/api/v1/world", "dashboard" to "/api/v1/player/dashboard",
            "market" to "/api/v1/market/assets", "missions" to "/api/v1/missions",
            "achievements" to "/api/v1/achievements", "organizations" to "/api/v1/organizations",
            "news" to "/api/v1/news", "events" to "/api/v1/world-events"
        )
        val results = coroutineScope { endpoints.map { endpoint -> async { endpoint.first to api.get(endpoint.second) } }.map { it.await() } }
        var next = previous
        var unauthorized = false
        var hadFailure = false
        var usedCache = false
        for ((key, result) in results) {
            val response = result.getOrNull()
            if (response?.code == 401) { unauthorized = true; continue }
            if (result.isSuccess && response != null && response.code in 200..299) {
                val parsed = parseEndpoint(next, key, response.body)
                if (parsed.isSuccess) { api.cache(key, response); next = parsed.getOrThrow() }
                else {
                    hadFailure = true
                    api.cached(key)?.let { cached -> parseEndpoint(next, key, cached.body).getOrNull()?.let { usedCache = true; next = it } }
                }
            } else {
                hadFailure = true
                api.cached(key)?.let { cached -> parseEndpoint(next, key, cached.body).getOrNull()?.let { usedCache = true; next = it } }
            }
        }
        val state = when { unauthorized -> SyncState.AUTH; !hadFailure -> SyncState.CONNECTED; usedCache -> SyncState.DEGRADED; else -> SyncState.OFFLINE }
        val message = when {
            unauthorized -> "Your session expired. Sign in again."
            !hadFailure -> ""
            usedCache -> "Some live data is unavailable. Showing the last successful sync."
            else -> "World data is temporarily unavailable. Check the server connection and retry."
        }
        DashboardRefresh(next.copy(lastSync = if (!hadFailure && !unauthorized) nowLabel() else previous.lastSync), state, message, unauthorized)
    }

    private fun parseEndpoint(data: DashboardView, key: String, body: String): Result<DashboardView> = runCatching {
        when (key) {
            "world" -> data.copy(world = parseWorld(body))
            "dashboard" -> data.copy(player = parsePlayerDashboard(body))
            "market" -> data.copy(market = parseMarket(body))
            "missions" -> data.copy(missions = parseMissions(body))
            "achievements" -> data.copy(achievements = parseAchievements(body))
            "organizations" -> data.copy(organizations = parseOrganizations(body))
            "news" -> data.copy(news = parseNews(body))
            "events" -> data.copy(events = parseEvents(body))
            else -> error("Unknown dashboard endpoint: " + key)
        }
    }

    private fun nowLabel(): String = SimpleDateFormat("HH:mm:ss", Locale.getDefault()).format(Date())
}