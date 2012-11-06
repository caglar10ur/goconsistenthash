# Consistent Hash

Reimplemantation of (http://www.lexemetech.com/2007/11/consistent-hashing.html "Consistent Hashing by Tom White") in Go

## Consistent hashing (from the web page mentioned above)

Consistent hashing is a scheme that provides a hash table functionality in a way that the adding or removing of one slot does not significantly change the mapping of keys to slots.

The need for consistent hashing arose from limitations experienced while running collections of caching machines - web caches, for example. If you have a collection of n cache machines then a common way of load balancing across them is to put object o in cache machine number hash(o) mod n. This works well until you add or remove cache machines (for whatever reason), for then n changes and every object is hashed to a new location. This can be catastrophic since the originating content servers are swamped with requests from the cache machines. It's as if the cache suddenly disappeared. Which it has, in a sense. (This is why you should care - consistent hashing is needed to avoid swamping your servers!)

It would be nice if, when a cache machine was added, it took its fair share of objects from all the other cache machines. Equally, when a cache machine was removed, it would be nice if its objects were shared between the remaining machines. This is exactly what consistent hashing does - consistently maps objects to the same cache machine, as far as is possible, at least.

The basic idea behind the consistent hashing algorithm is to hash both objects and caches using the same hash function. The reason to do this is to map the cache to an interval, which will contain a number of object hashes. If the cache is removed then its interval is taken over by a cache with an adjacent interval.

## Status

[![Build Status](https://secure.travis-ci.org/caglar10ur/goconsistenthash.png)](http://travis-ci.org/caglar10ur/goconsistenthash)
