# logrus-gcloud-format

Logrus formatter compatible with Google Cloud's Stackdriver logging agent

## Installation
To install formatter, use `go get`:

```sh
$ go get github.com/x-cray/logrus-prefixed-formatter
```

## Usage
Here is how it should be used:

```go
package main

import (
	"github.com/Sirupsen/logrus"
	lgcloud "github.com/everflow-io/logrus-gcloud-format"
)

var log = logrus.New()

func init() {
	log.Formatter = &lgcloud.GCloudFormatter{}
	log.Level = logrus.DebugLevel
}

```
