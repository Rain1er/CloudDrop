error_reporting(0);
header('Content-Type: text/html; charset=UTF-8');



function main($path = "",$newName="")
{
	
	$dlen = strlen($path);
	$cdlen = $dlen - strrpos($path,"/");
	$pdn = substr($path,0,-$cdlen);
	$res=rename($path,$pdn."/".$newName);
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