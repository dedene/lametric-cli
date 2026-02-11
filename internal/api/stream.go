package api

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

// LMSP protocol constants.
const (
	LMSPPort    = 7575
	LMSPVersion = 0x0100

	EncodingRAW  = 0
	EncodingPNG  = 1
	EncodingJPEG = 2
	EncodingGIF  = 3

	lmspMagic     = "lmsp"
	lmspHeaderLen = 36
)

// StreamSession represents an active LMSP streaming session.
type StreamSession struct {
	SessionID string // UUID from start response
	DeviceIP  string
	Width     int // 37 for TIME, 8 for SKY
	Height    int // always 8
	conn      *net.UDPConn
	uuidBytes [16]byte
}

// StreamStartResponse is returned by the streaming start endpoint.
type StreamStartResponse struct {
	SessionID string `json:"session_id"`
}

// StreamStatus is returned by the streaming status endpoint.
type StreamStatus struct {
	Active    bool   `json:"active"`
	SessionID string `json:"session_id,omitempty"`
}

// StartStream initiates a streaming session via POST /api/v2/device/streaming/start.
func (c *Client) StartStream(ctx context.Context) (*StreamStartResponse, error) {
	var resp StreamStartResponse
	if err := c.Post(ctx, "/api/v2/device/streaming/start", nil, &resp); err != nil {
		return nil, fmt.Errorf("start stream: %w", err)
	}
	return &resp, nil
}

// StopStream ends the streaming session via DELETE /api/v2/device/streaming/stop.
func (c *Client) StopStream(ctx context.Context) error {
	if err := c.Delete(ctx, "/api/v2/device/streaming/stop"); err != nil {
		return fmt.Errorf("stop stream: %w", err)
	}
	return nil
}

// GetStreamStatus returns current streaming status via GET /api/v2/device/streaming.
func (c *Client) GetStreamStatus(ctx context.Context) (*StreamStatus, error) {
	var status StreamStatus
	if err := c.Get(ctx, "/api/v2/device/streaming", &status); err != nil {
		return nil, fmt.Errorf("get stream status: %w", err)
	}
	return &status, nil
}

// NewStreamSession creates a StreamSession and opens the UDP connection.
func NewStreamSession(deviceIP, sessionID string, width, height int) (*StreamSession, error) {
	s := &StreamSession{
		SessionID: sessionID,
		DeviceIP:  deviceIP,
		Width:     width,
		Height:    height,
	}

	uuid, err := parseUUID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("parse session UUID: %w", err)
	}
	s.uuidBytes = uuid

	addr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort(deviceIP, fmt.Sprintf("%d", LMSPPort)))
	if err != nil {
		return nil, fmt.Errorf("resolve UDP addr: %w", err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("dial UDP: %w", err)
	}
	s.conn = conn

	return s, nil
}

// Close closes the UDP connection.
func (s *StreamSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// SendFrame sends a single frame via UDP with LMSP header.
func (s *StreamSession) SendFrame(encoding byte, data []byte) error {
	header := s.buildHeader(encoding, len(data))
	packet := make([]byte, 0, len(header)+len(data))
	packet = append(packet, header...)
	packet = append(packet, data...)

	_, err := s.conn.Write(packet)
	if err != nil {
		return fmt.Errorf("send frame: %w", err)
	}
	return nil
}

// buildHeader creates the 36-byte LMSP header (little-endian).
func (s *StreamSession) buildHeader(encoding byte, dataLen int) []byte {
	h := make([]byte, lmspHeaderLen)

	// 0-3: magic "lmsp"
	copy(h[0:4], lmspMagic)

	// 4-5: version
	binary.LittleEndian.PutUint16(h[4:6], LMSPVersion)

	// 6-21: session UUID (16 bytes)
	copy(h[6:22], s.uuidBytes[:])

	// 22: encoding
	h[22] = encoding

	// 23: reserved
	h[23] = 0

	// 24: canvas count
	h[24] = 1

	// 25: reserved
	h[25] = 0

	// 26-27: X offset
	binary.LittleEndian.PutUint16(h[26:28], 0)

	// 28-29: Y offset
	binary.LittleEndian.PutUint16(h[28:30], 0)

	// 30-31: width
	binary.LittleEndian.PutUint16(h[30:32], uint16(s.Width))

	// 32-33: height
	binary.LittleEndian.PutUint16(h[32:34], uint16(s.Height))

	// 34-35: data length
	binary.LittleEndian.PutUint16(h[34:36], uint16(dataLen))

	return h
}

// parseUUID converts a UUID string (with or without dashes) to 16 bytes.
func parseUUID(s string) ([16]byte, error) {
	var uuid [16]byte
	clean := strings.ReplaceAll(s, "-", "")
	if len(clean) != 32 {
		return uuid, fmt.Errorf("invalid UUID length: %q", s)
	}
	b, err := hex.DecodeString(clean)
	if err != nil {
		return uuid, fmt.Errorf("invalid UUID hex: %w", err)
	}
	copy(uuid[:], b)
	return uuid, nil
}
