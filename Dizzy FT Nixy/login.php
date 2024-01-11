<?php
    error_reporting(0);
   
    try
    {
        $user = $_POST["user"];
        $pass = $_POST["pass"];
        
        $path = "../data/users.json";
        
        if(!file_exists($path) or filesize($path)==false)
        {
               $mfile = fopen($path, "w");
            fwrite($mfile,"{}");
            $content = "{}";
            fclose($mfile);
        }
        $file = fopen($path, "r+");
        $content = fread($file,filesize($path));
        
        fclose($file);
        
        $data = json_decode($content,true);
        if(array_key_exists($user,$data))
        {
            $realPass = $data[$user];
            if($realPass == $pass)
            {
                //Correct password
                showMsg("Anda berhasil Login","Hallo,".$user."");
            }    
            else
            {
                //Incorrect password
                showMsg("Login Gagal","Password Salah.");
            }
        }
        else
        {
            //The account is not existed in users.json
            showMsg("gagal login","Username/password salah atau anda belum<a href='../reg.html'>Register</a> now.");
        }
    }
    catch(Exception $e)
    {
        showMsg("Error","Error message:".$e->getMessage().".");
    }
    
    function showMsg($title,$msg)
    {
        $file = fopen("../alert.html","r");
        echo (str_replace("MSG",$msg,str_replace("TITLE",$title,fread($file,filesize("../alert.html")))));
        fclose($file);
    }
?>