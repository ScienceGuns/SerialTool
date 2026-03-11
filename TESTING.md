Testing
-------
In order to test Serial Tool, you will need to create a config file at `~/.config/scienceguns/serialtool.json` as such:

```json
{
  "data_url": "https://scienceguns.com/data/",
  "mongodb": {
    "uri": "mongodb://drkrieger:guest@localhost:27017",
    "database": "firearms",
    "collection": "Serials"
  }
}
```
Then start Docker using the [docker-compose](testing/docker-compose.yml) file in the [testing](testing) folder.
This is required as the encoding functionality relies on a database in order to track counts.

Note
----
Currently Serial Tool is build for `darwin-arm64` and `linux-amd64`, however interactive testing is only performed on the macOS variant.
