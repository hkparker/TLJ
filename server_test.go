package tlj

import (
	"testing"
	"net"
	"reflect"
	"time"
	//"encoding/json"
	"fmt"
)

func TagSocketAll(socket *net.Conn, server *Server) {
    server.Tags[socket] = append(server.Tags[socket], "all")
    server.Sockets["all"] = append(server.Sockets["all"], socket)
}

func TestServerIsCorrectType(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		t.Errorf("could not start test server on localhost:5000")
	}
	defer listener.Close()
	type_store := NewTypeStore()
	server := NewServer(listener, TagSocketAll, &type_store)
	if reflect.TypeOf(server) != reflect.TypeOf(Server{}) {
		t.Errorf("return value of NewServer() != tlj.Server")
	} 
}

func TestServerCanReceiveAndTagConnection(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		t.Errorf("could not start test server on localhost:5000")
	}
	defer listener.Close()
	
	type_store := NewTypeStore()
	server := NewServer(listener, TagSocketAll, &type_store)
	
	client_socket, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		t.Errorf("could not connect test client to localhost:5000")
	}
	defer client_socket.Close()
	
	time.Sleep(5 * time.Millisecond)		// wait for server to process incoming connection
	server_conns := server.Sockets["all"]
	if len(server_conns) != 1 {
		t.Errorf("socket did not get tagged as all")
	}
	if server.Tags[server_conns[0]][0] != "all" {
		t.Errorf("socket did not get tagged as all")
	}
	
	t.Errorf(fmt.Sprintf("%s", <- server.FailedServer))
}

func TestCanUseClientSocketInServer(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		t.Errorf("could not start test server on localhost:5000")
	}
	defer listener.Close()

	client_socket, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		t.Errorf("could not connect test client to localhost:5000")
	}
	defer client_socket.Close()

	other_listener, err := net.Listen("tcp", "localhost:5001")
	if err != nil {
		t.Errorf("could not start test server on localhost:5001")
	}
	defer other_listener.Close()
	
	type_store := NewTypeStore()
	server := NewServer(other_listener, TagSocketAll, &type_store)
	server.Insert(client_socket)

	server_conns := server.Sockets["all"]
	if len(server_conns) != 1 {
		t.Errorf("socket did not get tagged as all")
	}
	if server.Tags[server_conns[0]][0] != "all" {
		t.Errorf("socket did not get tagged as all")
	}
}
/*
func TestServerCanRunAcceptCallbacks(t *testing.T) {
	server_filename := "server_test-ipc-" + uuid.NewV4().String()
	listener, err := net.Listen("unix", server_filename)
	if err != nil {
		t.Errorf("could not start unix server")
	}
	defer listener.Close()
	defer os.RemoveAll(server_filename)
	type_store := NewTypeStore()
	type_store.AddType(reflect.TypeOf(Thingy{}), BuildThingy)
	server := NewServer(listener, TagSocketAll, &type_store)
	
	server.Accept("all", reflect.TypeOf(Thingy{}), func(iface interface{}) {
		if thingy, ok :=  iface.(*Thingy); ok {
			f, err := os.Create(fmt.Sprintf("%s-%d.0", thingy.Name, thingy.ID))
			t.Log(f)
			if err != nil {
				t.Errorf("could not write test file")
			}
			f.Close()
		}
	})
	server.Accept("all", reflect.TypeOf(Thingy{}), func(iface interface{}) {
		if thingy, ok :=  iface.(*Thingy); ok {
			f, err := os.Create(fmt.Sprintf("%s-%d.1", thingy.Name, thingy.ID))
			if err != nil {
				t.Errorf("could not write test file")
			}
			f.Close()
		}
	})
	if server.Events["all"] == nil {
		t.Errorf("calls to Accept didn't create records in server Events")
	}
	code, _ := type_store.LookupCode(reflect.TypeOf(Thingy{}))
	if server.Events["all"][code] == nil {
		t.Errorf("calls to Accept didn't create records in server Events")
	}
	if len(server.Events["all"][code]) != 2 {
		t.Errorf("two calls to Accept didn't create two records in server Events")
	}

	client_socket, err := net.Dial("unix", server_filename)
	if err != nil {
		t.Errorf("could not connect to unix server")
	}
	defer client_socket.Close()
	
	thingy := Thingy {
		Name:	"accept-test",
		ID:		1,
	}
	formatted_thingy, err := format(thingy, &type_store)
	if err != nil {
			t.Errorf("could not format thingy")
	}
	client_socket.Write(formatted_thingy)
	time.Sleep(50 * time.Millisecond)		// wait for server to process incoming struct
	// try to read struct directly
	//obj, err := nextStruct(*server.Sockets["all"][0], server.TypeStore)
	//if err != nil { }//panic(err) }
	//t.Log(obj)
	//t.Errorf(fmt.Sprintf("%s", <- server.FailedSockets))
	
	//thingy_bytes, _ := json.Marshal(thingy)
	//server.Events["all"][code][0](type_store.BuildType(code, thingy_bytes))
	//server.Events["all"][code][1](type_store.BuildType(code, thingy_bytes))
	
	if _, err := os.Stat("accept-test-1.0"); os.IsNotExist(err) {
		t.Errorf("test file not created when struct recieved")
	}
	if _, err := os.Stat("accept-test-1.1"); os.IsNotExist(err) {
		t.Errorf("test file not created when struct recieved")
	}
	os.RemoveAll("accept-test-1.0")
	os.RemoveAll("accept-test-1.1")
}

func TestResponderCanSendResponse(t *testing.T) {
}

func TestServerCanRunAcceptRequestCallbacks(t *testing.T) {
}
*/
