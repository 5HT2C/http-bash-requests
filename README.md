# http-bash-requests

A simple wrapper to make bash (or other shell) commands with http.

**DISCLAIMER:**
This intentionally has no authentication built in, and only works on localhost. As such, be VERY cautious about giving it elevated permissions.
This also allows **other users** on your machine to delete your entire home folder (even without elevated permissions), if they would like, among other things.

I do not take any responsibility for the consequences of you running this, or any other, software on any computer.

**OTHER DISCLAIMER:**
This also has the ability to execute *any* executable file on your computer as your current user.
The `X-Bin-Path`, `X-Bin-Arg` and `X-Body-Split` HTTP headers make this possible. This is very dangerous.

I really do not advise you run this on any open machine that isn't a scrappable VM (READ: DO NOT RUN ON A REAL SYSTEM).

## Why?

I wanted a very quick hack to allow my docker containers to reboot themselves.

## Usage

```bash
# Run
go run main.go -port 6016

# Build to dir
go build -o ~/.local/bin/http-bash-requests .

curl localhost:6016 -d "echo test"
```

## Service

You can also use the services to keep it enabled.

```bash
mkdir -p ~/.config/systemd/user/
curl https://raw.githubusercontent.com/5HT2C/http-bash-requests/master/http-bash-requests.service -o ~/.config/systemd/user/http-bash-requests.service
curl https://raw.githubusercontent.com/5HT2C/http-bash-requests/master/http-bash-requests.timer -o ~/.config/systemd/user/http-bash-requests.timer
# Enable the service and timer for the current user
systemctl --user enable --now http-bash-requests.timer
```

Verify that the service is working like so, and ensure "test" is in the status log:

![](https://raw.githubusercontent.com/5HT2C/http-bash-requests/master/img.png)

## Library

Here is a basic example for using the library:

```go
package main

import (
	"github.com/5HT2C/http-bash-requests/httpBashRequests"
	"log"
	"net/http"
)

func main() {
	// Setup only needed once
	client := httpBashRequests.Client{Addr: "http://localhost:6016", HttpClient: &http.Client{Timeout: 5 * time.Minute}}
	httpBashRequests.Setup(&client)

	// Now we can run bash requests over http
	log.Println(httpBashRequests.Run("ls"))
}
```
