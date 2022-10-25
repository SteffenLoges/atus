package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) writePump(h *Hub) {
	ticker := time.NewTicker(h.PingPeriod)
	defer ticker.Stop()
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.sendChan:
			c.conn.SetWriteDeadline(time.Now().Add(h.WriteWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// append queued messages
			n := len(c.sendChan)
			for i := 0; i < n; i++ {
				w.Write(messageSeperator)
				w.Write(<-c.sendChan)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(h.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
