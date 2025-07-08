# Sentence
If the sentences in the text are independent, the generated sentences wil not interleave.

# Impl detail
- map do not need any mutex: even if more goroutines operate on the structure, if ops are read-only then it is thread-safe
- graph nodes id is the position of the word in the text; the same word appear multiple times (every node is unique)
- to handle cycles, in a sentence generation, a node may not appear more than N (parameter) times; otherwise the generation will be aborted; $N \ge max_i\{|str_i|\}$ 
  doing this, we are not excluding valid long sentences, but we impose an hard limit on (possibly) infinite cycles.