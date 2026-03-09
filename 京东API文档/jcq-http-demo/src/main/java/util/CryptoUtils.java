package util;

import java.math.BigInteger;
import java.nio.charset.Charset;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;


public class CryptoUtils {

    private static final String ALGORITHM_MD5 = "md5";

    static String md5(final String content) {
        if (content == null) {
            return null;
        }
        String result = null;
        try {
            MessageDigest messageDigest = MessageDigest.getInstance(ALGORITHM_MD5);
            messageDigest.update(content.getBytes(Charset.forName("utf-8")));
            BigInteger resultInteger = new BigInteger(1, messageDigest.digest());
            result = String.format("%032x", resultInteger);
        } catch (NoSuchAlgorithmException e) {
        }
        return result;
    }
}
