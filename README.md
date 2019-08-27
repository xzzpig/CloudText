# CloudText
An awesome tool to sync your clipboard

## Usage
### Server
```
Usage of CloudTextServer:
  -h string
        The http server host on. (default "0.0.0.0:23451")
  -u string
        username (default "cloudtext")
  -p string
        password (default random string)
  -t string
        http path (default "/cloudtext/text")
  -w string
        websocket path (default "/cloudtext/ws")
```

### Windows Client
```
Double click the CloudText.exe
Fill the config window
Enjoy!
```

### iOS Shortcut
1. Download the .shortcut file from `release`
2. Import it to into shortcut app
3. Edit it to config something
 * URL : change it to your CloudTextServer
 * username : the `-u` argument you set in CloudTextServer
 * password : the `-p` argument you set in CloudTextServer
4. Enjoy!
