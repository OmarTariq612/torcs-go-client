package client

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	c "github.com/OmarTariq612/torcs-go-client/controller"
)

type Client struct {
	serverAddr *net.UDPAddr
	c.CarController
	action c.CarAction
	state  c.CarState
}

func New(serverAddr *net.UDPAddr, controller c.CarController) *Client {
	return &Client{serverAddr: serverAddr, CarController: controller}
}

func (c *Client) Start() error {
	fmt.Println(c.serverAddr)
	udpConn, err := net.DialUDP(c.serverAddr.Network(), nil, c.serverAddr)
	if err != nil {
		return err
	}
	defer udpConn.Close()

	var initMessage strings.Builder
	var buffer [100]byte
	initMessage.WriteString("SCR") // client id
	initMessage.WriteString("(init")
	angles := c.InitAngles()
	for _, angle := range angles {
		initMessage.WriteRune(' ')
		initMessage.WriteString(string(strconv.AppendFloat(buffer[:0], angle, 'f', -1, 64)))
	}
	initMessage.WriteRune(')')
	initMessageBytes := []byte(initMessage.String())

	var udpMessage [1000]byte

	for {
		if _, err := udpConn.Write(initMessageBytes); err != nil {
			return err
		}
		fmt.Printf("<< %s\n\n", string(initMessageBytes))
		if err := udpConn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			return err
		}
		n, err := udpConn.Read(udpMessage[:])
		if err != nil {
			if nErr, ok := err.(net.Error); !ok || !nErr.Timeout() {
				return err
			} else {
				fmt.Printf("<> didn't get response from server within %v sec\n", 1)
				continue
			}
		}
		fmt.Printf(">> %s\n\n", string(udpMessage[:n]))
		if strings.HasPrefix(string(udpMessage[:n]), "***identified***") {
			break
		}
	}

	if err := udpConn.SetReadDeadline(time.Time{}); err != nil {
		return err
	}

	for {
		n, err := udpConn.Read(udpMessage[:])
		if err != nil {
			return err
		}
		fmt.Printf(">> %s\n\n", string(udpMessage[:n]))

		str := string(udpMessage[:n])

		switch {
		case strings.HasPrefix(str, "***shutdown***"):
			fmt.Println("shutting down")
			return nil

		case strings.HasPrefix(str, "***restart***"):
			fmt.Println("restarting")
			return nil

		default:
			if err := c.state.UnmarshalText(udpMessage[:n]); err != nil {
				return err
			}
			c.Control(&c.state, &c.action)
			actionBytes, err := c.action.MarshalText()
			fmt.Printf("<< %s\n\n", string(actionBytes))
			if err != nil {
				return err
			}
			if _, err := udpConn.Write(actionBytes); err != nil {
				return err
			}
		}
	}
}
