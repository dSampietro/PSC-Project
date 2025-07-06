# Impl detail
- map do not need any mutex: even if more goroutines operate on the structure, if ops are read-only then it is thread-safe
- graph nodes id is the position of the word in the text; the same word appear multiple times (every node is unique)