package com.worldbuster.mobile

import android.content.Context
import java.net.HttpURLConnection
import java.net.URL

data class ApiResponse(val code:Int,val body:String)

class ApiClient(context:Context, private val baseUrl:String) {
    private val prefs=context.getSharedPreferences("worldbuster",Context.MODE_PRIVATE)

    private fun open(path:String,method:String):HttpURLConnection{
        val c=URL(baseUrl.trimEnd('/')+path).openConnection() as HttpURLConnection
        c.requestMethod=method
        c.connectTimeout=5000
        c.readTimeout=5000
        c.useCaches=false
        prefs.getString("session_cookie",null)?.let{c.setRequestProperty("Cookie",it)}
        return c
    }

    fun request(path:String,method:String="GET",body:String?=null):Result<ApiResponse> = runCatching {
        val c=open(path,method)
        try{
            if(body!=null){
                c.doOutput=true
                c.setRequestProperty("Content-Type","application/json")
                c.outputStream.use{it.write(body.toByteArray(Charsets.UTF_8))}
            }
            val code=c.responseCode
            val stream=if(code in 200..299)c.inputStream else c.errorStream
            val text=stream?.bufferedReader()?.use{it.readText()}.orEmpty()
            c.headerFields["Set-Cookie"]?.firstOrNull()?.substringBefore(';')?.let{
                prefs.edit().putString("session_cookie",it).apply()
            }
            ApiResponse(code,text)
        }finally{c.disconnect()}
    }

    fun hasSession():Boolean=!prefs.getString("session_cookie",null).isNullOrBlank()
    fun clearSession(){prefs.edit().remove("session_cookie").apply()}
    fun login(username:String,password:String)=request("/api/v1/auth/login","POST","{\"username\":\""+escape(username)+"\",\"password\":\""+escape(password)+"\"}")
    fun register(username:String,password:String)=request("/api/v1/auth/register","POST","{\"username\":\""+escape(username)+"\",\"password\":\""+escape(password)+"\"}")
    fun logout():Result<ApiResponse>{val r=request("/api/v1/auth/logout","POST");clearSession();return r}
    fun get(path:String)=request(path)
    private fun escape(value:String)=value.replace("\\","\\\\").replace("\"","\\\"")
}
