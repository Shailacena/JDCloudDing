package util;

import request.RequestHeaders;
import request.RequestParameters;
import request.SendMessageParameters;

import java.beans.IntrospectionException;
import java.beans.Introspector;
import java.beans.PropertyDescriptor;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;


public class RequestSignSourceGenerator {
    private static final String FIELD_NAME_MESSAGE = "messages";
    private static final String FIELD_NAME_SIGNATURE = "signature";


    /**
     * 用户缓存PropertyDescriptor信息 <Class, Map<propertyName, PropertyDescriptor>>
     */
    private static ConcurrentMap<Class, ConcurrentMap<String, PropertyDescriptor>> propertyDescriptorMap = new ConcurrentHashMap<Class, ConcurrentMap<String, PropertyDescriptor>>();

    public static String getSignSource(final RequestParameters requestParameters, final String accessKey, final String datetime, final String token){
        RequestHeaders header= new RequestHeaders(accessKey, "", datetime, token);
        return getSignSource(header, requestParameters);
    }

    public static String getSignSource(final RequestParameters requestParameters, final String accessKey, final String datetime){
        RequestHeaders header= new RequestHeaders(accessKey, "", datetime, null);
        return getSignSource(header, requestParameters);
    }

    /**
     * 通过反射获得request的属性，并按属性名称排序拼接成签名的源串
     *
     * @return String
     */
    public static String getSignSource(final RequestHeaders requestHeaders, final RequestParameters requestParameters) {
        TreeMap<String, String> fieldValueMap = new TreeMap<String, String>();
        fieldValueMap.putAll(parse(requestHeaders));
        fieldValueMap.putAll(parse(requestParameters));
        return objectToString(fieldValueMap);
    }

    public static String getSignSource(final Object object) {
        TreeMap<String, String> fieldValueMap = new TreeMap<String, String>();
        fieldValueMap.putAll(parse(object));
        return objectToString(fieldValueMap);
    }

    private static HashMap<String, String> parse(Object object) {
        if (object == null) {
            return null;
        }

        ConcurrentMap<String, PropertyDescriptor> propertyDescriptors = propertyDescriptorMap.get(object.getClass());
        if (propertyDescriptors == null) {
            try {
                propertyDescriptors = new ConcurrentHashMap<String, PropertyDescriptor>(16);
                PropertyDescriptor[] descriptors = Introspector.getBeanInfo(object.getClass(), Object.class).getPropertyDescriptors();
                for (PropertyDescriptor descriptor : descriptors) {
                    propertyDescriptors.put(descriptor.getName(), descriptor);
                }
                propertyDescriptorMap.put(object.getClass(), propertyDescriptors);
            } catch (IntrospectionException e) {
                return null;
            }
        }

        HashMap<String, String> fieldValueMap = new HashMap<String, String>(propertyDescriptors.size());
        for (Map.Entry<String, PropertyDescriptor> entry : propertyDescriptors.entrySet()) {
            String fieldName = entry.getKey();
            PropertyDescriptor propertyDescriptor = entry.getValue();
            Method readMethod = propertyDescriptor.getReadMethod();
            readMethod.setAccessible(true);
            try {
                Object fieldValue = readMethod.invoke(object);
                if (fieldValue == null || FIELD_NAME_SIGNATURE.equals(fieldName)) {
                    continue;
                }
                if (propertyDescriptor.getPropertyType() == Map.class) {
                    for (Object objectKey : ((Map) fieldValue).keySet()) {
                        Object objectValue = ((Map) fieldValue).get(objectKey);
                        if (objectValue != null) {
                            fieldValueMap.put(objectToString(objectKey), objectToString(objectValue));
                        }
                    }
                } else if (propertyDescriptor.getPropertyType() == List.class) {
                    List<String> listValues = new ArrayList<String>();
                    for(Object value:((List)fieldValue)){
                        listValues.add(objectToString(value));
                    }

                    StringBuilder joinStr = new StringBuilder();
                    for(int i=0; i<listValues.size();i++)
                    {
                        joinStr.append(listValues.get(i));
                        if(i!=listValues.size()-1) {
                            joinStr.append(",");
                        }
                    }

                    fieldValueMap.put(fieldName, joinStr.toString());
                } else {
                    fieldValueMap.put(fieldName, objectToString(fieldValue));
                }
            } catch (IllegalAccessException | InvocationTargetException e) {
            }
        }

        return fieldValueMap;
    }

    private static String objectToString(Object object) {
        if (object.getClass() == SendMessageParameters.Message.class) {
            return CryptoUtils.md5(getSignSource(object));
        }
        if (object.getClass() == TreeMap.class) {
            StringBuilder stringBuilder = new StringBuilder();
            TreeMap<Object,Object> maps = (TreeMap<Object, Object>) object;
            for(Object key:maps.keySet()) // default是升序
            {
                stringBuilder.append(objectToString(key)).append("=").append(objectToString(maps.get(key))).append("&");
            }

            return stringBuilder.substring(0, Math.max(0, stringBuilder.length() - 1));
            //return stringBuilder.toString();
        }
        return object.toString();
    }

}