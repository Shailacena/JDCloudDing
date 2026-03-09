package response;

import java.io.Serializable;


public class HttpProxyResponse implements Serializable {
    private String requestId;

    public HttpProxyResponse() {}

    public HttpProxyResponse(final String requestId) {
        this.requestId = requestId;
    }

    public String getRequestId() {
        return requestId;
    }

    public void setRequestId(String requestId) {
        this.requestId = requestId;
    }
}
