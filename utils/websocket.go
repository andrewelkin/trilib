package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	ul "github.com/andrewelkin/trilib/utils"
	"github.com/andrewelkin/trilib/utils/logger"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const nSpace = "WS"

type Empty struct {
}

type Socket struct {
	ctx               context.Context
	Conn              *websocket.Conn
	WebsocketDialer   *websocket.Dialer
	Url               string
	ConnectionOptions ConnectionOptions
	RequestHeader     http.Header
	OnConnected       func(socket *Socket)
	OnTextMessage     func(message string, socket *Socket)
	OnBinaryMessage   func(data []byte, socket *Socket)
	OnReadError       func(err error, socket *Socket)
	OnConnectError    func(err error, socket *Socket)
	OnDisconnected    func(err error, socket *Socket)
	OnPingReceived    func(data string, socket *Socket)
	OnPongReceived    func(data string, socket *Socket)
	IsConnected       bool
	sendMu            *sync.Mutex // Prevent "concurrent write to websocket connection"

	log logger.Logger
}

type ConnectionOptions struct {
	UseCompression bool
	UseSSL         bool
	Proxy          func(*http.Request) (*url.URL, error)
	Subprotocols   []string
}

type ReconnectionOptions struct {
}

func NewWebSocket(ctx context.Context, url string, logger logger.Logger, header http.Header, readBufSize, writeBufSize int) *Socket {

	return &Socket{
		ctx:           ctx,
		Url:           url,
		log:           logger,
		RequestHeader: header,
		ConnectionOptions: ConnectionOptions{
			UseCompression: false,
			UseSSL:         true,
		},
		WebsocketDialer: &websocket.Dialer{
			ReadBufferSize:  readBufSize,
			WriteBufferSize: writeBufSize,
		},

		sendMu: &sync.Mutex{},
	}
}

func (socket *Socket) setConnectionOptions() {
	socket.WebsocketDialer.EnableCompression = socket.ConnectionOptions.UseCompression
	socket.WebsocketDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: socket.ConnectionOptions.UseSSL}
	socket.WebsocketDialer.Proxy = socket.ConnectionOptions.Proxy
	socket.WebsocketDialer.Subprotocols = socket.ConnectionOptions.Subprotocols
}

func (socket *Socket) Connect() error {
	var err error
	socket.setConnectionOptions()
	socket.Conn, _, err = socket.WebsocketDialer.DialContext(socket.ctx, socket.Url, socket.RequestHeader)
	if err != nil {
		socket.log.Warnf(nSpace, "Error while connecting to server '%s' : %v", socket.Url, err)
		socket.IsConnected = false
		if socket.OnConnectError != nil {
			socket.OnConnectError(err, socket)
		}
		return err
	}

	if socket.OnConnected != nil {
		socket.IsConnected = true
		socket.OnConnected(socket)
	} else {
		socket.log.Debugf(nSpace, "Connected to server %v", socket.Url)
	}

	socket.Conn.SetCloseHandler(nil)
	socket.Conn.SetPingHandler(nil)
	socket.Conn.SetPongHandler(nil)
	defaultPingHandler := socket.Conn.PingHandler()
	socket.Conn.SetPingHandler(func(appData string) error {
		//socket.log.Debugf(nSpace, "Received PING from server")
		if socket.OnPingReceived != nil {
			socket.OnPingReceived(appData, socket)
		}
		return defaultPingHandler(appData)
	})
	defaultPongHandler := socket.Conn.PongHandler()
	socket.Conn.SetPongHandler(func(appData string) error {
		//socket.log.Debugf(nSpace, "Received PONG from server")
		if socket.OnPongReceived != nil {
			socket.OnPongReceived(appData, socket)
		}
		return defaultPongHandler(appData)
	})
	defaultCloseHandler := socket.Conn.CloseHandler()
	socket.Conn.SetCloseHandler(func(code int, text string) error {
		result := defaultCloseHandler(code, text)
		if socket.OnDisconnected != nil {
			socket.IsConnected = false
			socket.OnDisconnected(errors.New(text), socket)
		} else {
			socket.log.Debugf(nSpace, "Disconnected from the server %v", socket.Url)
		}

		return result
	})

	go func() {

		ul.TryBlock{
			Try: func() {
				for {
					messageType, message, err := socket.Conn.ReadMessage()
					if err != nil {
						if socket.OnReadError != nil {
							socket.OnReadError(err, socket)
						} else {
							socket.log.Warnf(nSpace, "read error: %v", err)
						}
						return
					}

					switch messageType {
					case websocket.TextMessage:
						if socket.OnTextMessage != nil {
							socket.OnTextMessage(string(message), socket)
						}
					case websocket.BinaryMessage:
						if socket.OnBinaryMessage != nil {
							socket.OnBinaryMessage(message, socket)
						}
					}
				}
			},
			Catch: func(e ul.Exception) {
				if socket.OnReadError != nil {
					go socket.OnReadError(err, socket)
				} else {
					socket.log.Debugf(nSpace, "Exception while reading from the server %v", socket.Url)
				}
				return
			},
			Finally: func() {
			},
		}.Do()

	}()
	return nil
}

func (socket *Socket) SendPingFrame() error {
	if !socket.IsConnected {
		return nil
	}
	// logger.Debugf(nSpace, "Forcing PING frame")
	err := socket.Conn.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(time.Second))
	if err == websocket.ErrCloseSent {
		return nil
	} else if e, ok := err.(net.Error); ok && e.Temporary() {
		return nil
	}
	return err
}

func (socket *Socket) SendPongFrame() error {
	if !socket.IsConnected {
		return nil
	}
	// logger.Debugf(nSpace, "Forcing PING frame")
	err := socket.Conn.WriteControl(websocket.PongMessage, []byte(""), time.Now().Add(time.Second))
	if err == websocket.ErrCloseSent {
		return nil
	} else if e, ok := err.(net.Error); ok && e.Temporary() {
		return nil
	}
	return err
}

func (socket *Socket) SendText(message string) error {
	return socket.send(websocket.TextMessage, []byte(message))
}

func (socket *Socket) SendBinary(data []byte) error {
	return socket.send(websocket.BinaryMessage, data)
}

func (socket *Socket) send(messageType int, data []byte) error {
	if !socket.IsConnected {
		return fmt.Errorf("can't send, disconnected socket")
	}
	socket.sendMu.Lock()
	err := socket.Conn.WriteMessage(messageType, data)
	socket.sendMu.Unlock()
	return err
}

func (socket *Socket) Close() {
	if !socket.IsConnected {
		return
	}
	err := socket.send(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		socket.log.Warnf(nSpace, "Tried to close websocket gracefully, got error:", err)
	}
	socket.Conn.Close()
	if socket.OnDisconnected != nil {
		socket.IsConnected = false
		socket.OnDisconnected(err, socket)
	}
}
