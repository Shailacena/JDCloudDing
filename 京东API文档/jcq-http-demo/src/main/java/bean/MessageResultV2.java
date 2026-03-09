package bean;

import com.alibaba.fastjson.JSON;

import java.util.List;

public class MessageResultV2 {
    private String topicName;
    private String ackIndex;
    private List<MessageResult.Message> messages;

    public MessageResultV2() {
    }

    public MessageResultV2(String topicName, String ackIndex, List<MessageResult.Message> messages) {
        this.topicName = topicName;
        this.ackIndex = ackIndex;
        this.messages = messages;
    }

    @Override
    public String toString() {
        return JSON.toJSONString(this);
    }

    public String getTopicName() {
        return topicName;
    }

    public void setTopicName(String topicName) {
        this.topicName = topicName;
    }

    public String getAckIndex() {
        return ackIndex;
    }

    public void setAckIndex(String ackIndex) {
        this.ackIndex = ackIndex;
    }

    public List<MessageResult.Message> getMessages() {
        return messages;
    }

    public void setMessages(List<MessageResult.Message> messages) {
        this.messages = messages;
    }


}
