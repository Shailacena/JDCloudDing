package request;

import com.alibaba.fastjson.JSON;
import util.StringUtils;

import java.util.List;
import java.util.Map;


public class SendMessageParameters extends RequestParameters {
    private static final int MESSAGES_MAX_SIZE = 32;

    private String topic;
    private String type;
    private List<Message> messages;

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

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public List<Message> getMessages() {
        return messages;
    }

    public void setMessages(List<Message> messages) {
        this.messages = messages;
    }

    public class Message {
        private String body;
        private Integer delaySeconds;
        private String tag;
        private Map<String, String> properties;

        public boolean validate() {
            return isBodyValid() && isDelaySecondsValid() && isTagValid();
        }

        private boolean isBodyValid() {
            return !StringUtils.isEmpty(getBody());
        }

        private boolean isDelaySecondsValid() {
            return getDelaySeconds() == null || getDelaySeconds() >= 0;
        }

        private boolean isTagValid() {
            return getTag() == null || getTag().length() <= 64;
        }

        public String getBody() {
            return body;
        }

        public void setBody(String body) {
            this.body = body;
        }

        public Integer getDelaySeconds() {
            return delaySeconds;
        }

        public void setDelaySeconds(Integer delaySeconds) {
            this.delaySeconds = delaySeconds;
        }

        public String getTag() {
            return tag;
        }

        public void setTag(String tag) {
            this.tag = tag;
        }

        public Map<String, String> getProperties() {
            return properties;
        }

        public void setProperties(Map<String, String> properties) {
            this.properties = properties;
        }
    }
}
