package demo;

import request.AckMessageParametersV2;
import request.GetMessageParameters;
import request.SendMessageParameters;
import response.HttpProxyGetMessagesV2ResultResponse;
import response.HttpProxyResultResponse;
import response.HttpProxySendMessagesResultResponse;
import util.JCQHttpUtil;
import util.StringUtils;

import java.util.ArrayList;
import java.util.Calendar;
import java.util.HashMap;
import java.util.TimeZone;

public class JCQHttpProcessor {


    /**
     * 用户AccessKey
     */
    public static final String AK = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA";
    /**
     * 用户secretKey
     */
    public static final String Sk = "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB";

    /**
     * 元数据服务器地址
     */
    public static final String ENDPOINT = "$HTTP_ENDPOINTS";
    /**
     * topic名称
     */
    public static final String TOPIC_NAME = "topic_001";


    /**
     * 消费组
     */
    public static final String CONSUMER_GROUP_ID = "CG_001";



    public void sendMessage() throws Exception {
        SendMessageParameters sendMsgParas = new SendMessageParameters();
        SendMessageParameters.Message msg1 = sendMsgParas.new Message();
        HashMap<String, String> map1 = new HashMap<>();
        map1.put("key1", "value1");
        msg1.setBody("this is https");
        msg1.setProperties(map1);
        ArrayList<SendMessageParameters.Message> msgs = new ArrayList<>();
        msgs.add(msg1);
        sendMsgParas.setMessages(msgs);
        sendMsgParas.setTopic(TOPIC_NAME);
        sendMsgParas.setType("NORMAL");
        sendMsgParas.setDateTime(getUtcNow());
        // bHttps参数说明: true（https方式）, false(http方式)
        HttpProxySendMessagesResultResponse response = JCQHttpUtil.sendMessages(JCQHttpProcessor.AK, JCQHttpProcessor.Sk,  sendMsgParas, ENDPOINT, null, false);
        System.out.println("response is " + response);
    }

    private void pullMessageAutoAck() throws Exception {
        GetMessageParameters getMsgParas = new GetMessageParameters();
        getMsgParas.setTopic(TOPIC_NAME);
        getMsgParas.setConsumerGroupId(CONSUMER_GROUP_ID);
        getMsgParas.setSize(1); // 一次最多拉取消息条数，全局顺序消息推荐设置为1
        getMsgParas.setConsumerId(null);
        getMsgParas.setConsumeFromWhere(null);
        getMsgParas.setFilterExpressionType(null);
        getMsgParas.setFilterExpression(null);
        getMsgParas.setAck(true); //这里有区别
        getMsgParas.setDateTime(getUtcNow());
        // bHttps参数说明: true（https方式）, false(http方式)
        HttpProxyGetMessagesV2ResultResponse response = JCQHttpUtil.getMessages(JCQHttpProcessor.AK, JCQHttpProcessor.Sk, getMsgParas, ENDPOINT, false, null, false);
        System.out.println("response is " + response);
    }

    private void pullMessageManualAck() throws Exception{
        GetMessageParameters getMsgParas = new GetMessageParameters();
        getMsgParas.setTopic(TOPIC_NAME);
        getMsgParas.setConsumerGroupId(CONSUMER_GROUP_ID);
        getMsgParas.setSize(1); // 一次最多拉取消息条数，全局顺序消息推荐设置为1
        getMsgParas.setConsumerId(null);
        getMsgParas.setConsumeFromWhere(null);
        getMsgParas.setFilterExpressionType(null);
        getMsgParas.setFilterExpression(null);
        getMsgParas.setAck(false); //这里有区别
        getMsgParas.setDateTime(getUtcNow());
        // bHttps参数说明: true（https方式）, false(http方式)
        HttpProxyGetMessagesV2ResultResponse response = JCQHttpUtil.getMessages(JCQHttpProcessor.AK, JCQHttpProcessor.Sk, getMsgParas, ENDPOINT, false, null, false);
        System.out.println("manual ack pull response is " + response);
        if(null != response.getError() || null == response.getResult() || StringUtils.isEmpty(response.getResult().getAckIndex())){
            return;
        }

        String ackIndex = response.getResult().getAckIndex();
        if(ackIndex != null){
            AckMessageParametersV2 ackParas = new AckMessageParametersV2();
            ackParas.setTopic(TOPIC_NAME);
            ackParas.setDateTime(getUtcNow());
            ackParas.setConsumerGroupId(CONSUMER_GROUP_ID);
            ackParas.setAckAction("SUCCESS");
            ackParas.setAckIndex(ackIndex);
            // bHttps参数说明: true（https方式）, false(http方式)
            HttpProxyResultResponse ackResponse = JCQHttpUtil.ack(JCQHttpProcessor.AK, JCQHttpProcessor.Sk, ackParas, ENDPOINT, null, false);

            System.out.println("ack response is " + ackResponse);
        }


    }

    public static String getUtcNow(){
        java.text.SimpleDateFormat simpleDateFormat =new java.text.SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'");
        simpleDateFormat.setTimeZone(TimeZone.getTimeZone("UTC"));
        return simpleDateFormat.format(Calendar.getInstance().getTime());
    }

    public static void main(String[] args) throws Exception {
        JCQHttpProcessor processor = new JCQHttpProcessor();
        processor.sendMessage();
        processor.sendMessage();
        processor.pullMessageAutoAck();
        processor.pullMessageManualAck();
    }

}
