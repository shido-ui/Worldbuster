package com.worldbuster.mobile

import android.content.Context
import android.os.Build
import java.net.ConnectException
import java.net.HttpURLConnection
import java.net.SocketTimeoutException
import java.net.URL
import java.net.UnknownHostException

data class ApiResponse(val code: Int, val body: String)

class ApiClient(context: Context) {
    private val prefs = context.getSharedPreferences("worldbuster", Context.MODE_PRIVATE)
    private val emulatorDefault = "http://10.0.2.2:8080"

    fun isPhysicalDevice(): Boolean {
        val fingerprint = Build.FINGERPRINT.lowercase()
        val model = Build.MODEL.lowercase()
        val product = Build.PRODUCT.lowercase()
        return !(fingerprint.contains("generic") || fingerprint.contains("emulator") ||
            model.contains("emulator") || model.contains("android sdk") || product.contains("sdk"))
    }

    fun baseUrl(): String =
        prefs.getString("server_url", null) ?: if (isPhysicalDevice()) "" else emulatorDefault

    fun setBaseUrl(value: String) {
        prefs.edit().putString("server_url", value.trim().trimEnd('/')).apply()
    }

    private fun open(path: String, method: String): HttpURLConnection {
        val root = baseUrl()
        require(root.isNotBlank()) { "Server address is not configured." }
        require(root.startsWith("http://") || root.startsWith("https://")) {
            "Server address must start with http:// or https://."
        }
        val c = URL(root + path).openConnection() as HttpURLConnection
        c.requestMethod = method
        c.connectTimeout = 8000
        c.readTimeout = 8000
        c.useCaches = false
        c.setRequestProperty("Accept", "application/json")
        prefs.getString("session_cookie", null)?.let { c.setRequestProperty("Cookie", it) }
        return c
    }

    fun request(path: String, method: String = "GET", body: String? = null): Result<ApiResponse> = runCatching {
        val c = open(path, method)
        try {
            if (body != null) {
                c.doOutput = true
                c.setRequestProperty("Content-Type", "application/json")
                c.outputStream.use { it.write(body.toByteArray(Charsets.UTF_8)) }
            }
            val code = c.responseCode
            val stream = if (code in 200..299) c.inputStream else c.errorStream
            val text = stream?.bufferedReader()?.use { it.readText() }.orEmpty()
            if (code in 200..299) {
                c.headerFields["Set-Cookie"]?.firstOrNull()?.substringBefore(';')?.let {
                    prefs.edit().putString("session_cookie", it).apply()
                }
            }
            ApiResponse(code, text)
        } finally {
            c.disconnect()
        }
    }

    fun health(): Result<ApiResponse> = request("/health")
    fun hasSession(): Boolean = !prefs.getString("session_cookie", null).isNullOrBlank()
    fun clearSession() { prefs.edit().remove("session_cookie").apply() }
    fun login(username: String, password: String) =
        request("/api/v1/auth/login", "POST", credentials(username, password))
    fun register(username: String, password: String) =
        request("/api/v1/auth/register", "POST", credentials(username, password))
    fun logout(): Result<ApiResponse> {
        val result = request("/api/v1/auth/logout", "POST")
        clearSession()
        return result
    }
    fun get(path: String) = request(path)

    fun describeFailure(result: Result<ApiResponse>): String {
        val error = result.exceptionOrNull() ?: return "The server did not respond."
        return when (error) {
            is UnknownHostException -> "Server address could not be resolved."
            is ConnectException -> "Server refused the connection. Check that Worldbuster is running and the address is reachable."
            is SocketTimeoutException -> "Server timed out. Check the server address and local network."
            else -> error.message ?: "Unable to reach the Worldbuster server."
        }
    }

    fun describeHttpError(response: ApiResponse): String {
        val detail = response.body.trim().take(240)
        return if (detail.isNotBlank()) "Server returned " + response.code + ": " + detail
        else "Server returned HTTP " + response.code + "."
    }

    private fun credentials(username: String, password: String): String =
        "{\"username\":\"" + escape(username) + "\",\"password\":\"" + escape(password) + "\"}"

    private fun escape(value: String): String =
        value.replace("\\", "\\\\").replace("\"", "\\\"")
}
