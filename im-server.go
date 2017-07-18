package main;

import "net";
import "encoding/binary";
import "encoding/json";
import "fmt";
import "io";

type Message struct {
    Id int64 `json:"id"`
    Data string `json:"data"`
}

func handle_error(err error) {
    if err != nil {
        panic(err);
    }
}

func main() {
    l, err := net.Listen("tcp", "127.0.0.1:1987");
    handle_error(err);
    defer l.Close();

    for {
        conn, err := l.Accept();
        handle_error(err);
        go func(conn net.Conn) {
            defer conn.Close();
            data := "";
            for {
                buf := make([]byte, 1024);
                n, err := conn.Read(buf);
                if err != nil {
                    switch err {
                        case io.EOF:
                            return;
                        default:
                            handle_error(err);
                    }
                }
                data += string(buf[0 : n]);
                k := -1;
                for i, _ := range data {
                    if data[i] == '\x87' {
                        k = i + 1;
                        break;
                    }
                }
                if k > 0 {
                    l1 := int(binary.LittleEndian.Uint32([]byte(data[k : k + 4])));
                    l2 := len(data);
                    if l2 - k - 4 >= l1 {
                        var m Message;
                        json.Unmarshal([]byte(data[k + 4 : k + 4 + l1]), &m);
                        data = data[k + 4 + l1 : l2];
                        fmt.Println(m.Data);
                    }
                }
            }
        } (conn);
    }
}
