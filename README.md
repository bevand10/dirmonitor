# dirmonitor
Go iNotify based directory monitor with http notification and support for UTF-8 filenames.

## Dependencies
https://golang.org/doc/go1.8 - specifically net/url - url.PathEscape()

## Usage
dirmonitor is configured via 3 environment variables. In the examples below, `%%PATH%%` will be replaced with the full pathname of the target file stripped of the supplied suffix match filter.

env | role | example
--- | --- | ---
WATCHFOLDER | The directory you want to monitor for changes. Note that network mounts (e.g. CIFS, NFS etc) do NOT generate iNotify events. | /path/to/be/monitored
WORKFLOWURL | A URL to be called when CLOSE events have been detected on matched files (see below) in the monitored directory. | http://something.url/some/endpoing.php?path=%%PATH%%
||| http://something.url/rest/endpoint/action/%%PATH%%
FILESUFFIX | A filename match filter string, not regex. | .xml

For example:

```
#!/bin/sh
export WATCHFOLDER=/path/to/be/monitored
export WORKFLOWURL=http://something.url/inotify.php?path=%%PATH%%
export FILESUFFIX=.xml
/usr/sbin/dirmonitor
```

In all circumstances, the param `FILESUFFIX` will be appended to the constructed URL.

Given a `CLOSE` event observed on `/path/to/be/monitored/105_8339634_5290303والمعلومات.xml`, the following URL will be called:

`http://octojex.labs.jupiter.bbc.co.uk/api2/jupiterarrival/%2Fhome%2Focto_from_jupiter%2F105_8339634_5290303%D9%88%D8%A7%D9%84%D9%85%D8%B9%D9%84%D9%88%D9%85%D8%A7%D8%AA?FILESUFFIX=.xml`
