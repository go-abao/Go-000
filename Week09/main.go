package main

// 用 Go 实现一个 tcp server ，用两个 goroutine 读写 conn，两个 goroutine 通过 chan 可以传递 message，能够正确退出
import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	server := NewServer("172.16.11.227", "8088")

	g.Go(func() error {
		go func() {
			<-ctx.Done()
			server.Shutdown()
		}()
		return server.Start(ctx)
	})

	g.Go(func() error {

		exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
		c := make(chan os.Signal, len(exitSignals))
		signal.Notify(c, exitSignals...)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return server.Shutdown()
			}
		}
	})

	g.Go(func() error {
		time.Sleep(time.Second)
		conn, err := net.Dial("tcp", "172.16.11.227:8088")
		if err != nil {
			fmt.Printf("Tcp dial error: %v\n", err)
			return err
		}
		fmt.Println("Connect server success")
		defer conn.Close()

		go func() {
			for {
				data := make([]byte, 1024)
				_, err := conn.Read(data)
				if err != nil {
					return
				}
				fmt.Println(string(data))
			}
		}()

		for i := 0; i < 10; i++ {
			_, err = conn.Write([]byte(fmt.Sprintf("[Client Write]: %s \n", "test message "+strconv.Itoa(time.Now().Nanosecond()))))
			if err != nil {
				fmt.Printf("Client write error %v\n", err)
				return err
			}
			time.Sleep(time.Second)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Println("done:", err)
	}
}

// Server Server
type Server struct {
	IP     string
	Port   string
	Listen net.Listener
}

// NewServer NewServer
func NewServer(ip, port string) *Server {
	return &Server{
		IP:   ip,
		Port: port,
	}
}

// Start Start
func (s *Server) Start(ctx context.Context) (err error) {
	fmt.Println("Server start")
	s.Listen, err = net.Listen("tcp", fmt.Sprintf("%s:%s", s.IP, s.Port))
	if err != nil {
		return err
	}

	for {
		conn, err := s.Listen.Accept()
		if err != nil {
			return err
		}

		m := make(chan string, 1)
		go handlerRConn(ctx, conn, m)
		go handlerWConn(ctx, conn, m)

		fmt.Printf("Connection close")
	}
}

// handlerRConn handlerRConn
func handlerRConn(ctx context.Context, conn net.Conn, m chan<- string) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Server read conn close caused by context done\n")
			close(m)
			return
		default:
			line, _, err := r.ReadLine()
			fmt.Println(string(line))
			if err != nil {
				fmt.Printf("Server read conn err %v\n", err)
				return
			}
			m <- string(line)
		}
	}
}

// handlerWConn handlerWConn
func handlerWConn(ctx context.Context, conn net.Conn, m chan string) {
	defer conn.Close()
	wr := bufio.NewWriter(conn)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Server write conn close caused by context done\n")
			close(m)
			return
		default:
			for c := range m {
				_, err := wr.WriteString(fmt.Sprintf("[Server reply]: received [%s]", c))
				if err != nil {
					fmt.Printf("Server write error: %v\n", err)
					return
				}
				wr.Flush()
			}
		}
	}
}

// Shutdown Shutdown
func (s *Server) Shutdown() error {
	fmt.Println("Server shutdown")
	return s.Listen.Close()
}
