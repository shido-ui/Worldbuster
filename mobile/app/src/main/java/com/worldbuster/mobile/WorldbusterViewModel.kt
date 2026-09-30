package com.worldbuster.mobile

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

data class WorldbusterUiState(
    val authenticated: Boolean = false,
    val dashboard: DashboardView = DashboardView(),
    val syncState: SyncState = SyncState.AUTH,
    val loading: Boolean = false,
    val error: String = "",
    val authBusy: Boolean = false,
    val checkingServer: Boolean = false,
    val serverConnected: Boolean = false,
    val serverMessage: String = ""
)

class WorldbusterViewModel(application: Application) : AndroidViewModel(application) {
    private val repository = WorldbusterRepository(application)

    private val _uiState = MutableStateFlow(
        WorldbusterUiState(
            authenticated = repository.hasSession(),
            syncState = if (repository.hasSession()) SyncState.OFFLINE else SyncState.AUTH
        )
    )
    val uiState: StateFlow<WorldbusterUiState> = _uiState.asStateFlow()

    val serverUrl: String
        get() = repository.serverUrl

    val physicalDevice: Boolean
        get() = repository.isPhysicalDevice()

    init {
        if (_uiState.value.authenticated) refresh()
    }

    fun setServerUrl(value: String) {
        repository.setServerUrl(value)
        _uiState.value = _uiState.value.copy(serverConnected = false, serverMessage = "")
    }

    fun testConnection(value: String = serverUrl) {
        if (value.isBlank()) {
            _uiState.value = _uiState.value.copy(serverConnected = false, serverMessage = "Enter a server address first.")
            return
        }
        repository.setServerUrl(value)
        _uiState.value = _uiState.value.copy(checkingServer = true, serverConnected = false, serverMessage = "")
        viewModelScope.launch {
            val result = repository.testConnection()
            val connected = result.isSuccess && result.getOrNull()?.code in 200..299
            val message = if (connected) "Server reachable" else {
                result.exceptionOrNull()?.let { ApiClientErrorFormatter.message(it) }
                    ?: result.getOrNull()?.let { ApiClientErrorFormatter.http(it) }
                    ?: "Unable to reach the server."
            }
            _uiState.value = _uiState.value.copy(checkingServer = false, serverConnected = connected, serverMessage = message)
        }
    }

    fun authenticate(username: String, password: String, register: Boolean) {
        _uiState.value = _uiState.value.copy(authBusy = true, error = "")
        viewModelScope.launch {
            val result = repository.authenticate(username, password, register)
            val response = result.getOrNull()
            if (result.isFailure) {
                val message = ApiClientErrorFormatter.message(result.exceptionOrNull()!!)
                _uiState.value = _uiState.value.copy(authBusy = false, serverConnected = false, serverMessage = message, error = message)
                return@launch
            }
            if (response == null || response.code !in 200..299) {
                val message = response?.let(ApiClientErrorFormatter::http) ?: "The server returned no response."
                _uiState.value = _uiState.value.copy(authBusy = false, error = message)
                return@launch
            }
            if (register) {
                _uiState.value = _uiState.value.copy(authBusy = false, serverConnected = true, serverMessage = "Server reachable", error = "Account created. Sign in to continue.")
            } else {
                _uiState.value = _uiState.value.copy(authenticated = true, authBusy = false, serverConnected = true, serverMessage = "Server reachable", syncState = SyncState.OFFLINE, error = "")
                refresh()
            }
        }
    }

    fun refresh() {
        if (!_uiState.value.authenticated) return
        _uiState.value = _uiState.value.copy(loading = true)
        viewModelScope.launch {
            val result = repository.refresh(_uiState.value.dashboard)
            if (result.unauthorized) {
                _uiState.value = _uiState.value.copy(authenticated = false, dashboard = DashboardView(), syncState = SyncState.AUTH, loading = false, error = result.message)
            } else {
                _uiState.value = _uiState.value.copy(dashboard = result.dashboard, syncState = result.syncState, loading = false, error = result.message)
            }
        }
    }

    fun logout() {
        viewModelScope.launch {
            repository.logout()
            _uiState.value = WorldbusterUiState(syncState = SyncState.AUTH)
        }
    }
}

object ApiClientErrorFormatter {
    fun message(error: Throwable): String = when (error) {
        is java.net.UnknownHostException -> "The server address could not be resolved."
        is java.net.ConnectException -> "The server refused the connection. Check that Worldbuster is running and reachable."
        is java.net.SocketTimeoutException -> "The server timed out. Check the address and local network."
        else -> error.message ?: "Unable to reach the Worldbuster server."
    }

    fun http(response: ApiResponse): String {
        val detail = runCatching { org.json.JSONObject(response.body).optString("error").trim() }.getOrNull().orEmpty()
        return if (detail.isNotBlank()) "Server " + response.code + ": " + detail else "Server returned HTTP " + response.code + "."
    }
}