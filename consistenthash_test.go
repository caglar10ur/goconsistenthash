package consistenthash

import "testing"
import "strconv"
import "fmt"

func TestNew(t *testing.T) {
    numberOfReplicas := 32
    ch := New(numberOfReplicas)
    if ch == nil {
        t.Errorf("New failed")
    }
}

func TestAdd(t *testing.T) {
    numberOfReplicas := 32
    ch := New(numberOfReplicas)
    ch.Add("node1")
    if len(ch.Circle) != numberOfReplicas || len(ch.Hashes) != numberOfReplicas {
        t.Errorf("Circle or Hashes failed")
    }
    ch.Add("node2")
    if len(ch.Circle) != 2*numberOfReplicas || len(ch.Hashes) != 2*numberOfReplicas {
        t.Errorf("Circle or Hashes failed")
    }
}

func TestRemove(t *testing.T) {
    numberOfReplicas := 32
    ch := New(numberOfReplicas)
    ch.Add("node1")
    if len(ch.Circle) != numberOfReplicas || len(ch.Hashes) != numberOfReplicas {
        t.Errorf("Remove failed")
    }
    ch.Remove("node1")
    if len(ch.Circle) != 0 || len(ch.Hashes) != 0 {
        t.Errorf("Remove failed")
    }
}

func TestRemoveNonExisting(t *testing.T) {
    numberOfReplicas := 32
    ch := New(numberOfReplicas)
    ch.Add("node1")
    ch.Remove("node2")
    if len(ch.Circle) != numberOfReplicas || len(ch.Hashes) != numberOfReplicas {
        t.Errorf("RemoveNonExisting failed")
    }
}

func TestGet(t *testing.T) {
    numberOfReplicas := 32
    nodes := []string{"node1", "node2", "node3", "node4"}
    keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "ch", "y", "z"}

    m := make(map[string]int)

    distribution := func() {
        for k, v := range m {
            fmt.Println(k, v)
        }
    }

    ch := New(numberOfReplicas)
    for _, v := range nodes {
        ch.Add(v)
    }

    for _, v := range keys {
        m[ch.Get(v)]++
        if (v == "a" && ch.Get(v) != "node1") || (v == "l" && ch.Get(v) != "node2") || (v == "p" && ch.Get(v) != "node3") || (v == "w" && ch.Get(v) != "node4") {
            t.Errorf("Get failed")
        }
    }
    distribution()

    fmt.Println()
    ch.Add("node5")
    m = make(map[string]int)
    for _, v := range keys {
        m[ch.Get(v)]++
        if (v == "a" && ch.Get(v) != "node1") || (v == "l" && ch.Get(v) != "node2") || (v == "p" && ch.Get(v) != "node3") || (v == "w" && ch.Get(v) != "node5") {
            t.Errorf("Get failed")
        }
    }
    distribution()

    fmt.Println()
    ch.Remove("node5")
    m = make(map[string]int)
    for _, v := range keys {
        m[ch.Get(v)]++
        if (v == "a" && ch.Get(v) != "node1") || (v == "l" && ch.Get(v) != "node2") || (v == "p" && ch.Get(v) != "node3") || (v == "w" && ch.Get(v) != "node4") {
            t.Errorf("Get failed")
        }
    }
    distribution()
}

func BenchmarkCycle(b *testing.B) {
    ch := New(32)
    ch.Add("foo")
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ch.Add("foo" + strconv.Itoa(i))
        ch.Remove("foo" + strconv.Itoa(i))
    }
}

func BenchmarkCycleLarge(b *testing.B) {
    ch := New(32)
    for i := 0; i < 10; i++ {
        ch.Add("start" + strconv.Itoa(i))
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ch.Add("foo" + strconv.Itoa(i))
        ch.Remove("foo" + strconv.Itoa(i))
    }
}

func BenchmarkGet(b *testing.B) {
    ch := New(32)
    ch.Add("node")
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ch.Get("foo")
    }
}

func BenchmarkGetLarge(b *testing.B) {
    ch := New(32)
    for i := 0; i < 10; i++ {
        ch.Add("node" + ":" + strconv.Itoa(i))
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ch.Get("foo")
    }
}
