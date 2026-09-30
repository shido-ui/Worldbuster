package com.worldbuster.mobile

import org.json.JSONArray
import org.json.JSONObject

data class WorldView(
    val tick: Long = 0,
    val day: Int = 1,
    val time: String = "--:--",
    val online: Int = 0,
    val status: String = "UNKNOWN"
)

data class PlayerView(
    val name: String = "Unknown",
    val level: Int = 1,
    val xp: Long = 0,
    val cash: Long = 0,
    val energy: Int = 0,
    val inventoryUsed: Int = 0,
    val inventoryCapacity: Int = 100
)

data class MarketAssetView(
    val symbol: String,
    val name: String,
    val category: String,
    val price: Long,
    val supply: Long,
    val demand: Long
) {
    val pressure: Float
        get() = if (supply + demand == 0L) 0f else demand.toFloat() / (supply + demand).toFloat()
}

data class MissionView(
    val title: String,
    val description: String,
    val rewardCash: Long,
    val rewardXp: Long,
    val target: Long,
    val type: String
)

data class AchievementView(
    val title: String,
    val description: String,
    val progress: Long,
    val target: Long,
    val points: Int
)

data class OrganizationView(
    val name: String,
    val type: String,
    val level: Int,
    val reputation: Int,
    val members: Int
)

data class NewsView(val headline: String, val body: String, val category: String, val importance: Int)
data class EventView(val code: String, val status: String, val severity: Int, val location: String)

data class DashboardView(
    val world: WorldView = WorldView(),
    val player: PlayerView = PlayerView(),
    val market: List<MarketAssetView> = emptyList(),
    val missions: List<MissionView> = emptyList(),
    val achievements: List<AchievementView> = emptyList(),
    val organizations: List<OrganizationView> = emptyList(),
    val news: List<NewsView> = emptyList(),
    val events: List<EventView> = emptyList(),
    val lastSync: String = ""
)

private fun JSONObject.long(name: String, fallback: Long = 0) =
    if (has(name)) optLong(name, fallback) else fallback

private fun JSONObject.int(name: String, fallback: Int = 0) =
    if (has(name)) optInt(name, fallback) else fallback

private fun JSONObject.text(name: String, fallback: String = "") =
    optString(name, fallback).takeIf { it.isNotBlank() } ?: fallback

private fun JSONArray.objects(): List<JSONObject> =
    (0 until length()).mapNotNull { optJSONObject(it) }

fun parseWorld(body: String): WorldView {
    val o = JSONObject(body)
    return WorldView(o.long("tick"), o.int("day", 1), o.text("time", "--:--"), o.int("onlineCount"), o.text("status", "UNKNOWN"))
}

fun parsePlayerDashboard(body: String): PlayerView {
    val root = JSONObject(body)
    val profile = root.optJSONObject("profile") ?: JSONObject()
    val economy = root.optJSONObject("economy") ?: JSONObject()
    val inventory = root.optJSONObject("inventory") ?: JSONObject()
    val stacks = inventory.optJSONArray("stacks")
    val used = stacks?.objects()?.sumOf { it.int("quantity") } ?: 0
    val capacity = inventory.int("capacity", 100)
    return PlayerView(
        profile.text("displayName", "Unknown"),
        profile.int("level", 1),
        profile.long("xp"),
        economy.long("balance", profile.long("cash")),
        profile.int("energy"),
        used,
        capacity
    )
}

fun parseMarket(body: String): List<MarketAssetView> =
    (JSONObject(body).optJSONArray("assets") ?: JSONArray()).objects().map {
        MarketAssetView(
            it.text("symbol"),
            it.text("name", it.text("symbol")),
            it.text("category", "goods"),
            it.long("currentPrice"),
            it.long("supply"),
            it.long("demand")
        )
    }

fun parseMissions(body: String): List<MissionView> =
    (JSONObject(body).optJSONArray("missions") ?: JSONArray()).objects().map {
        MissionView(it.text("title"), it.text("description"), it.long("rewardCash"), it.long("rewardXP"), it.long("targetValue", 1), it.text("missionType", "MISSION"))
    }

fun parseAchievements(body: String): List<AchievementView> =
    (JSONObject(body).optJSONArray("achievements") ?: JSONArray()).objects().map {
        AchievementView(it.text("title"), it.text("description"), it.long("progress"), it.long("targetProgress", 1), it.int("points"))
    }

fun parseOrganizations(body: String): List<OrganizationView> =
    (JSONObject(body).optJSONArray("organizations") ?: JSONArray()).objects().map {
        OrganizationView(it.text("name"), it.text("type", "ORGANIZATION"), it.int("level", 1), it.int("reputation"), it.int("maxMembers"))
    }

fun parseNews(body: String): List<NewsView> =
    (JSONObject(body).optJSONArray("news") ?: JSONArray()).objects().map {
        NewsView(it.text("headline"), it.text("body"), it.text("category", "WORLD"), it.int("importance", 1))
    }

fun parseEvents(body: String): List<EventView> =
    (JSONObject(body).optJSONArray("events") ?: JSONArray()).objects().map {
        EventView(it.text("code"), it.text("status"), it.int("severity", 1), it.text("locationId", "WORLD"))
    }
