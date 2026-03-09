package response;


public class HttpProxyErrorResponse extends HttpProxyResponse {
    private ErrorResponse error;

    public HttpProxyErrorResponse() {}

    public HttpProxyErrorResponse(String requestId, int code, String status, String message) {
        super(requestId);
        this.error = new ErrorResponse(code, status, message);
    }

    public ErrorResponse getError() {
        return error;
    }

    public void setError(ErrorResponse error) {
        this.error = error;
    }
}
