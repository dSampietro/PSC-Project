# Sentence
We define a sentence as a path fro a node to the final node "."
If the sentences in the text are independent, the generated sentences wil not interleave.

# Impl detail
- map do not need any mutex: even if more goroutines operate on the structure, if ops are read-only then it is thread-safe
- graph nodes id is the position of the word in the text; the same word appear multiple times (every node is unique)
- to handle cycles, in a sentence generation,  the generated sentence must not be longer than $max_depth$ (param); otherwise that generation will be aborted; 
  doing this, we are not excluding valid long sentences, but we impose an hard limit on (possibly) infinite cycles.

# Scalability
The first solution had a fixed buffer size; however 
1) for bigger texts
2) for varying depth, limit paramters
this can cause panics or crashes.

My solution was to opt for unbounded channels. (implemented as In channel, queue as buffer, Out channel)