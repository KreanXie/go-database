package lruk

import (
	"container/list"
	"sync"
	"time"
)

type Node struct {
	history     *list.List
	k           int
	frameId     int
	isEvictable bool
}

func (l *Node) SetEvictable(setEvictable bool) {
	l.isEvictable = setEvictable
}

type Replacer struct {
	nodeStore    map[int]*Node
	timeStamp    int
	curSize      int
	replacerSize int
	k            int
	mu           *sync.Mutex
}

func NewReplacer(numFrames, k int) *Replacer {
	return &Replacer{
		nodeStore:    make(map[int]*Node),
		timeStamp:    int(time.Now().Unix()),
		curSize:      0,
		replacerSize: numFrames,
		k:            k,
		mu:           &sync.Mutex{},
	}
}

func (lru *Replacer) Evict(frameId int) (int, error) {
	// error handle
	switch {
	case lru.nodeStore == nil:
		// uninitialized lru replacer

		return -1, ErrUnInitialized
	case frameId < 0 || frameId >= lru.replacerSize:
		// frameId should not be greater than replacerSize

		return -1, ErrInvalidFrameId
	}

	lru.mu.Lock()
	defer lru.mu.Unlock()

	earliestTime := time.Now().Unix()
	earliestFrameId := -1

	// find the frame with largest kth backward distance
	for _, node := range lru.nodeStore {
		if node.history.Len() < lru.k {
			continue
		}
		kthBackwardTime := node.history.Front().Value.(int64)
		if node.isEvictable && kthBackwardTime < earliestTime {
			earliestTime = kthBackwardTime
			earliestFrameId = node.frameId
		}
	}

	if earliestFrameId == -1 {
		return -1, ErrNoEvictableFrame
	}

	delete(lru.nodeStore, earliestFrameId)
	lru.curSize--

	return earliestFrameId, nil
}

func (lru *Replacer) RecordAccess(frameId int, accessType int) error {
	// error handle
	switch {
	case lru.nodeStore == nil:
		// uninitialized lru replacer

		return ErrUnInitialized
	case frameId < 0 || frameId >= lru.replacerSize:
		// frameId should not be greater than replacerSize

		return ErrInvalidFrameId
	case accessType < 0 || accessType > 3:
		// accessType should be in [0,3], but it's not used here anyway

		return ErrUnknownAccessType
	default:
	}

	lru.mu.Lock()
	defer lru.mu.Unlock()

	// if this frame is not seen.
	if _, ok := lru.nodeStore[frameId]; !ok {
		lru.nodeStore[frameId] = &Node{
			history:     list.New(),
			k:           lru.k,
			frameId:     frameId,
			isEvictable: false,
		}
	}

	history := lru.nodeStore[frameId].history
	// if length of history is greater than k, then pop latest one record
	if history.Len() == lru.k {
		_ = history.Remove(lru.nodeStore[frameId].history.Front())
	}

	_ = history.PushBack(time.Now().Unix())
	lru.curSize++

	return nil
}

func (lru *Replacer) SetEvictable(frameId int, setEvictable bool) error {
	// error handle
	switch {
	case lru.nodeStore == nil:
		// uninitialized lru replacer

		return ErrUnInitialized
	case frameId < 0 || frameId >= lru.replacerSize:
		// frameId should not be greater than replacerSize

		return ErrInvalidFrameId
	default:
	}

	lru.mu.Lock()
	defer lru.mu.Unlock()

	if _, ok := lru.nodeStore[frameId]; !ok {
		return ErrInvalidFrameId
	}

	lru.nodeStore[frameId].SetEvictable(setEvictable)

	return nil
}

func (lru *Replacer) Remove(frameId int) error {
	// error handle
	switch {
	case lru.nodeStore == nil:
		// uninitialized lru replacer

		return ErrUnInitialized
	case frameId < 0 || frameId >= lru.replacerSize:
		// frameId should not be greater than replacerSize

		return ErrInvalidFrameId
	case !lru.nodeStore[frameId].isEvictable:
		// un evictable frame

		return ErrUnRemovableFrame
	default:
	}

	lru.mu.Lock()
	defer lru.mu.Unlock()

	delete(lru.nodeStore, frameId)
	lru.curSize--

	return nil
}

func (lru *Replacer) Size() int {
	return lru.curSize
}
