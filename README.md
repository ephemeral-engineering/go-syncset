# Go SyncSet

SyncSet is a golang Set implementation using sync.Map

## Getting Started

### Installation

```sh
go get github.com/ephemeral-engineering/go-syncset
```

### Usage

```go
package main

import (
  "log"
  "github.com/ephemeral-engineering/go-syncset"
)

func main() {
    syncSet := syncset.NewSyncSet[string]()
    syncSet.Add("item1")
    syncSet.Add("item2")
    syncSet.Add("item3")

    if syncSet.Size() != 3 {
        log.Errorf("Expected 3 items, got %d", syncSet.Size())
    }
}
```

## Make a new version

```sh
go test -coverprofile=cover.out
go tool cover -html cover.out -o cover.html
firefox cover.html
git commit ...
git tag "v1.1.0"
git push --tags
```
