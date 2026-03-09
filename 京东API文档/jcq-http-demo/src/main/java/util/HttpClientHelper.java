package util;

import sun.misc.BASE64Encoder;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.URISyntaxException;
import java.net.URL;
import java.util.HashMap;
import java.util.Map;


public class HttpClientHelper {
    private static HttpClientHelper instance = null;
    private static HttpClientUtil httpClientUtil = new HttpClientUtil();
    protected HttpClientHelper() {
    }
    public static HttpClientHelper getInstance() {
        if (instance == null) {
            instance = new HttpClientHelper();
        }
        return instance;
    }

    public String doHttpProxyGet(final String url, final String accessKey, final String datetime, final String signature) throws IOException,URISyntaxException {
        return doHttpProxyGet(url, accessKey, datetime, signature, null);
    }

    public String doHttpProxyGet(final String url, final String accessKey, final String datetime, final String signature, final String accessToken) throws IOException,URISyntaxException {
        HashMap<String,String> kvs = new HashMap<>();
        kvs.put("accessKey", accessKey);
        kvs.put("signature", signature);
        if(!datetime.isEmpty()) {
            kvs.put("dateTime", datetime);
        }
        if(accessToken!=null && !accessToken.isEmpty()){
            kvs.put("token",accessToken);
        }
        // return doMethod(url, "GET", kvs);
        return httpClientUtil.get(url, kvs);
    }


    public String doGet(final String url, final String headerJvesselParams) throws IOException {
        HashMap<String,String> kvs = new HashMap<>();
        kvs.put("JVESSEL-Params", headerJvesselParams);
        return doMethod(url, "GET", kvs);
    }

    public String doDelete(final String url,final String headerJvesselParams) throws IOException {
        HashMap<String,String> kvs = new HashMap<>();
        kvs.put("JVESSEL-Params", headerJvesselParams);
        return doMethod(url, "DELETE",kvs);
    }

    // connection.setRequestProperty("JVESSEL-Params", headerJvesselParams);
    // connection.addRequestProperty("accesskey", accessKey);
    // connection.addRequestProperty("signature", signature);
    // connection.addRequestProperty("datetime", datetime);
    private String doMethod(final String url, final String httpMethod, final HashMap<String,String> kvs) throws IOException{
        BufferedReader in = null;
        StringBuilder resultJsonString = new StringBuilder();
        try {
            URL realUrl = new URL(url);
            HttpURLConnection connection = (HttpURLConnection)realUrl.openConnection();

            connection.setRequestProperty("Content-Type","application/json; charset=utf-8");
            connection.setRequestProperty("x-jdcloud-service","jcq");
            String data = (new BASE64Encoder()).encode(ConfigHelper.getUserPin().getBytes());
            connection.setRequestProperty("x-jdcloud-pin", data);

//            if(kvs!=null)
//            {
//                for(String key:kvs.keySet()){
//                    connection.setRequestProperty(key,kvs.get(key));
//                }
//            }


            connection.setDoInput(true);
            connection.setDoOutput(false);
            connection.setUseCaches(false);
            connection.setRequestMethod(httpMethod);
            in = new BufferedReader(new InputStreamReader(connection.getInputStream()));
            String line;
            while ((line = in.readLine()) != null) {
                resultJsonString.append(line);
            }
        } catch (MalformedURLException e) {
            throw e;
        } catch (IOException e) {
            throw e;
        } finally {
            try {
                if (in != null) {
                    in.close();
                }
            } catch (IOException e) {
            }
        }
        return resultJsonString.toString();
    }

    public String doHttpProxyPost(final String url, final String objectJson, final String accessKey, final String signature, final String datetime) throws Exception{
        return doHttpProxyPost(url, objectJson, accessKey, signature, datetime, null);
    }

    public String doHttpProxyPost(final String url, final String objectJson, final String accessKey, final String signature, final String datetime, final String accessToken) throws Exception{
        HashMap<String,String> kvs = new HashMap<>();
        kvs.put("accessKey", accessKey);
        kvs.put("signature", signature);
        if(!datetime.isEmpty()) {
            kvs.put("dateTime", datetime);
        }
        if(accessToken!=null && !accessToken.isEmpty()) {
            kvs.put("token", accessToken);
        }

        // return doPost(url, objectJson, kvs);
        return httpClientUtil.postJSON(url, objectJson, kvs);
    }

    public String doPut(final String url, final String objectJson, Map<String,String> kvs) throws Exception {
        // return doPost(url, objectJson, kvs);
        return httpClientUtil.putJSON(url, objectJson, kvs);
    }

