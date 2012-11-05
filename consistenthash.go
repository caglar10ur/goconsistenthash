// Reimplemantation of http://www.lexemetech.com/2007/11/consistent-hashing.html in Go

package consistenthash

import "fmt"
import "crypto/md5"
import "io"
import "strconv"
import "sort"
import "log"

const SEPERATOR = ":"

type ConsistentHash struct {
    //replicas indicates how many virtual points should be used per node and they are required to improve the distribution
    NumberOfReplicas    int
    // Circle
    Circle              map[string]string
    // Sorted list of hashes
    Hashes              []string
}

// helper function for generating md5 hash
func generateHash(input string) string {
    key := md5.New()
    _, err := io.WriteString(key, input)
    if err != nil {
        log.Fatal(err)
    }
    return fmt.Sprintf("%x", key.Sum(nil))
}

// returns an initialized ConsistentHash
func New(NumberOfReplicas int) *ConsistentHash {
    ch := new(ConsistentHash)
    ch.NumberOfReplicas = NumberOfReplicas
    ch.Circle = make(map[string]string)
    return ch
}

// adds a node to Circle
func (ch *ConsistentHash) Add(node string) {
    for i := 0; i < ch.NumberOfReplicas; i++ {
        // use SEPERATOR for protecting ourselves from collisions. If we concatenate directly then following happens;
        //
        // We have 20 nodes which are represented by strings as "node1"..."node20" with 100 replicas for each node
        ///
        // Following "virtual nodes" are generated "node1":
        // "node1" + 0 = "node10"
        // ...
        // "node1" + 10 = "node110"
        // "node1" + 11 = "node111"
        // ...
        //
        // At the same, following "virtual nodes" are generated for "node11":
        // "node11" + 0 = "node110"
        // "node11" + 1 = "node111"
        // ...
        hash := generateHash(node + SEPERATOR + strconv.Itoa(i))
        ch.Circle[hash] = node
        ch.Hashes = append(ch.Hashes, hash)
    }
    // keep Hashes sorted as we are going to use binary search on them
    sort.Strings(ch.Hashes)
}

// removes a node from Circle
func (ch *ConsistentHash) Remove(node string) {
    for i := 0; i < ch.NumberOfReplicas; i++ {
        hash := generateHash(node + ":" + strconv.Itoa(i))
        delete(ch.Circle, hash)
    }

    // As long as I see Go has append but not remove...
    ch.Hashes = ch.Hashes[:0]

    for k,_ := range ch.Circle {
        ch.Hashes = append(ch.Hashes, k)
    }
    // keep Hashes sorted as we are going to use binary search on them
    sort.Strings(ch.Hashes)
}

// returns the index of node that holds given key
func (ch *ConsistentHash) search(key string) (int) {
    // use binary search to find a proper node for fiven key
    index := sort.Search(len(ch.Hashes), func(i int) bool { return ch.Hashes[i] > key })

    if index == len(ch.Hashes) {
        index = 0
    }

    return index
}

// returns the node for given key
func (ch *ConsistentHash) Get(key string) (string) {
    index := ch.search(generateHash(key))
    return ch.Circle[ch.Hashes[index]]
}
