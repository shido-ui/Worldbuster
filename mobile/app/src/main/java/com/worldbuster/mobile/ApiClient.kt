package com.worldbuster.mobile

import java.net.HttpURLConnection
import java.net.URL

class ApiClient(private val baseUrl: String) {
    fun get(path: String): String {
        val c = URL(baseUrl.trimEnd('/') + path).openConnection() as HttpURLConnection
        c.requestMethod = "GET"
        c.connectTimeout = 5000
        c.readTimeout = 5000
        return c.inputStream.bufferedReader().use { it.readText() }
    }
}
