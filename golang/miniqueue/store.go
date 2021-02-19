package miniqueue

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"sync"
)

type storeError string

const (
	errTopicEmpty      = storeError("topic is empty")
	errTopicNotExist   = storeError("topic does not exist")
	errAckMsgNotExist  = storeError("msg to ack does not exist")
	errNackMsgNotExist = storeError("msg to nack does not exist")
	errBackMsgNotExist = storeError("msg to back does not exist")
)

func (s storeError) Error() string {
	return string(s)
}

type storer interface {
	// Insert inserts a new record for a given topic
	Insert(topic string, value value) error

	// GetNext will retrieve the next value in the topic
	// as well as the AckKey allowing future acking/nacking of the value
	GetNext(topic string) (val value, ackOffset int, err error)

	// Ack will acknowledge the processing of a value, removing it from
	// the topic entirely
	Ack(topic string, ackOffset int) error

	// Nack will negatively acknowledge the value on a given topic,
	// returning it to the front of the consumption queue
	Nack(topic string, ackOffset int) error

	// ... returning it to the back of the ...
	Back(topic string, ackOffset int) error

	// Close closes the store
	Close() error

	// Destroy removes the store from persistence.
	// This is a destructive operation
	Destroy()
}

const (
	topicFmt      = "%s-%d"
	headPosKeyFmt = "%s-head"
	tailPosKeyFmt = "%s-tail"

	ackTopicFmt      = "%s-ack-%d"
	ackTailPosKeyFmt = "%s-ack-head"
)

// store handles the underlying leveldb implementation
type store struct {
	path string
	db *leveldb.DB
	sync.Mutex
}

func newStore(dbPath string) storer {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open levelDB")
	}

	return &store{
		path: dbPath,
		db:   db,
	}
}

func (s *store) Insert(topic string, value value) error {
	s.Lock()
	defer s.Unlock()

	headPosKey := []byte(fmt.Sprintf(headPosKeyFmt, topic))
	tailPosKey := []byte(fmt.Sprintf(tailPosKeyFmt, topic))
	ackTailPosKey := []byte(fmt.Sprintf(ackTailPosKeyFmt, topic))

	exists, err := s.db.Has(tailPosKey, nil)
	if err != nil {
		return fmt.Errorf("checking for has: %v", err)
	}

	// The key already exists
	if exists {
		if _, err := appendValue(s.db, tailPosKeyFmt, topicFmt, topic, value); err != nil {
			return err
		}
		return nil
	}

	// write initial head position
	headPos := make([]byte, 8)
	binary.PutVarint(headPos, 0)
	if err := s.db.Put(headPosKey, headPos, nil); err != nil {
		return fmt.Errorf("putting head position value: %v", err)
	}

	// write initial ack topic head position
	ackTailPos := make([]byte, 8)
	binary.PutVarint(ackTailPos, 0)
	if err := s.db.Put(ackTailPosKey, ackTailPos, nil); err != nil {
		return fmt.Errorf("putting ack head position value: %v", err)
	}

	// write initial tail position
	tailPos := make([]byte, 8)
	binary.PutVarint(tailPos, 1)
	if err := s.db.Put(tailPosKey, tailPos, nil); err != nil {
		return fmt.Errorf("putting tail position value: %v", err)
	}

	newKey := []byte(fmt.Sprintf(topicFmt, topic, 0))
	if err := s.db.Put(newKey, value, nil); err != nil {
		return fmt.Errorf("putting first value for topic: %v", err)
	}

	return nil
}

// retrieves the first record for a topic, incrementing the head position
// of the main array and pushing the value onto the ack array
func (s *store) GetNext(topic string) (val value, ackOffset int, err error) {
	s.Lock()
	defer s.Unlock()

	headOffset, err := getPos(s.db, headPosKeyFmt, topic)
	if err != nil {
		return nil, 0, err
	}

	val, err = getValue(s.db, topicFmt, topic, headOffset)
	if err != nil {
		return nil, 0, err
	}

	insertOffset, err := appendValue(s.db, ackTailPosKeyFmt, ackTopicFmt, topic, val)
	if err != nil {
		return nil, 0, err
	}

	if _, _, err := addPos(s.db, headPosKeyFmt, topic, 1); err != nil {
		return nil, 0, err
	}

	return val, insertOffset, nil
}

// Ack will acknowledge the processing of a value, removing it from the topic
// entirely.
func (s *store) Ack(topic string, ackOffset int) error {
	s.Lock()
	defer s.Unlock()

	key := fmt.Sprintf(ackTopicFmt, topic, ackOffset)
	if err := s.db.Delete([]byte(key), nil); err != nil {
		return fmt.Errorf("delete from ack topic: %v", err)
	}

	return nil
}

