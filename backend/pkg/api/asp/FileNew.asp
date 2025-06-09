Function GetFso()
	Dim Fso,Key
	Key="Scripting.FileSystemObject"
	Set Fso=server.CreateObject(Key)
	Set GetFso=Fso
End Function

Sub FileNew(path)
  on error resume next
  result="Failed"
  Const adTypeBinary = 1
  Const adSaveCreateOverWrite = 2
  Dim BinaryStream
  Set BinaryStream = CreateObject("ADODB.Stream")
  BinaryStream.Type = adTypeBinary
  BinaryStream.Open
  BinaryStream.SaveToFile path, adSaveCreateOverWrite
  If Err Then 
	result=Err.Description
	Else
	result="Success" 
  End If
Response.write(result)
End Sub

Sub main(path)
	call FileNew(path)
End Sub