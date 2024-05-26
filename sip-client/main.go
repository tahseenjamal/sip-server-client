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
	clientIP   = "127.0.0.1"
	clientPort = 5061
)

func generateCallID() string {
	return fmt.Sprintf("%d@%s", rand.Int(), clientIP)
}

func generateFromTag() string {
	return fmt.Sprintf("%dSIPpTag%03d", rand.Int(), rand.Int())
}

func createSIPMessage(method, callID, fromTag, toTag, branch string, cSeq int) string {

	return fmt.Sprintf("%s sip:service@%s:%d SIP/2.0\r\n"+
		"Call-ID: %s\r\n"+
		"Contact: sip:sipp@%s:%d\r\n"+
		"Content-Length: %d\r\n"+
		"CSeq: %d %s\r\n"+
		"From: sipp <sip:sipp@%s:%d>;tag=%s\r\n"+
		"Max-Forwards: 70\r\n"+
		"Subject: Performance Test\r\n"+
		"To: service <sip:service@%s:%d>%s\r\n"+
		"Via: SIP/2.0/UDP %s:%d;branch=%s\r\n"+
		"\r\n"+
		"%s",
		method, serverIP, serverPort,
		callID,
		clientIP, clientPort,
		getContentLength(method),
		cSeq, method, clientIP, clientPort, fromTag,
		serverIP, serverPort, formatToTag(toTag), clientIP, clientPort, branch,
		getSDPBody(method))
}

func formatToTag(toTag string) string {
	if toTag != "" && !strings.HasPrefix(toTag, ";tag=") {
		return ";tag=" + toTag
	}
	return toTag
}

func getContentLength(method string) int {
	contentLength := len(getSDPBody(method))
	if contentLength <= 0 {
		return 0
	}
	return contentLength
}

func getSDPBody(method string) string {
	if method == "INVITE" || method == "BYE" {

		return "v=0\r\n" +
			"o=user1 53655765 2353687637 IN IP6 [::1]\r\n" +
			"s=-\r\n" +
			"c=IN IP6 ::1\r\n" +
			"t=0 0\r\n" +
			"m=audio 6000 RTP/AVP 0\r\n" +
			"a=rtpmap:0 PCMU/8000\r\n" +
			"a=sendrecv"

	}
	return ""
}

func generateBranch() string {
	return fmt.Sprintf("z9hG4bK-%d", rand.Int())
}

func sendSIPMessage(conn *net.UDPConn, message string) ([]byte, error) {
	message = strings.TrimSpace(message)
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return nil, err
	}
	fmt.Printf("Sent message:\n%s\n", message)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buffer := make([]byte, 1024)
	return buffer, nil
}

func readSIPResponse(conn *net.UDPConn, buffer []byte) (*sip.Msg, error) {
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, fmt.Errorf("error receiving response: %v", err)
	}

	msg, err := sip.ParseMsg(buffer[:n])
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	fmt.Printf("Received response:\n%s\n", msg)
	return msg, nil
}

func main() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(serverIP),
		Port: serverPort,
	})
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	callID := generateCallID()
	fromTag := generateFromTag()

	// Sending SIP INVITE message
	inviteMessage := createSIPMessage("INVITE", callID, fromTag, "", generateBranch(), 1)
	buffer, err := sendSIPMessage(conn, inviteMessage)
	if err != nil {
		return
	}
	msg, err := readSIPResponse(conn, buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extracting toTag from the received response
	toTag := msg.To.Param.Get("tag").Value

	// Sending SIP ACK message
	ackMessage := createSIPMessage("ACK", callID, fromTag, toTag, generateBranch(), 1)
	buffer, err = sendSIPMessage(conn, ackMessage)
	if err != nil {
		return
	}
	if _, err = readSIPResponse(conn, buffer); err != nil {
		fmt.Println(err)
		return
	}

	// Sending SIP BYE message
	byeMessage := createSIPMessage("BYE", callID, fromTag, toTag, generateBranch(), 2)
	buffer, err = sendSIPMessage(conn, byeMessage)
	if err != nil {
		return
	}
	if _, err = readSIPResponse(conn, buffer); err != nil {
		fmt.Println(err)
		return
	}
}
