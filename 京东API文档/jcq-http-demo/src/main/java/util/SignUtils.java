package util;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import javax.xml.bind.DatatypeConverter;

public class SignUtils {
    private static final String HMAC_SHA1 = "HmacSHA1";
    private static final String UTF8 = "UTF-8";

    public SignUtils() {
    }

    public static String signWithHmacSha1(String source, String key) throws Exception {
        if (!StringUtils.isEmpty(source) && !StringUtils.isEmpty(key)) {
            Mac mac = Mac.getInstance("HmacSHA1");
            mac.init(new SecretKeySpec(key.getBytes("UTF-8"), "HmacSHA1"));
            return DatatypeConverter.printBase64Binary(mac.doFinal(source.getBytes("UTF-8")));
        } else {
            return null;
        }
    }
}
