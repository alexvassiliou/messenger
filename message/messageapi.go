package messageapi

import (
	context "context"
	fmt "fmt"
	"sync"
)

type Connection struct {
	stream MessageService_CreateStreamServer
	id     string
	active bool
	error  chan error
}

// Server object containing the streaming connections
type Server struct {
	Connection []*Connection
}

// CreateStream opens a stream of messages from the client
func (s *Server) CreateStream(connect *Connect, stream MessageService_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     connect.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn)
	return <-conn.error
}

// BroadcastStream broadcasts the message received from CreateStream back to the client
func (s *Server) BroadcastStream(ctx context.Context, m *Message) (*Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		fmt.Println(conn.id)
		fmt.Println(conn.stream)
		wait.Add(1)

		go func(m *Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(m)
				fmt.Println("sending message to: ", conn.stream)

				if err != nil {
					fmt.Printf("Error with Stream: %v - Error: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(m, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &Close{}, nil
}
