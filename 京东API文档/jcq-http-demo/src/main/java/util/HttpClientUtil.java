package util;

import org.apache.http.HttpResponse;
import org.apache.http.NameValuePair;
import org.apache.http.client.ResponseHandler;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.client.entity.UrlEncodedFormEntity;
import org.apache.http.client.methods.*;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.message.BasicNameValuePair;
import org.apache.http.util.EntityUtils;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;


public class HttpClientUtil {
    private HttpClientBuilder httpClientBuilder = HttpClientBuilder.create();
    private RequestConfig.Builder requestConfigBuilder = RequestConfig.custom();
    private int timeout = 1000;
    private Charset charset = Charset.forName("UTF-8");

    public HttpClientUtil() {
    }

    public HttpClientUtil(int timeout) {
        this(timeout, null);
    }

    public HttpClientUtil(int timeout, Charset charset) {
        this.timeout = timeout > 0 ? timeout : this.timeout;
        this.charset = charset != null ? charset : this.charset;
        requestConfigBuilder.setConnectTimeout(timeout);
        requestConfigBuilder.setSocketTimeout(timeout);
        requestConfigBuilder.setConnectionRequestTimeout(timeout);
        httpClientBuilder.setDefaultRequestConfig(requestConfigBuilder.build());
    }

//    public String get(String url) throws IOException {
//        HttpGet httpGet = new HttpGet(url);
//        return request(httpGet, getStringResponseHandler());
//    }

    //url 包含？keyvalue
    //public String get(String url, Map<String,String> headers, Map<String, String> params) throws IOException, URISyntaxException {
    public String get(String url, Map<String,String> headers) throws IOException, URISyntaxException {
        HttpGet httpGet = new HttpGet("");
        if(headers!=null) {
            for (String key : headers.keySet()) {
                httpGet.setHeader(key, headers.get(key));
            }
        }

       // String query = EntityUtils.toString(new UrlEncodedFormEntity(toNameValuePairs(params), charset));
        httpGet.setURI(new URI(url));
       // httpGet.setURI(new URI(httpGet.getURI().toString() + "?" + query));
        return request(httpGet, getStringResponseHandler());
    }

    public String post(String url, Map<String,String> headers, Map<String, String> params) throws IOException {
        HttpPost httpPost = new HttpPost(url);

        httpPost.setEntity(new UrlEncodedFormEntity(toNameValuePairs(params), charset));
        return request(httpPost, getStringResponseHandler());
    }

    public String patch(String url, Map<String,String> headers, Map<String, String> params) throws IOException {
        //HttpPost httpPost = new HttpPost(url);
        HttpPatch httpPatch = new HttpPatch(url);

        httpPatch.setEntity(new UrlEncodedFormEntity(toNameValuePairs(params), charset));
        return request(httpPatch, getStringResponseHandler());
    }

    public String putJSON(String url, String json, Map<String,String> headers) throws IOException {
        return put(url, json, "application/json", headers);
    }

    public String postJSON(String url, String json, Map<String,String> headers) throws IOException {
        return post(url, json, "application/json", headers);
    }

    public String put(String url, String content, String contentType, Map<String,String> headers) throws IOException {
        HttpPut httpPut = new HttpPut(url);

        if(headers!=null) {
            for (String key : headers.keySet()) {
                httpPut.setHeader(key, headers.get(key));
            }
        }

        StringEntity stringEntity = new StringEntity(content, charset);
        if (contentType != null) {
            stringEntity.setContentType(contentType);
        }
        httpPut.setEntity(stringEntity);
        return request(httpPut, getStringResponseHandler());
    }

    public String post(String url, String content, String contentType, Map<String,String> headers) throws IOException {
        HttpPost httpPost = new HttpPost(url);

        if(headers!=null) {
            for (String key : headers.keySet()) {
                httpPost.setHeader(key, headers.get(key));
            }
        }

        StringEntity stringEntity = new StringEntity(content, charset);
        if (contentType != null) {
            stringEntity.setContentType(contentType);
        }
        httpPost.setEntity(stringEntity);
        return request(httpPost, getStringResponseHandler());
    }

//    curl \
//            -H "signature:lLTKGzWXT92gZHsrHwFrcQ0lr+k=" \
//            -H "Content-Type: application/json" \
//            -H "x-jdcloud-pin:amNsb3VkaWFhczI=" \
//            -X DELETE http://10.226.134.77:8080/v1/regions/cn-north-internet/instance/test2-pub-0730-5?jvessel_region=cn-north-1

    public String delete(String url, String content, String contentType, Map<String,String> headers) throws IOException {
        HttpDelete httpDelete = new HttpDelete(url);

        if(headers!=null) {
            for (String key : headers.keySet()) {
                httpDelete.setHeader(key, headers.get(key));
            }
        }

        StringEntity stringEntity = new StringEntity(content, charset);
        if (contentType != null) {
            stringEntity.setContentType(contentType);
        }

        // httpDelete.setEntity(stringEntity);
        return request(httpDelete, getStringResponseHandler());
    }

    public String deleteJSON(String url, String json, Map<String,String> headers) throws IOException {
        return patch(url, json, "application/json", headers);
    }

    public String patchJSON(String url, String json, Map<String,String> headers) throws IOException {
        return patch(url, json, "application/json", headers);
    }

    public String patch(String url, String content, String contentType, Map<String,String> headers) throws IOException {
        HttpPatch httpPatch = new HttpPatch(url);
        if(headers!=null) {
            for (String key : headers.keySet()) {
                httpPatch.setHeader(key, headers.get(key));
            }
        }

        StringEntity stringEntity = new StringEntity(content, charset);
        if (contentType != null) {
            stringEntity.setContentType(contentType);
        }
        httpPatch.setEntity(stringEntity);
        return request(httpPatch, getStringResponseHandler());
    }


    public String putJSON(String url, String json) throws IOException {
        return put(url, json, "application/json");
    }

    public String put(String url, String content, String contentType) throws IOException {
        HttpPut httpPut = new HttpPut(url);
        StringEntity stringEntity = new StringEntity(content, charset);
        if (contentType != null) {
            stringEntity.setContentType(contentType);
        }
        httpPut.setEntity(stringEntity);
        return request(httpPut, getStringResponseHandler());
    }

    public <T> T request(HttpUriRequest httpUriRequest, ResponseHandler<T> responseHandler) throws IOException {
        CloseableHttpClient httpclient = httpClientBuilder.build();

        T result = null;
        try {
            result = httpclient.execute(httpUriRequest, responseHandler);

        } finally {
            httpclient.close();
        }
        return result;
    }

    public List<NameValuePair> toNameValuePairs(Map<String, String> map) {
        List<NameValuePair> nameValuePairs = new ArrayList<NameValuePair>();
        for (Map.Entry<String, String> entry : map.entrySet()) {
            nameValuePairs.add(new BasicNameValuePair(entry.getKey(), entry.getValue()));
        }
        return nameValuePairs;
    }


    public ResponseHandler<String> getStringResponseHandler() {
        return new ResponseHandler<String>() {
            @Override
            public String handleResponse(final HttpResponse httpResponse) throws IOException {
                return EntityUtils.toString(httpResponse.getEntity(), HttpClientUtil.this.charset);
            }
        };
    }

    public void setTimeout(int timeout) {
        this.timeout = timeout;
    }

    public void setCharset(Charset charset) {
        this.charset = charset;
    }

}
