package request;

import com.alibaba.fastjson.JSON;

public class AckMessageParameters extends RequestParameters {
    private String topic;
    private String consumerGroupId;
    private String ackAction;
    private Long ackIndex;



    private boolean isAckIndexValid() {
        return getAckIndex() != null && getAckIndex() >= 0;
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

    public String getAckAction() {
        return ackAction;
    }

    public void setAckAction(String ackAction) {
        this.ackAction = ackAction;
    }

    public Long getAckIndex() {
        return ackIndex;
    }

    public void setAckIndex(Long ackIndex) {
        this.ackIndex = ackIndex;
    }
}
