package pearl

import (
	"io"
	"net"
	"sync"

	"github.com/mmcloughlin/pearl/tls"

	"github.com/mmcloughlin/pearl/log"
	"github.com/pkg/errors"
)

// CellSender can send a Cell.
type CellSender interface {
	SendCell(Cell) error
}

// CellReceiver can receive Cells.
type CellReceiver interface {
	ReceiveCell() (Cell, error)
}

// CellReceiver can receive legacy Cells (circ ID length 2).
type LegacyCellReceiver interface {
	CellReceiver
	ReceiveLegacyCell() (Cell, error)
}

// Link is a Cell communication layer.
type Link interface {
	CellSender
	CellReceiver
}

type link struct {
	CellSender
	CellReceiver
}

func NewLink(s CellSender, r CellReceiver) Link {
	return link{
		CellSender:   s,
		CellReceiver: r,
	}
}

type CellChan chan Cell

func (ch CellChan) SendCell(cell Cell) error {
	ch <- cell
	return nil
}

func (ch CellChan) ReceiveCell() (Cell, error) {
	cell, ok := <-ch
	if !ok {
		return nil, errors.New("closed channel") // XXX correct return?
	}
	return cell, nil
}

type CircuitLink interface {
	CircID() CircID
	Link
}

type circLink struct {
	id CircID
	Link
}

func NewCircuitLink(id CircID, lk Link) CircuitLink {
	return circLink{
		id:   id,
		Link: lk,
	}
}

func (c circLink) CircID() CircID { return c.id }

// Connection encapsulates a router connection.
type Connection struct {
	router      *Router
	tlsCtx      *TLSContext
	tlsConn     *tls.Conn
	fingerprint []byte
	outbound    bool

	channels *ChannelManager

	rw io.ReadWriter
	CellReceiver
	CellSender

	logger log.Logger
}

// NewServer constructs a server connection.
func NewServer(r *Router, conn net.Conn, logger log.Logger) (*Connection, error) {
	tlsCtx, err := NewTLSContext(r.IdentityKey())
	if err != nil {
		return nil, err
	}
	tlsConn := tlsCtx.ServerConn(conn)
	c := newConnection(r, tlsCtx, tlsConn, logger.With("role", "server"))
	c.outbound = false
	return c, nil
}

// NewClient constructs a client-side connection.
func NewClient(r *Router, conn net.Conn, logger log.Logger) (*Connection, error) {
	tlsCtx, err := NewTLSContext(r.IdentityKey())
	if err != nil {
		return nil, err
	}
	tlsConn := tlsCtx.ClientConn(conn)
	c := newConnection(r, tlsCtx, tlsConn, logger.With("role", "client"))
	c.outbound = true
	return c, nil
}

func newConnection(r *Router, tlsCtx *TLSContext, tlsConn *tls.Conn, logger log.Logger) *Connection {
	rw := tlsConn // TODO(mbm): use bufio
	return &Connection{
		router:      r,
		tlsCtx:      tlsCtx,
		tlsConn:     tlsConn,
		fingerprint: nil,

		channels: NewChannelManager(),

		rw:           rw,
		CellReceiver: NewCellReader(rw, logger),
		CellSender:   NewCellWriter(rw, logger),

		logger: log.ForConn(logger, tlsConn),
	}
}

func (c *Connection) newHandshake() *Handshake {
	return &Handshake{
		Conn:        c.tlsConn,
		Link:        NewHandshakeLink(c.rw, c.logger),
		TLSContext:  c.tlsCtx,
		IdentityKey: &c.router.idKey.PublicKey,
		logger:      c.logger,
	}
}

// Fingerprint returns the fingerprint of the connected peer.
func (c *Connection) Fingerprint() (Fingerprint, error) {
	if c.fingerprint == nil {
		return Fingerprint{}, errors.New("peer fingerprint not established")
	}
	return NewFingerprintFromBytes(c.fingerprint)
}