func (s *store) Nack(topic string, ackOffset int) error {
	s.Lock()
	defer s.Unlock()

	nackKey := []byte(fmt.Sprintf(ackTopicFmt, topic, ackOffset))

	tx, err := s.db.OpenTransaction()
	if err != nil {
		return fmt.Errorf("opening transaction: %v", err)
	}

	exists, err := tx.Has(nackKey, nil)
	if err != nil {
		tx.Discard()
		return fmt.Errorf("checking for has: %v", err)
	}

	if !exists {
		tx.Discard()
		return errNackMsgNotExist
	}

	val, err := getOffsetTx(tx, ackTopicFmt, topic, ackOffset)
	if err != nil {
		tx.Discard()
		return fmt.Errorf("getting ack msg from topic %s at offset %d: %v", topic, ackOffset, err)
	}

	if _, err := prependValueTx(tx, headPosKeyFmt, topicFmt, topic, val); err != nil {
		tx.Discard()
		return fmt.Errorf("prepending value to topic %s: %v", topic, err)
	}

	if err := tx.Delete(nackKey, nil); err != nil {
		tx.Discard()
		return fmt.Errorf("deleting ackKey %s: %v", nackKey, err)
	}

	if err := tx.Commit(); err != nil {
		tx.Discard()
		return fmt.Errorf("committing nack transaction: %v", err)
	}

	return nil
}

func (s *store) Back(topic string, ackOffset int) error {
	s.Lock()
	defer s.Unlock()

	backKey := []byte(fmt.Sprintf(ackTopicFmt, topic, ackOffset))

	tx, err := s.db.OpenTransaction()
	if err != nil {
		return fmt.Errorf("opening transaction: %v", err)
	}

	exists, err := tx.Has(backKey, nil)
	if err != nil {
		tx.Discard()
		return fmt.Errorf("checking for has: %v", err)
	}
	if !exists {
		tx.Discard()
		return errBackMsgNotExist
	}

	val, err := getOffsetTx(tx, ackTopicFmt, topic, ackOffset)
	if err != nil {
		tx.Discard()
		return fmt.Errorf("getting ack msg from topic %s at offset %d: %v", topic, ackOffset, err)
	}

	if _, err := appendValueTx(tx, tailPosKeyFmt, topicFmt, topic, val); err != nil {
		tx.Discard()
		return fmt.Errorf("appending value to topic %s: %v", topic, err)
	}

	if err := tx.Delete(backKey, nil); err != nil {
		tx.Discard()
		return fmt.Errorf("deleting ackKey %s: %v", backKey, err)
	}

	if err := tx.Commit(); err != nil {
		tx.Discard()
		return fmt.Errorf("committing nack transaction: %v", err)
	}

	return nil
}

func (s *store) Close() error {
	return s.db.Close()
}

func (s *store) Destroy() {
	_ = s.Close()
	_ = os.RemoveAll(s.path)
}

func getOffsetTx(db *leveldb.Transaction, topicFmt string, topic string, offset int) (value, error) {
	key := fmt.Sprintf(topicFmt, topic, offset)

	val, err := db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	return val, nil
}

// getPos gets the integer position value for topic and key format
func getPos(db *leveldb.DB, keyFmt string, topic string) (int, error) {
	key := []byte(fmt.Sprintf(keyFmt, topic))
	pos, err := db.Get(key, nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		return 0, errTopicNotExist
	}

	if err != nil {
		return 0, fmt.Errorf("getting offset position position: %v", err)
	}

	i, err := binary.ReadVarint(bytes.NewReader(pos))
	if err != nil {
		return 0, fmt.Errorf("reading offset position varint: %v", err)
	}

	return int(i), nil
}

// getPosTx gets the integer position value (aka offset) for topic and key format.
func getPosTx(tx *leveldb.Transaction, keyFmt string, topic string) (int, error) {
	key := []byte(fmt.Sprintf(keyFmt, topic))

	pos, err := tx.Get(key, nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		return 0, errTopicNotExist
	}
	if err != nil {
		return 0, fmt.Errorf("getting offset position position: %v", err)
	}

	i, err := binary.ReadVarint(bytes.NewReader(pos))
	if err != nil {
		return 0, fmt.Errorf("reading offset position varint: %v", err)
	}

	return int(i), nil
}

func getValue(db *leveldb.DB, keyFmt string, topic string, offset int) (value, error) {
	key := fmt.Sprintf(keyFmt, topic, offset)

	val, err := db.Get([]byte(key), nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		return nil, errTopicEmpty
	}
	if err != nil {
		return nil, fmt.Errorf("getting value with fmt [%s] from topic %s at offset %d: %v", keyFmt, topic, offset, err)
	}

	return val, nil
}