    public String doPost(final String url, final String objectJson, Map<String,String> kvs) throws Exception {
        // return doPost(url, objectJson, kvs);
        return httpClientUtil.postJSON(url, objectJson, kvs);
    }

    public String doDelete(final String url, final String objectJson, Map<String,String> kvs) throws Exception {
        // return doPost(url, objectJson, kvs);
        return httpClientUtil.deleteJSON(url, objectJson, kvs);
    }

    //public String doPost(final String url, final String objectJson, final HashMap<String,String> headerJvesselParams) throws IOException {
    public String doPost(final String url, final String objectJson, final String userPin) throws IOException {
        OutputStream out = null;
        BufferedReader in = null;
        StringBuilder resultJsonString = new StringBuilder();
        try {
            URL hostUrl = new URL(url);
            HttpURLConnection connection = (HttpURLConnection) hostUrl.openConnection();

            connection.setRequestProperty("Content-Type","application/json; charset=utf-8");
            connection.setRequestProperty("x-jdcloud-service","jcq");
            String data = (new BASE64Encoder()).encode(userPin.getBytes());
            connection.setRequestProperty("x-jdcloud-pin", data);

            connection.setDoInput(true);
            connection.setDoOutput(true);
            connection.setUseCaches(false);
            connection.setRequestMethod("POST");

            if(objectJson!=null && !objectJson.isEmpty()) {
                out = connection.getOutputStream();
                out.write(objectJson.getBytes("UTF-8"));
                out.flush();
            }
            InputStream inputStream = connection.getInputStream();
            in = new BufferedReader(new InputStreamReader(inputStream));

            String line;
            while ((line = in.readLine()) != null) {
                resultJsonString.append(line);
            }
            connection.disconnect();
        } catch (MalformedURLException e) {
            throw e;
        } catch (IOException e) {
            throw e;
        }
        catch (Exception e){
//            logger.error("%s", e.getMessage());
        }
        finally {
            try {
                if (out != null) {
                    out.close();
                }
                if (in != null) {
                    in.close();
                }
            } catch (IOException e) {
//                logger.error("Exception occurs in closing connection for POST request: {}", e.toString());
            }
        }
        return resultJsonString.toString();
    }

//    public String doPatch(final String url, final String headerJvessel) throws Exception{
//        return doPatch(url,null,);
//    }

    public String doPut(final String url, final String headerJvessel) throws Exception{
        return doPut(url,null,headerJvessel);
    }

    public String doPatch(final String url, final String objectJson, Map<String,String> headers) throws IOException {
        // return doPatchOrPut("PATCH", url, objectJson, headerJvessel);
        return httpClientUtil.patchJSON(url, objectJson, headers);
    }

    public String doPut(final String url, final String objectJson, final String headerJvessel) throws IOException {
        return doPatchOrPut("PUT", url, objectJson, headerJvessel);
    }

    private String doPatchOrPut(final String method, final String url, final String objectJson, final String headerJvessel) throws IOException {
        OutputStream out = null;
        BufferedReader in = null;
        StringBuilder resultJsonString = new StringBuilder();
        try {
            URL hostUrl = new URL(url);
            HttpURLConnection connection = (HttpURLConnection) hostUrl.openConnection();
            connection.setRequestProperty("Content-Type","application/json; charset=utf-8");
            connection.setRequestProperty("JVESSEL-Params",headerJvessel);
            connection.setDoInput(true);
            connection.setDoOutput(true);
            connection.setUseCaches(false);
            connection.setRequestMethod(method);
            if(objectJson != null && !objectJson.isEmpty()) {
                out = connection.getOutputStream();
                out.write(objectJson.getBytes("UTF-8"));
                out.flush();
            }

            in = new BufferedReader(new InputStreamReader(connection.getInputStream()));
            String line;
            while ((line = in.readLine()) != null) {
                resultJsonString.append(line);
            }
            connection.disconnect();
        } catch (MalformedURLException e) {
//            logger.error("Exception occurs in constructing URL object for PUTT request: {}", e.toString());
            throw e;
        } catch (IOException e) {
//            logger.error("Exception occurs in constructing URLConnection object for PUTT request: {}", e.toString());
            throw e;
        } finally {
            try {
                if (out != null) {
                    out.close();
                }
                if (in != null) {
                    in.close();
                }
            } catch (IOException e) {
//                logger.error("Exception occurs in closing connection for PUTT request: {}", e.toString());
            }
        }
        return resultJsonString.toString();
    }
}
