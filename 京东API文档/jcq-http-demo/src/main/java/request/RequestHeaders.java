package request;

public class RequestHeaders {

    private String accessKey;
    private String signature;
    private String dateTime;
    private String token;

    public RequestHeaders() {}

    public RequestHeaders(final String accessKey, final String signature, final String dateTime, final String token) {
        this.accessKey = accessKey;
        this.signature = signature;
        this.dateTime = dateTime;
        this.token = token;

    }

    public String getAccessKey() {
        return accessKey;
    }

    public void setAccessKey(String accessKey) {
        this.accessKey = accessKey;
    }

    public String getSignature() {
        return signature;
    }

    public void setSignature(String signature) {
        this.signature = signature;
    }

    public String getDateTime() {
        return dateTime;
    }

    public void setDateTime(String dateTime) {
        this.dateTime = dateTime;
    }

    public String getToken() {
        return token;
    }

    public void setToken(String token) {
        this.token = token;
    }
}
