package util;


import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.util.Properties;

public class ConfigHelper {
    private static Properties properties = null;

    private static void load() {
        if (properties == null) {
            properties = new Properties();
            try {
                InputStream inputStream = ConfigHelper.class.getClassLoader().getResourceAsStream("jcqTestConfig.properties");
                BufferedReader bf = new BufferedReader(new InputStreamReader(inputStream, "UTF-8"));
                properties.load(bf);
            } catch (IOException e) {
            }
        }
    }

    static {
        load();
    }

    public static String getAdminServerAddress() {
        return getValue("AdminServerAddress");
    }

    public static String getAclAddress() {
        return getValue("AclAddress");
    }

    public static String getEtcdEndPoint() {
        return getValue("EtcdEndPoint");
    }

//    public static String getHttpProxyEndPoint(){
//        return getValue("HttpProxyEndPoint");
//    }

    public static String getHttpProxyVersion() {
        return getValue("HttpProxyVersion");
    }

    public static String getOpenAPIAddress() {
        return getValue("OpenAPIAddress");
    }

    public static String getRegionId() {
        return getValue("RegionId");
    }

    public static String getExistNormalTopic(){
        return getValue("ExistNormalTopic");
    }

    public static String getExistNormalCGID(){
        return getValue("ExistNormalCGId");
    }

    public static String getExistNormalCGID_Sub(){
        return getValue("ExistNormalCGId_Sub");
    }

    public static String getExistNormalCGID_Auth(){
        return getValue("ExistNormalCGId_Auth");
    }

    public static String getExistOrderTopic(){
        return getValue("ExistOrderTopic");
    }

    public static String getExistOrderCGID(){
        return getValue("ExistOrderCGId");
    }

    public static String getExistOrderCGID_Sub(){
        return getValue("ExistOrderCGId_Sub");
    }

    public static String getExistOrderCGID_Auth(){
        return getValue("ExistOrderCGId_Auth");
    }

//    public static String getManagerAddress() {
//        return getValue("ManagerAddress");
//    }

    public static String getEtcdServerAddress() {
        return getValue("EtcdServerAddress");
    }

    public static String getUserId() {
        return getValue("UserId");
    }

    public static String getUserPin() {
        return getValue("UserPin");
    }

    public static String getAuthorizedUserPin() {
        return getValue("AuthorizedUserPin");
    }

    public static String getMySqlDBConnectionStr() {
        return getValue("MySqlDBConnection");
    }

    public static String getBrokerAddress() {
        return getValue("BrokerAddress");
    }

    public static String getAuthorizedUserId() {
        return getValue("AuthorizedUserId");
    }

    public static String getSubAccountUserId() {
        return getValue("SubUserId");
    }

    public static String getGatewayUseProxy() {
        return getValue("GatewayUseProxy");
    }

    public static String getGatewayProxyEndpoint() {
        return getValue("GatewayProxyEndpoint");
    }

    public static String getRealGatewayEndpoint() {
        return getValue("RealGatewayEndpoint");
    }
    public static String getGatewayAddree() {
        return getValue("GatewayAddree");
    }


    private static void checkNullOrEmpty(String key, String value) {
        if (value == null || value.isEmpty()) {
        }
    }

    private static String getValue(String key) {
        String value = System.getenv(key);
        if(value == null || value.isEmpty()){
            value = properties.getProperty(key);
            checkNullOrEmpty(key, value);
        }

        return value;
    }
}
