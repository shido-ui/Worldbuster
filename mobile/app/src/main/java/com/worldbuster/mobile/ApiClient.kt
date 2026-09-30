package com.worldbuster.mobile

import android.content.Context
import android.net.ConnectivityManager
import android.net.NetworkCapabilities
import android.os.Build
import org.json.JSONObject
import java.net.ConnectException
import java.net.HttpURLConnection
import java.net.SocketTimeoutException
import java.net.URL
import java.net.UnknownHostException
import kotlin.math.min

data class ApiResponse(val code: Int, val body: String)

class ApiClient(context: Context) {
    private val appContext = context.applicationContext
    private val prefs = appContext.getSharedPreferences("worldbuster", Context.MODE_PRIVATE)
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
        prefs.edit().putString("server_url", normalizeBaseUrl(value)).apply()
    }

    fun isNetworkAvailable(): Boolean {
        val manager = appContext.getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager
        val network = manager.activeNetwork ?: return false
        val capabilities = manager.getNetworkCapabilities(network) ?: return false
        return capabilities.hasCapability(NetworkCapabilities.NET_CAPABILITY_INTERNET) &&
            capabilities.hasCapability(NetworkCapabilities.NET_CAPABILITY_VALIDATED)
    }

    private fun normalizeBaseUrl(value: String): String = value.trim().trimEnd('/')

    private fun open(path: String, method: String): HttpURLConnection {
        require(path.startsWith("/")) { "API path must start with /." }
        val root = normalizeBaseUrl(baseUrl())
        require(root.isNotBlank()) { "Server address is not configured." }
        require(root.startsWith("http://") || root.startsWith("https://")) {
            "Server address must start with http:// or https://."
        }
        val c = URL(root + path).openConnection() as HttpURLConnection
        c.requestMethod = method
        c.connectTimeout = 8000
        c.readTimeout = 8000
        c.useCaches = false
        c.instanceFollowRedirects = false
        c.setRequestProperty("Accept", "application/json")
        prefs.getString("session_cookie", null)?.let { c.setRequestProperty("Cookie", it) }
        return c
    }

    fun request(path: String, method: String = "GET", body: String? = null): Result<ApiResponse> {
        val attempts = if (method == "GET") 2 else 1
        var last: Result<ApiResponse> = Result.failure(IllegalStateException("request not attempted"))
        repeat(attempts) { attempt ->
            last = runCatching { execute(path, method, body) }
            val response = last.getOrNull()
            val retryable = last.isFailure ||
                response?.code == 408 || response?.code == 429 ||
                (response?.code ?: 0) in 500..599
            if (!retryable || attempt == attempts - 1) return last
            Thread.sleep(min(750L, 250L * (attempt + 1)))
        }
        return last
    }

    private fun execute(path: String, method: String, body: String?): ApiResponse {
        val c = open(path, method)
        try {
            if (body != null) {
                c.doOutput = true
                c.setRequestProperty("Content-Type", "application/json; charset=utf-8")
                c.outputStream.use { it.write(body.toByteArray(Charsets.UTF_8)) }
            }
            val code = c.responseCode
            val stream = if (code in 200..299) c.inputStream else c.errorStream
            val text = stream?.bufferedReader()?.use { it.readText() }.orEmpty()
            if (code in 200..299) {
                c.headerFields.entries.firstOrNull { it.key.equals("Set-Cookie", ignoreCase = true) }
                    ?.value?.firstOrNull()?.substringBefore(';')?.let {
                        prefs.edit().putString("session_cookie", it).apply()
                    }
            }
            return ApiResponse(code, text)
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

    fun cache(key: String, response: ApiResponse) {
        if (response.code in 200..299) prefs.edit().putString("cache:" + key, response.body).apply()
    }

    fun cached(key: String): ApiResponse? =
        prefs.getString("cache:" + key, null)?.let { ApiResponse(200, it) }

    fun describeFailure(result: Result<ApiResponse>): String {
        val error = result.exceptionOrNull() ?: return "The server did not respond."
        return when (error) {
            is UnknownHostException -> "The server address could not be resolved."
            is ConnectException -> "The server refused the connection. Check that Worldbuster is running and reachable."
            is SocketTimeoutException -> "The server timed out. Check the address and local network."
            else -> error.message ?: "Unable to reach the Worldbuster server."
        }
    }

    fun describeHttpError(response: ApiResponse): String {
        val detail = runCatching { JSONObject(response.body).optString("error").trim() }
            .getOrNull().orEmpty()
        return if (detail.isNotBlank()) "Server " + response.code + ": " + detail
        else "Server returned HTTP " + response.code + "."
    }

    private fun credentials(username: String, password: String): String =
        JSONObject().put("username", username).put("password", password).toString()
}