package main

/*
 * Written: Dave Bevan
 * Date   : 23 Oct '17
 *
 * Purpose: To monitor a local folder, watching for file events using iNotify.
 * Specifically, trigger when CLOSE events have been received.
 *
 * Configured via 3 environment variables:
 *
 *   WATCHFOLDER=/path/to/be/monitored - can ONLY be a local folder, not network shares (as inotify events are absent from anything other than local paths)
 *
 *   %%PATH%% will be replaced with the full pathname to the file concerned:
 *
 *   E.g. WORKFLOWURL=http://something.url/some/endpoint.php?path=%%PATH%%
 *        WORKFLOWURL=http://something.url/rest/endpoint/action/%%PATH%%
 *
 *   FILESUFFIX=.xml - only call WORKFLOWURL for CLOSE events on files ending in $FILESUFFIX
 *
 */
import (
    "github.com/bevand10/dirmonitor/pkg/dirmonitor"
)

func main() {
    dirmonitor.Run()
}
