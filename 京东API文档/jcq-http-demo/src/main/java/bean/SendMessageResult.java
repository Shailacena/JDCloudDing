package bean;


import java.util.List;

public class SendMessageResult {
    private List<String> messageIds;

    public List<String> getMessageIds() {
        return messageIds;
    }

    public void setMessageIds(List<String> messageIds) {
        this.messageIds = messageIds;
    }

    @Override
    public String toString() {
        return "SendMessageResult{" +
                "messageIds=" + messageIds +
                '}';
    }
}
