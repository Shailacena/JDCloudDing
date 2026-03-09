package response;


public enum ResponseStatus {

    OK(200, "OK"),
    /**
     * 参数无效
     */
    INVALID_ARGUMENT(400, "INVALID_ARGUMENT"),
    /**
     * 预置条件不满足
     */
    FAILED_PRECONDITION(400,"FAILED_PRECONDITION"),
    /**
     * 折扣无效
     */
    OUT_OF_RANGE(400, "OUT_OF_RANGE"),
    /**
     *  验证失败
     */
    UNAUTHENTICATED(401,"UNAUTHENTICATED"),
    /**
     * 没有权限
     */
    PERMISSION_DENIED(403,"PERMISSION_DENIED"),
    /**
     * 资源没有找到
     */
    NOT_FOUND(404, "NOT_FOUND"),

    DELETED(404, "DELETED"),
    /**
     * 冲突
     */
    CONFLICT(409, "CONFLICT"),
    /**
     * 资源状态
     */
    STATE(409,"STATE"),

    /**
     * 操作终止
     */
    ABORTED(409, "ABORTED"),
    /**
     *对象已存在
     */
    ALREADY_EXISTS(409, "ALREADY_EXISTS"),
    /*
     * 资源处理中
     */
    IN_PROCESSING(409, "IN PROCESSING"),
    /**
     *配额不足
     */
    QUOTA_EXCEEDED(429,"QUOTA_EXCEEDED"),
    /**
     * 请求过频繁
     */
    RATE_LIMIT(429,"RATE_LIMIT"),
    /**
     * 取消操作
     */
    CANCELLED(499,"CANCELLED"),
    /**
     * 内部错误
     */
    INTERNAL_ERROR(500, "INTERNAL"),

    DATA_ERROR(500, "DATA"),
    /**
     * 未知错误
     */
    UNKNOWN(500,"UNKNOWN"),
    /**
     * 操作未实现
     */
    NOT_IMPLEMENTED(500,"NOT_IMPLEMENTED"),


    UNAVAILABLE(503,"UNAVAILABLE"),
    /**
     * 超时
     */
    DEADLINE_EXCEEDED(504,"DEADLINE_EXCEEDED")

    ;

    /**
     * 响应编号
     */
    private int code;
    /**
     * 响应信息字符串
     */
    private String status;

    ResponseStatus(int code, String status) {
        this.code = code;
        this.status = status;
    }

    public int getCode() {
        return code;
    }

    public String getStatus() {
        return status;
    }
}
