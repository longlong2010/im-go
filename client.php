<?php
$sock = fsockopen('tcp://127.0.0.1:1987');
$data = json_encode(array(
    'id' => 1234,
    'data' => "hello",
));
fprintf($sock, "\x87%s%s", pack("I", strlen($data)), $data);
sleep(10);
