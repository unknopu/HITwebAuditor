<?php
    $user_name = $_GET['user_name'];
    $sql = "SELECT id FROM users WHERE name='$user_name'";
    $query_result = mysql_query($sql);
    var_dump($query_result);
?>


