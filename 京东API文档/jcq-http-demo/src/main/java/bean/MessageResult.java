package bean;

import com.alibaba.fastjson.JSON;

import java.util.List;
import java.util.Map;


public class MessageResult {
    private String topicName;
    private Long ackIndex;
    private List<Message> messages;

    @Override
    public String toString() {
        return JSON.toJSONString(this);
    }

    public MessageResult() {}

    public MessageResult(String topicName, Long ackIndex, List<Message> messages) {
        this.topicName = topicName;
        this.ackIndex = ackIndex;
        this.messages = messages;
    }

    public String getTopicName() {
        return topicName;
    }

    public void setTopicName(String topicName) {
        this.topicName = topicName;
    }

    public Long getAckIndex() {
        return ackIndex;
    }

    public void setAckIndex(Long ackIndex) {
        this.ackIndex = ackIndex;
    }

    public List<Message> getMessages() {
        return messages;
    }

    public void setMessages(List<Message> messages) {
        this.messages = messages;
    }

    public static class Message {
        private String messageId;
        private String messageBody;
        private Map<String, String> properties;

        public Message() {}

        public Message(final String messageId, final String messageBody, final Map<String, String> properties) {
            this.messageId = messageId;
            this.messageBody = messageBody;
            this.properties = properties;
        }

        @Override
        public String toString() {
            return JSON.toJSONString(this);
        }

        public String getMessageId() {
            return messageId;
        }

        public void setMessageId(String messageId) {
            this.messageId = messageId;
        }

        public String getMessageBody() {
            return messageBody;
        }

        public void setMessageBody(String messageBody) {
            this.messageBody = messageBody;
        }

        public Map<String, String> getProperties() {
            return properties;
        }

        public void setProperties(Map<String, String> properties) {
            this.properties = properties;
        }
    }
}
