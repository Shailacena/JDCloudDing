package response;


public class HttpProxyResultResponse<T> extends HttpProxyResponse {
    private T result;
    private ErrorResponse error;
    public ErrorResponse getError() {
        return error;
    }

    public void setError(ErrorResponse error) {
        this.error = error;
    }

    public HttpProxyResultResponse() {}

    public HttpProxyResultResponse(T result) {this.result = result;}

    public HttpProxyResultResponse(String requestId, T result) {
        super(requestId);
        this.result = result;
    }

    public T getResult() {
        return result;
    }

    public void setResult(T result) {
        this.result = result;
    }

    @Override
    public String toString() {
        return "HttpProxyResultResponse{" +
                "result=" + result +
                ", error=" + error +
                '}';
    }
}
