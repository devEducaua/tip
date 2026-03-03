package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"net"
	"path/filepath"
	"strings"
)

type ResponseCode int;
const (
	CodeSuccess ResponseCode = iota
	CodeRedirect
	CodeError
)

type Req struct {
	addr string
	path string
}

type Res struct {
	code ResponseCode
	message string
}

var PATH = "./serve";

func main() {
	argv := os.Args[1:];

	if len(argv) > 0 {
		PATH = argv[0];
	}

	listener, err := net.Listen("tcp", ":1979");
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not listen: %v\n", err);
		os.Exit(1);
	}

	for {
		conn, err := listener.Accept();
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: could not accept the connection: %v\n", err);
			os.Exit(1);
		}

		go handleConnection(conn);
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close();

	reader := bufio.NewReader(conn);
	rawReq, err := reader.ReadString('\n');
	if err != nil {
		if err == io.EOF {
			return;
		}
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err);
		return;
	}

	req := parseRequest(rawReq);
	fmt.Printf("%v -> %v\n", conn.RemoteAddr().String(), req.path);

	res := getResponse(req)
	fmt.Fprintf(conn, "%v\n%v\n", res.code, res.message);
}

func getResponse(req Req) Res {
	var res Res;

	files, err := readDirFromPath(PATH);
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not read the dir %v\n", err);
		return Res{}; 
	}

	found := files[req.path[1:]];
	if !found {
		res.message = fmt.Sprintf("file not found: %v\n", req.path)
		res.code = CodeError;
		return res;
	}

    dat, err := os.ReadFile(PATH+req.path);
	if err != nil {
		res.message = fmt.Sprintf("could not open the file: %v\n", req.path);
		res.code = CodeError;
		return res;
	}

	res.message = string(dat);
	res.code = CodeSuccess;
	return res;
}

func parseRequest(req string) Req {
	splitted := strings.SplitN(req, " ", 2);
	addr := splitted[0];
	path := strings.TrimSpace(splitted[1]);

	if strings.HasSuffix(path, "/") {
		path += "main.tip";
	}
	return Req{addr, path};
}

func readDirFromPath(root string) (map[string]bool, error) {
	m := make(map[string]bool);

	err := filepath.WalkDir(root, func (path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err;
		}

		relpath, err := filepath.Rel(root, path);
		if err != nil {
			return err;
		}
		if relpath == "." {
			return nil;
		}

		m[relpath] = true;
		return nil;
	})
	if err != nil {
		return nil, err;

	}
	return m, nil;
}
