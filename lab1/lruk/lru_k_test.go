package lruk

import (
	"errors"
	"testing"
	"time"
)

func TestLRUKNode_SetEvictable(t *testing.T) {
	node := LRUKNode{}
	node.SetEvictable(true)
	if node.isEvictable != true {
		t.Error("evictable should be true")
	}
}

func TestNewLRUKReplacer(t *testing.T) {
	lruKReplacer := NewLRUKReplacer(10, 10)
	if lruKReplacer == nil {
		t.Error("lruKReplacer should not be nil")
	}
}

func TestLRUKReplacer_Evict(t *testing.T) {
	// assume numFrames is 5 and k is 3
	numFrames, k := 5, 3

	t.Run("uninitialized", func(t *testing.T) {
		lruKReplacer := new(LRUKReplacer)
		err := lruKReplacer.Evict(0)
		if !errors.Is(err, ErrUnInitialized) {
			t.Error("should return ErrUnInitialized")
		}
	})

	t.Run("invalid frame id", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)
		err := lruKReplacer.Evict(-1)
		if !errors.Is(err, ErrInvalidFrameId) {
			t.Error("should return ErrInvalidFrameId")
		}

		err = lruKReplacer.Evict(numFrames)
		if !errors.Is(err, ErrInvalidFrameId) {
			t.Error("should return ErrInvalidFrameId")
		}
	})

	t.Run("no evictable frame", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)
		err := lruKReplacer.Evict(0)
		if !errors.Is(err, ErrNoEvictableFrame) {
			t.Error("should return ErrNoEvictableFrame")
		}
	})

	t.Run("no enough data", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)
		err := lruKReplacer.RecordAccess(0, 0)
		if err != nil {
			t.Error("should no error")
		}
		err = lruKReplacer.Evict(0)
		if !errors.Is(err, ErrNoEvictableFrame) {
			t.Error("should return ErrNoEvictableFrame")
		}
	})

	t.Run("normal case", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)
		for i := 0; i < numFrames; i++ {
			for j := 0; j < k; j++ {
				err := lruKReplacer.RecordAccess(i, 0)
				if err != nil {
					t.Error(err)
				}
			}
			err := lruKReplacer.SetEvictable(i, true)
			if err != nil {
				t.Error(err)
			}
		}

		// sleep for a second in case the timestamp conflict
		time.Sleep(1 * time.Second)
		err := lruKReplacer.Evict(0)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestLRUKReplacer_RecordAccess(t *testing.T) {
	// assume numFrames is 5 and k is 3
	numFrames, k := 5, 3

	t.Run("uninitialized", func(t *testing.T) {
		lruKReplacer := new(LRUKReplacer)
		err := lruKReplacer.RecordAccess(0, 0)
		if !errors.Is(err, ErrUnInitialized) {
			t.Error("should return ErrUnInitialized")
		}
	})

	t.Run("invalid frame id", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)

		err := lruKReplacer.RecordAccess(-1, 0)
		if !errors.Is(err, ErrInvalidFrameId) {
			t.Error("should return ErrInvalidFrameId")
		}
	})

	t.Run("invalid access type", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)

		err := lruKReplacer.RecordAccess(0, -1)
		if !errors.Is(err, ErrUnknownAccessType) {
			t.Error("should return ErrUnknownAccessType")
		}
	})

	t.Run("access times greater than k", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)

		for i := 0; i < k; i++ {
			err := lruKReplacer.RecordAccess(0, 0)
			if err != nil {
				t.Error(err)
			}
		}

		err := lruKReplacer.RecordAccess(0, 0)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestLRUKReplacer_SetEvictable(t *testing.T) {
	numFrames, k := 5, 3

	t.Run("uninitialized", func(t *testing.T) {
		lruKReplacer := new(LRUKReplacer)
		err := lruKReplacer.SetEvictable(0, true)
		if !errors.Is(err, ErrUnInitialized) {
			t.Error("should return ErrUnInitialized")
		}
	})

	t.Run("invalid frame id", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)
		err := lruKReplacer.SetEvictable(-1, true)
		if !errors.Is(err, ErrInvalidFrameId) {
			t.Error("should return ErrInvalidFrameId")
		}

		err = lruKReplacer.SetEvictable(0, true)
		if !errors.Is(err, ErrInvalidFrameId) {
			t.Error("should return ErrInvalidFrameId")
		}
	})

	t.Run("normal case", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)

		err := lruKReplacer.RecordAccess(0, 0)
		if err != nil {
			t.Error(err)
		}

		err = lruKReplacer.SetEvictable(0, true)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestLRUKReplacer_Remove(t *testing.T) {
	numFrames, k := 5, 3

	t.Run("uninitialized", func(t *testing.T) {
		lruKReplacer := new(LRUKReplacer)
		err := lruKReplacer.Remove(0)
		if !errors.Is(err, ErrUnInitialized) {
			t.Error("should return ErrUnInitialized")
		}
	})

	t.Run("invalid frame id", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)
		err := lruKReplacer.Remove(-1)
		if !errors.Is(err, ErrInvalidFrameId) {
			t.Error("should return ErrInvalidFrameId")
		}
	})

	t.Run("non-evictable frame", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)

		err := lruKReplacer.RecordAccess(0, 0)
		if err != nil {
			t.Error(err)
		}

		err = lruKReplacer.Remove(0)
		if !errors.Is(err, ErrUnRemovableFrame) {
			t.Error("should return ErrNoEvictableFrame")
		}
	})

	t.Run("normal case", func(t *testing.T) {
		lruKReplacer := NewLRUKReplacer(numFrames, k)

		err := lruKReplacer.RecordAccess(0, 0)
		if err != nil {
			t.Error(err)
		}

		err = lruKReplacer.SetEvictable(0, true)
		if err != nil {
			t.Error(err)
		}

		err = lruKReplacer.Remove(0)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestLRUKReplacer_Size(t *testing.T) {
	numFrames, k := 5, 3

	lruKReplacer := NewLRUKReplacer(numFrames, k)

	if lruKReplacer.Size() != 0 {
		t.Error("lruKReplacer.Size() should be 0")
	}

	err := lruKReplacer.RecordAccess(0, 0)
	if err != nil {
		t.Error(err)
	}

	if lruKReplacer.Size() != 1 {
		t.Error("lruKReplacer.Size() should be 1")
	}
}
