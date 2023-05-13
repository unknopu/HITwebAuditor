<?PHP 
    $file = $_GET["file"];
    $handle = fopen($file, 'r');
    $poem = fread($handle, 1);
    fclose($handle);
    echo $poem;
?>


