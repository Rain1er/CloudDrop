error_reporting(0);
header('Content-Type: text/html; charset=UTF-8');

function main($srcPath = "",$toPath="")
{
	$zip = new ZipArchive();
	$openRes = $zip->open($srcPath);
	$msg="Success";
	if ($openRes === TRUE) {
	  $zip->extractTo($toPath);
	  $zip->close();
	}
	else{
		$msg="Failed";
	}
	echo $msg;
}