func appendValue(db *leveldb.DB, tailPosKeyFmt, keyFmt, topic string, val value) (offset int, err error) {
	tailPosKey := []byte(fmt.Sprintf(tailPosKeyFmt, topic))

	// Fetch the current tail position
	tailPosVal, err := db.Get(tailPosKey, nil)
	if err != nil {
		return 0, fmt.Errorf("getting tail position from db: %v", err)
	}

	origOffset, err := binary.ReadVarint(bytes.NewReader(tailPosVal))
	if err != nil {
		return 0, fmt.Errorf("reading tail pos varint: %v", err)
	}

	newKey := []byte(fmt.Sprintf(keyFmt, topic, origOffset))
	if err := db.Put(newKey, val, nil); err != nil {
		return 0, fmt.Errorf("putting value: %v", err)
	}

	// update tail position
	tail := make([]byte, 8)
	binary.PutVarint(tail, origOffset+1)
	if err := db.Put(tailPosKey, tail, nil); err != nil {
		return 0, fmt.Errorf("putting new tail position: %v", err)
	}

	return int(origOffset), nil
}

// appendValueTx returns inserts a new value to the end of a topic given,
// returning the inserted offset.
func appendValueTx(tx *leveldb.Transaction, tailPosKeyFmt, keyFmt, topic string, val value) (offset int, err error) {
	tailPosKey := []byte(fmt.Sprintf(tailPosKeyFmt, topic))

	// Fetch the current tail position
	tailPosVal, err := tx.Get(tailPosKey, nil)
	if err != nil {
		return 0, fmt.Errorf("getting tail position from db: %v", err)
	}

	origOffset, err := binary.ReadVarint(bytes.NewReader(tailPosVal))
	if err != nil {
		return 0, fmt.Errorf("reading tail pos varint: %v", err)
	}

	// Write new record to next tail position
	newKey := []byte(fmt.Sprintf(keyFmt, topic, origOffset))

	if err := tx.Put(newKey, val, nil); err != nil {
		return 0, fmt.Errorf("putting value: %v", err)
	}

	// Update tail position
	tail := make([]byte, 8)
	binary.PutVarint(tail, origOffset+1)
	if err := tx.Put(tailPosKey, tail, nil); err != nil {
		return 0, fmt.Errorf("putting new tail position: %v", err)
	}

	return int(origOffset), nil
}

func prependValueTx(tx *leveldb.Transaction, headPosKeyFmt, keyFmt, topic string, val value) (offset int, err error) {
	headPosKey := []byte(fmt.Sprintf(headPosKeyFmt, topic))

	headPosVal, err := tx.Get(headPosKey, nil)
	if err != nil {
		return 0, fmt.Errorf("getting head position from db: %v", err)
	}

	headOffset, err := binary.ReadVarint(bytes.NewReader(headPosVal))
	if err != nil {
		return 0, fmt.Errorf("reading head pos varint: %v", err)
	}

	// Write new record to lower neighbouring position
	newHeadOffset := headOffset - 1
	newKey := []byte(fmt.Sprintf(keyFmt, topic, newHeadOffset))
	if err := tx.Put(newKey, val, nil); err != nil {
		return 0, fmt.Errorf("putting value: %v", err)
	}

	// Update head position
	_, newPosition, err := addPosTx(tx, headPosKeyFmt, topic, -1)
	if err != nil {
		return 0, fmt.Errorf("decrementing head pos by 1: %v", err)
	}

	return int(newPosition), nil
}
// addPos adds the an integer to a given position pointer
func addPos(db *leveldb.DB, posKeyFmt string, topic string, sum int) (oldPosition, newPosition int, err error) {
	oldPos, err := getPos(db, posKeyFmt, topic)
	if err != nil {
		return 0, 0, err
	}

	newPos := oldPos + sum
	newPosBytes := make([]byte, 8)
	binary.PutVarint(newPosBytes, int64(newPos))

	key := []byte(fmt.Sprintf(posKeyFmt, topic))

	if err := db.Put(key, newPosBytes, nil); err != nil {
		return 0, 0, fmt.Errorf("putting new increment position: %v", err)
	}

	return oldPos, newPos, nil
}

// addPosTx adds the an integer to a given position pointer.
func addPosTx(tx *leveldb.Transaction, posKeyFmt string, topic string, sum int) (oldPosition, newPosition int, err error) {
	oldPos, err := getPosTx(tx, posKeyFmt, topic)
	if err != nil {
		return 0, 0, err
	}

	newPos := oldPos + sum
	newPosBytes := make([]byte, 8)
	binary.PutVarint(newPosBytes, int64(newPos))

	key := []byte(fmt.Sprintf(posKeyFmt, topic))

	if err := tx.Put(key, newPosBytes, nil); err != nil {
		return 0, 0, fmt.Errorf("putting new increment position: %v", err)
	}

	return oldPos, newPos, nil
}
