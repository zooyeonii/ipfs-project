// pubsub
package pubsub

import (
	"encoding/json"
	"io"
	
	peer "github.com/libp2p/go-libp2p-core/peer"
)
//Message is a pubsub message.
//구조체 
type Message struct{
	From peer.ID
	Data []byte
	Seqno []byte
	TopicIDs []string
}

//PubSubSubscription allow you to receive pubsub records that where published
//on the network
// 클로저 (closer) 
type PubSubSubscription struct{
	resp io.Closer
	dec *json.Decoder
}

func NewPubSubSubscription(resp io.ReadCloser) *PubSubSubscription {
	return &PubSubSubscription{
		resp: resp,
		dec : json.NewDecoder(resp),
	}
}

//Next waits for the next record ans returns that.
func (s *PubSubSubscription) Next() (*Message, error){
	var r struct{
		From []byte 'json:"from,omitempty"' //omitempty : 비워두다 
		Data []byte 'json:"data,omitempty"'
		Seqno []byte 'json:"seqno,omitempty"'
		TopicIDs []string 'json:"topicIDs,omitempty"'
	}
	err := s.dec.Decode(&r)
	if err!=nil {
		return nil, err
	}
	
	from, err:= peer.IDFromBytes(r.From)
	if err!=nil{
		return nil, err
	}
	return &Message{
		From: from, 
		Data: r.Data,
		Seqno: r.Seqno,
		TopicIDs : r.TopicIDs,
	}, nil
}
//Cancels the given subscription
func (s *PubSubSubscription) Cancel() error{
	return s.resp.Close()
}    