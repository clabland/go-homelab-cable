# go-homelab-cable

A Go application for playing media in your homelab the way a television station would. Set up multiple "channels" within your network, each of which continuously play through a list of local media files at random.

## Dependencies
Only tested on linux machines, requires VLC and LibVLC.
```
sudo apt install vlc libvlc-dev
``` 

## Jump Start
### Start the Server
1. Sync this repo
1. `go run cmds/cli/main server --path "./local/path/to/folder/of/videos"`

### Interact with the Server via the Client
#### View what's currently playing
`go run cmds/cli/main client`

### Interact with the Server via the web UI
Navigate to http://localhost:3004/ in your browser.

### Interact with the Server via the HTTP API
**Request**
`curl http://localhost:3004/api/networks`
**Response**
```
[
    {
        "name": "Homelab Cable",
        "call_sign": "KHLC"
    }
]
```

**Request**
`curl http://localhost:3004/api/:call_sign/channels`
**Response**
[
    {
        "id": "bd63bc6a-27cd-47fb-81d5-642a216f335e",
        "playing": "title of currently playing media.mp4",
        "up_next": "title of the next show which will play.mp4",
        "live": true // whether or not this channel is outputing via VLC on the device
    },
    {
        "id": "0cd43939-082b-4c15-b642-0682a035cfb0",
        "playing": "some other currently playing media.mp4",
        "up_next": "some other show which will play next.mp4",
        "live": false // if false, the channel is a 'virtual' channel. it is cycling through its media list in memory every 30 minutes to mimic a channel that isn't being watched.
    }
]