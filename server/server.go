package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/jart/gosip/sip"
)

const (
	serverIP   = "127.0.0.1"
	serverPort = 5060
)

func createSIPResponse(
	statusCode int,
	reasonPhrase, method, callID, fromTag, toTag, viaBranch string,
	cSeq int,
) string {
	return fmt.Sprintf("SIP/2.0 %d %s\r\n"+
		"Call-ID: %s\r\n"+
		"CSeq: %d %s\r\n"+
		"From: sipp <sip:sipp@%s:%d>;tag=%s\r\n"+
		"To: service <sip:service@%s:%d>;tag=%s\r\n"+
		"Via: SIP/2.0/UDP %s:%d;branch=%s\r\n"+
		"Contact: <sip:%s:%d;transport=UDP>\r\n"+
		"Content-Length: 0\r\n"+
		"\r\n",
		statusCode, reasonPhrase,
		callID, cSeq, method,
		serverIP, serverPort, fromTag,
		serverIP, serverPort, toTag,
		serverIP, serverPort, viaBranch,
		serverIP, serverPort)
}

func generateTag() string {
	return fmt.Sprintf("%dSIPpTag%03d", rand.Int(), rand.Int())
}

func handleInvite(conn *net.UDPConn, addr *net.UDPAddr, msg *sip.Msg) {
	fromTag := msg.From.Param.Get("tag").Value
	toTag := generateTag()
	viaBranch := msg.Via.Param.Get("branch").Value
	callID := msg.CallID
	cSeq := msg.CSeq

	// Send 180 Ringing
	ringingResponse := createSIPResponse(
		180,
		"Ringing",
		"INVITE",
		callID,
		fromTag,
		toTag,
		viaBranch,
		cSeq,
	)
	_, err := conn.WriteToUDP([]byte(ringingResponse), addr)
	if err != nil {
		fmt.Printf("Error sending 180 Ringing response: %v\n", err)
		return
	}
	fmt.Printf("Sent 180 Ringing response:\n%s\n", ringingResponse)

	time.Sleep(1 * time.Second)

	// Send 200 OK
	okResponse := createSIPResponse(200, "OK", "INVITE", callID, fromTag, toTag, viaBranch, cSeq)
	_, err = conn.WriteToUDP([]byte(okResponse), addr)
	if err != nil {
		fmt.Printf("Error sending 200 OK response: %v\n", err)
		return
	}
	fmt.Printf("Sent 200 OK response:\n%s\n", okResponse)
}

func handleAck(conn *net.UDPConn, addr *net.UDPAddr, msg *sip.Msg) {
	fmt.Println("Received ACK")
}

func handleBye(conn *net.UDPConn, addr *net.UDPAddr, msg *sip.Msg) {
	fromTag := msg.From.Param.Get("tag").Value
	toTag := msg.To.Param.Get("tag").Value
	viaBranch := msg.Via.Param.Get("branch").Value
	callID := msg.CallID
	cSeq := msg.CSeq

	// Send 200 OK for BYE
	okResponse := createSIPResponse(200, "OK", "BYE", callID, fromTag, toTag, viaBranch, cSeq)
	_, err := conn.WriteToUDP([]byte(okResponse), addr)
	if err != nil {
		fmt.Printf("Error sending 200 OK response: %v\n", err)
		return
	}
	fmt.Printf("Sent 200 OK response:\n%s\n", okResponse)
}

func main() {
	addr := net.UDPAddr{
		IP:   net.ParseIP(serverIP),
		Port: serverPort,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Server listening on %s:%d\n", serverIP, serverPort)

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}

		msg, err := sip.ParseMsg(buffer[:n])
		if err != nil {
			fmt.Printf("Error parsing SIP message: %v\n", err)
			continue
		}

		switch {
		case strings.HasPrefix(msg.Method, "INVITE"):
			fmt.Printf("Received INVITE message:\n%s\n", msg)
			handleInvite(conn, clientAddr, msg)
		case strings.HasPrefix(msg.Method, "ACK"):
			fmt.Printf("Received ACK message:\n%s\n", msg)
			handleAck(conn, clientAddr, msg)
		case strings.HasPrefix(msg.Method, "BYE"):
			fmt.Printf("Received BYE message:\n%s\n", msg)
			handleBye(conn, clientAddr, msg)
		}
	}
}
