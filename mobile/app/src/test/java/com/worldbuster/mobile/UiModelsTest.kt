package com.worldbuster.mobile

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class UiModelsTest {
    @Test
    fun parsesWorldState() {
        val world = parseWorld("""{"tick":42,"day":3,"time":"12:30","onlineCount":17,"status":"RUNNING"}""")
        assertEquals(42L, world.tick)
        assertEquals(3, world.day)
        assertEquals(17, world.online)
        assertEquals("RUNNING", world.status)
    }

    @Test
    fun parsesPlayerDashboardAndInventoryUsage() {
        val player = parsePlayerDashboard(
            """{"profile":{"displayName":"Shido","level":4,"xp":1250,"energy":80},"economy":{"balance":900},"inventory":{"capacity":20,"stacks":[{"quantity":2},{"quantity":3}]}}"""
        )
        assertEquals("Shido", player.name)
        assertEquals(4, player.level)
        assertEquals(900L, player.cash)
        assertEquals(5, player.inventoryUsed)
        assertEquals(20, player.inventoryCapacity)
    }

    @Test
    fun parsesMarketAndCalculatesPressure() {
        val assets = parseMarket("""{"assets":[{"symbol":"WBX","name":"WBX Credit","category":"currency","currentPrice":100,"supply":25,"demand":75}]}""")
        assertEquals(1, assets.size)
        assertEquals("WBX", assets[0].symbol)
        assertEquals(0.75f, assets[0].pressure, 0.001f)
    }

    @Test
    fun missingCollectionsBecomeEmpty() {
        assertTrue(parseMarket("{}").isEmpty())
        assertTrue(parseMissions("{}").isEmpty())
        assertTrue(parseAchievements("{}").isEmpty())
        assertTrue(parseOrganizations("{}").isEmpty())
        assertTrue(parseNews("{}").isEmpty())
        assertTrue(parseEvents("{}").isEmpty())
    }

    @Test
    fun parsesMissionAndAchievementContracts() {
        val missions = parseMissions("""{"missions":[{"title":"Train","description":"Practice","rewardCash":50,"rewardXP":20,"targetValue":5,"missionType":"TRAINING"}]}""")
        assertEquals(50L, missions[0].rewardCash)
        assertEquals(20L, missions[0].rewardXp)
        assertEquals(5L, missions[0].target)

        val achievements = parseAchievements("""{"achievements":[{"title":"First","description":"Begin","progress":2,"targetProgress":10,"points":25}]}""")
        assertEquals(2L, achievements[0].progress)
        assertEquals(10L, achievements[0].target)
        assertEquals(25, achievements[0].points)
    }

    @Test
    fun parserRejectsMalformedJson() {
        var failed = false
        try {
            parseWorld("{not-json")
        } catch (_: Exception) {
            failed = true
        }
        assertTrue(failed)
    }
}
