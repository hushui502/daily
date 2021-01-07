# hsfs
implementation of io/fs.FS that appends SHA256 hashes to filename to allow for aggressive HTTP caching.

note that this library requires Go 1.16 or higher, you can download 1.16beta, and add io/fs package to old version.


## Usage
To use hsfs, first wrap your fs in a hsfs.FS filesystem.
```cassandraql
var embedFS embed.FS

var fsys = hashfs.NewFS(embedFS)
```
Then attach a hsfs.FileServer to your router.
```cassandraql
http.Handle("/assets", http.StripPrefix("/assets", hashfs.FileServer(fsys)))
```
Next, your html templating library can obtain the hashname of your file using the hsfs.FS.HashName() method.
```cassandraql
func renderHTML(w io.Writer) {
	fmt.Fprintf(w, `<html>`)
	fmt.Fprintf(w, `<script script="/assets/%s">`, fsys.HashName("scripts/main.js"))
	fmt.Fprintf(w, `</html>`)
}
```

