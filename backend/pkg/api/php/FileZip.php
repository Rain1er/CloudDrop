error_reporting(0);
header('Content-Type: text/html; charset=UTF-8');



function main($srcPath = "",$toPath="")
{
	$fileList =explode(",",$srcPath);
	$zip = new ZipArchive();
	$zip->open($toPath,ZipArchive::CREATE);
	foreach($fileList as $file){
		if(is_file($file)){
			$zip->addFile($file,basename($file));
		}
		else{
			$isend=(substr($file,-1) === "/");
			if($isend){
				$fdic=substr($file,1);;
			}
			else{
				$fdic=$file;
			}
			$cdic=substr($fdic,strripos($fdic,"/")+1);
			addFileToZip($cdic,$file,$file,$zip);
		}
	}
	$zip->close();
	echo "Success";
}

function addFileToZip($cdic,$root,$path,$zip){
    $handler=opendir($path);
    while(($filename=readdir($handler))!==false){
        if($filename!= "."&&$filename!= ".."){
			$cpath=$path."/".$filename;
			$capath=str_replace($root,"",$path."/".$filename);
			$capath=$cdic."/".substr($capath,1);
            if(is_dir($cpath)){
				$zip->addEmptyDir($capath);
                addFileToZip($cdic,$root,$path."/".$filename, $zip);
            }else{
                $zip->addFile($path."/".$filename,$capath);
            }
        }
    }
    @closedir($path);
}
