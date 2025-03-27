package timer

import (
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	var now = time.Now()
	var timer = AfterFunc(time.Second, func() {
		if time.Now().Sub(now) < time.Second {
			t.Fatal("Early func")
		}
	})
	timer.Start()
	time.Sleep(time.Second * 2)
}

func TestStop(t *testing.T) {
	var timer = AfterFunc(time.Second, func() {
		t.Fatal("Stopping is not working")
	})
	timer.Start()
	timer.Stop()
	time.Sleep(time.Second * 2)
}

func TestIdleStop(t *testing.T) {
	var timer = AfterFunc(time.Second, func() {
		t.Fatal("Stopping is not working")
	})
	timer.Start()
	timer.Stop()
	if timer.Stop() {
		t.Fatal("Invalid timer state")
	}
	if timer.Pause() {
		t.Fatal("Invalid timer state")
	}
	time.Sleep(time.Second * 2)
}

func TestPause(t *testing.T) {
	var now = time.Now()
	var timer = AfterFunc(time.Second, func() {
		if time.Now().Sub(now) < time.Second*2 {
			t.Fatal("Early func")
		}
	})
	timer.Start()
	time.Sleep(time.Microsecond * 500)
	timer.Pause()
	time.Sleep(time.Second)
	timer.Start()
	time.Sleep(time.Second)
}

func TestChanStart(t *testing.T) {
	var now = time.Now()
	var timer = NewTimer(time.Second)
	go func() {
		<-timer.C
		if time.Now().Sub(now) < time.Second {
			t.Fatal("Early func")
		}
	}()
	timer.Start()
	time.Sleep(time.Second * 2)
}

func TestChanStop(t *testing.T) {
	var timer = NewTimer(time.Second)
	go func() {
		<-timer.C
		t.Fatal("Stopping is not working")
	}()
	timer.Start()
	timer.Stop()
	time.Sleep(time.Second * 2)
}

func TestChanPause(t *testing.T) {
	var now = time.Now()
	var timer = NewTimer(time.Second)
	go func() {
		<-timer.C
		if time.Now().Sub(now) < time.Second*2 {
			t.Fatal("Early func")
		}
	}()
	timer.Start()
	time.Sleep(time.Microsecond * 500)
	timer.Pause()
	time.Sleep(time.Second)
	timer.Start()
	time.Sleep(time.Second)
}

func TestRestart(t *testing.T) {
	timer := NewTimer(time.Millisecond)
	c := make(chan bool)
	count := 0
	go func() {
		<-timer.C
		count++
		c <- true
		<-timer.C
		count++
		c <- true
	}()
	timer.Start()
	<-c
	timer.Restart(time.Millisecond)
	select {
	case <-c:
		if count != 2 {
			t.Fatalf("count = %d", count)
		}
	case <-time.After(time.Millisecond * 100):
		t.Fatal("timeout")
	}
}
