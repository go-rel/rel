# Instrumentation

REL provides hooks that can be used to log or instrument your queries.

{{ embed_code("examples/instrumentation.go", "instrumentation") }}

This is the list for available operations:

- `rel-aggregate`
- `rel-count`
- `rel-find`
- `rel-find-all`
- `rel-find-and-count-all`
- `rel-scan-one`
- `rel-scan-all`
- `rel-scan-multi`
- `rel-insert`
- `rel-insert-all`
- `rel-update`
- `rel-delete`
- `rel-delete-all`
- `rel-preload`
- `rel-transaction`
- `adapter-aggregate`
- `adapter-query`
- `adapter-exec`
- `adapter-begin`
- `adapter-commit`
- `adapter-rollback`
