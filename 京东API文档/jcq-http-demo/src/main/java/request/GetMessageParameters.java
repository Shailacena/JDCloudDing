package request;

import com.alibaba.fastjson.JSON;

public class GetMessageParameters extends RequestParameters {
    private String topic;
    private String consumerGroupId;
    private Integer size;
    private String consumerId;
    private String consumeFromWhere;
    private String filterExpressionType;
    private String filterExpression;
    private Boolean ack;

    public GetMessageParameters() {}

    public GetMessageParameters(final String topic, final String consumerGroupId, final Integer size,
                                final String consumerId, final String consumeFromWhere,
                                final String filterExpressionType, final String filterExpression, final Boolean ack) {
        this.topic = topic;
        this.consumerGroupId = consumerGroupId;
        this.size = size;
        this.consumerId = consumerId;
        this.consumeFromWhere = consumeFromWhere;
        this.filterExpressionType = filterExpressionType;
        this.filterExpression = filterExpression;
        this.ack = ack;
    }

    @Override
    public String toString() {
        return JSON.toJSONString(this);
    }

    public String getTopic() {
        return topic;
    }

    public void setTopic(String topic) {
        this.topic = topic;
    }

    public String getConsumerGroupId() {
        return consumerGroupId;
    }

    public void setConsumerGroupId(String consumerGroupId) {
        this.consumerGroupId = consumerGroupId;
    }

    public Integer getSize() {
        return size;
    }

    public void setSize(Integer size) {
        this.size = size;
    }

    public String getConsumerId() {
        return consumerId;
    }

    public void setConsumerId(String consumerId) {
        this.consumerId = consumerId;
    }

    public String getConsumeFromWhere() {
        return consumeFromWhere;
    }

    public void setConsumeFromWhere(String consumeFromWhere) {
        this.consumeFromWhere = consumeFromWhere;
    }

    public String getFilterExpressionType() {
        return filterExpressionType;
    }

    public void setFilterExpressionType(String filterExpressionType) {
        this.filterExpressionType = filterExpressionType;
    }

    public String getFilterExpression() {
        return filterExpression;
    }

    public void setFilterExpression(String filterExpression) {
        this.filterExpression = filterExpression;
    }

    public Boolean getAck() {
        return ack;
    }

    public void setAck(Boolean ack) {
        this.ack = ack;
    }
}
