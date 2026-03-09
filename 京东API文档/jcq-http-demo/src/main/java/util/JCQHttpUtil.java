package util;

import com.alibaba.fastjson.JSON;
import request.AckMessageParametersV2;
import request.GetMessageParameters;
import request.SendMessageParameters;
import response.HttpProxyGetMessagesV2ResultResponse;
import response.HttpProxyResultResponse;
import response.HttpProxySendMessagesResultResponse;

public class JCQHttpUtil {

    public static HttpProxySendMessagesResultResponse sendMessages(String accessKey, String secretKey, SendMessageParameters sendMsgs, String endPoint, String accessToken, boolean bHttps) throws Exception{
//        String dateTime = getUtcNow();
        String sourceSign = RequestSignSourceGenerator.getSignSource(sendMsgs, accessKey, sendMsgs.getDateTime(), accessToken);
        String signature = SignUtils.signWithHmacSha1(sourceSign, secretKey);
        String jsonStr = JSON.toJSONString(sendMsgs);
        String url = String.format("%s/%s/messages", GetWebSiteRootUrl(endPoint, bHttps), "v2");
        String returnObj = HttpClientHelper.getInstance().doHttpProxyPost(url,jsonStr, accessKey, signature, sendMsgs.getDateTime(), accessToken);
        return JSON.parseObject(returnObj, HttpProxySendMessagesResultResponse.class);
    }


    public static HttpProxyGetMessagesV2ResultResponse getMessages(
            String accessKey,
            String secretKey,
            GetMessageParameters getMsgParas,
            String endPoint,
            boolean removeIfEmpty,
            String accessToken,
            boolean bHttps) throws Exception {

        String sourceSign = RequestSignSourceGenerator.getSignSource(getMsgParas, accessKey,  getMsgParas.getDateTime(), accessToken);
        String signature = SignUtils.signWithHmacSha1(sourceSign, secretKey);

        String url = String.format("%s/%s/messages/?", GetWebSiteRootUrl(endPoint, bHttps), "v2");
        StringBuilder sb = new StringBuilder();
        sb.append(url);
        sb.append(String.format("topic=%s&", getMsgParas.getTopic()));
        sb.append(String.format("consumerGroupId=%s&", getMsgParas.getConsumerGroupId()));
        if (getMsgParas.getSize() != null) {
            sb.append(String.format("size=%s&", getMsgParas.getSize()));
        }
        if(removeIfEmpty)
        {
            sb.append(String.format("consumerId=%s&", getMsgParas.getConsumerId()));
            sb.append(String.format("consumeFromWhere=%s&", getMsgParas.getConsumeFromWhere()));
            sb.append(String.format("filterExpressionType=%s&", getMsgParas.getFilterExpressionType()));
            sb.append(String.format("filterExpression=%s&", getMsgParas.getFilterExpression()));
        }
        else {
            if (getMsgParas.getConsumerId() != null && !getMsgParas.getConsumerId().isEmpty()) {
                sb.append(String.format("consumerId=%s&", getMsgParas.getConsumerId()));
            }
            if (getMsgParas.getConsumeFromWhere() != null && !getMsgParas.getConsumeFromWhere().isEmpty()) {
                sb.append(String.format("consumeFromWhere=%s&", getMsgParas.getConsumeFromWhere()));
            }
            if (getMsgParas.getFilterExpressionType() != null && !getMsgParas.getFilterExpressionType().isEmpty()) {
                sb.append(String.format("filterExpressionType=%s&", getMsgParas.getFilterExpressionType()));
            }
            if (getMsgParas.getFilterExpression() != null && !getMsgParas.getFilterExpression().isEmpty()) {
                sb.append(String.format("filterExpression=%s&", getMsgParas.getFilterExpression()));
            }
        }
        sb.append(String.format("ack=%s", getMsgParas.getAck().toString()));
        String jsonString = HttpClientHelper.getInstance().doHttpProxyGet(sb.toString(), accessKey, getMsgParas.getDateTime(), signature, accessToken);
        return JSON.parseObject(jsonString, HttpProxyGetMessagesV2ResultResponse.class);
    }

    public static HttpProxyResultResponse ack(String accessKey, String secretKey, AckMessageParametersV2 ackParas, String endPoint, String accessToken, boolean bHttps) throws Exception{
        String sourceSign = RequestSignSourceGenerator.getSignSource(ackParas, accessKey, ackParas.getDateTime(), accessToken);
        String signature = SignUtils.signWithHmacSha1(sourceSign, secretKey);
        String jsonStr = JSON.toJSONString(ackParas);
        String url = String.format("%s/%s/ack", GetWebSiteRootUrl(endPoint, bHttps), "v2");
        String returnObj = HttpClientHelper.getInstance().doHttpProxyPost(url,jsonStr, accessKey, signature, ackParas.getDateTime(), accessToken);
        return JSON.parseObject(returnObj, HttpProxyResultResponse.class);
    }

    public static String GetWebSiteRootUrl(String httpProxy, boolean bHttps){
        if(bHttps){
            return String.format("https://%s", httpProxy);
        }
        return String.format("http://%s", httpProxy);
    }

}
