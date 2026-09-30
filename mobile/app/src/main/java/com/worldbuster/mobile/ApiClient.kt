package com.worldbuster.mobile

import java.net.HttpURLConnection
import java.net.URL

class ApiClient(private val baseUrl: String) {
    fun get(path: String): Result<String> = runCatching {
        val connection = URL(baseUrl.trimEnd('/') + path).openConnection() as HttpURLConnection
        connection.requestMethod = "GET"
        connection.connectTimeout = 5000
        connection.readTimeout = 5000
        try {
            val code = connection.responseCode
            val stream = if (code in 200..299) connection.inputStream else connection.errorStream
            val body = stream?.bufferedReader()?.use { it.readText() }.orEmpty()
            if (code !in 200..299) error("HTTP $code: $body")
            body
        } finally { connection.disconnect() }
    }
}