func (c *Connection) Serve() error {
	c.logger.Info("handle")

	h := c.newHandshake()
	err := h.Server()
	if err != nil {
		log.Err(c.logger, err, "server handshake failed")
		return nil
	}

	// TODO(mbm): register connection

	return c.readLoop()
}

func (c *Connection) StartClient() error {
	h := c.newHandshake()
	err := h.Client()
	if err != nil {
		log.Err(c.logger, err, "client handshake failed")
	}

	// TODO(mbm): register connection

	// TODO(mbm): goroutine management
	go c.readLoop()

	return nil
}

func (c *Connection) readLoop() error {
	var err error
	var cell Cell
	for {
		cell, err = c.ReceiveCell()
		if err != nil {
			return errors.Wrap(err, "could not read cell")
		}

		logger := CellLogger(c.logger, cell)
		logger.Trace("received cell")

		switch cell.Command() {
		// Cells to be handled by this Connection
		case Create2:
			Create2Handler(c, cell) // XXX error return
		// Cells related to a circuit
		case Created2, Relay, RelayEarly, Destroy:
			logger.Trace("directing cell to circuit channel")
			ch, ok := c.channels.Channel(cell.CircID())
			if !ok {
				// BUG(mbm): is logging the correct behavior
				logger.Error("unrecognized circ id")
				continue
			}
			ch <- cell
		// Cells to be ignored
		case Padding, Vpadding:
			logger.Debug("skipping padding cell")
		// Something which shouldn't happen
		default:
			logger.Error("no handler registered")
		}
	}
}

// GenerateCircuitLink
func (c *Connection) GenerateCircuitLink() CircuitLink {
	// BUG(mbm): what if c.proto has not been established
	id, ch := c.channels.New(c.outbound)
	return NewCircuitLink(id, NewLink(c, CellChan(ch)))
}

// NewCircuitLink
func (c *Connection) NewCircuitLink(id CircID) (CircuitLink, error) {
	ch, err := c.channels.NewWithID(id)
	if err != nil {
		return nil, err
	}
	return NewCircuitLink(id, NewLink(c, CellChan(ch))), nil
}

func CellLogger(l log.Logger, cell Cell) log.Logger {
	return l.With("cmd", cell.Command()).With("circid", cell.CircID())
}

// ChannelManager manages a collection of cell channels.
type ChannelManager struct {
	channels map[CircID]chan Cell

	sync.RWMutex
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channels: make(map[CircID]chan Cell),
	}
}

func (m *ChannelManager) New(outbound bool) (CircID, chan Cell) {
	m.Lock()
	defer m.Unlock()

	// Reference: https://github.com/torproject/torspec/blob/4074b891e53e8df951fc596ac6758d74da290c60/tor-spec.txt#L931-L933
	//
	//	   In link protocol version 4 or higher, whichever node initiated the
	//	   connection sets its MSB to 1, and whichever node didn't initiate the
	//	   connection sets its MSB to 0.
	//
	msb := uint32(0)
	if outbound {
		msb = uint32(1)
	}

	// BUG(mbm): potential infinite (or at least long) loop to find a new id
	for {
		id := GenerateCircID(msb)
		// 0 is reserved
		if id == 0 {
			continue
		}
		_, exists := m.channels[id]
		if exists {
			continue
		}
		ch := m.newWithID(id)
		return id, ch
	}
}

func (m *ChannelManager) NewWithID(id CircID) (chan Cell, error) {
	m.Lock()
	defer m.Unlock()
	_, exists := m.channels[id]
	if exists {
		return nil, errors.New("cannot override existing channel id")
	}
	return m.newWithID(id), nil
}

func (m *ChannelManager) newWithID(id CircID) chan Cell {
	ch := make(chan Cell)
	m.channels[id] = ch
	return ch
}

func (m *ChannelManager) Channel(id CircID) (chan Cell, bool) {
	m.RLock()
	defer m.RUnlock()
	ch, ok := m.channels[id]
	return ch, ok
}
