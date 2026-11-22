---
name: go-correctness
description: Runtime bugs in the asoai OpenAI client
applies_to: ["**/*.go"]
---
asoai is an OpenAI chat client with sessions and a history database
(buntdb). Review the changed Go code for bugs that show up when it runs:

- Errors ignored or replaced, so a failed request or write looks like
  success; `err` shadowed inside a block and never checked.
- Nil dereferences and out-of-range indexing on API responses (an empty
  `Choices`), on missing sessions, or on flags left unset.
- Streams: `io.EOF` treated as an error or other errors as the end, streams
  and responses not closed, partial output lost.
- buntdb: writes in a read-only transaction (`View`), reads and writes of a
  session that are not in one transaction when they must be, keys built
  inconsistently (`session:<name>`).
- Sessions: history appended twice or lost, the current session pointing at
  one that was deleted.

Only report what is visible in the diff and the files given.
