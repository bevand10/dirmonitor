package main

/*
 * Written: Dave Bevan
 * Date   : 23 Oct '17
 *
 * Purpose: To monitor a local folder, watching for file events using iNotify.
 * Specifically, look for CLOSE events (CLOSE_WRITE is monitored).
 *
 * Configured via 3 environment variables:
 *
 *   WATCHFOLDER=/path/to/be/monitored - can ONLY be a local folder, not network shares (as inotify events are absent from anything other than local paths)
 *
 *   %%PATH%% will be replaced with the full pathname to the file concerned:
 *
 *   E.g. WORKFLOWURL=http://something.url/some/endpoint?path=%%PATH%%
 *        WORKFLOWURL=http://something.url/rest/endpoint/action/%%PATH%%
 *
 *   FILESUFFIX=.xml - only call WORKFLOWURL for CLOSE events on files ending in $FILESUFFIX
 *
 */
import (
	"github.com/betim/fsnotify"
	"github.com/yookoala/realpath"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	watchfolder, ok := os.LookupEnv("WATCHFOLDER")
	if !ok {
		log.Fatal("Env WATCHFOLDER unset!")
	}

	watchfolder, err := realpath.Realpath(watchfolder)
	if err != nil {
		log.Fatal(error(err))
	}

	workflowurl, ok := os.LookupEnv("WORKFLOWURL")
	if !ok {
		log.Fatal("Env WORKFLOWURL unset!\n",
			"  E.g. WORKFLOWURL=http://something.url/some/endpoint.php?path=%%PATH%%\n",
			"       WORKFLOWURL=http://something.url/rest/endpoint/action/%%PATH%%")
	}

	filesuffix, ok := os.LookupEnv("FILESUFFIX")
	if !ok {
		log.Fatal("Env FILESUFFIX unset!")
	}

	startWatching(watchfolder, workflowurl, filesuffix)
}

func startWatching(watchfolder, workflowurl string, filesuffix string) {

	log.Println("WATCHFOLDER=" + watchfolder)
	log.Println("WORKFLOWURL=" + workflowurl)
	log.Println("FILESUFFIX=" + filesuffix)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				procEvent(event, workflowurl, filesuffix)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(watchfolder)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}

func procEvent(event fsnotify.Event, workflowurl string, filesuffix string) {

	if event.Op&fsnotify.Close == fsnotify.Close {

		log.Println("CLOSE:", event.Name)

		if strings.HasSuffix(event.Name, filesuffix) {

			p := strings.LastIndex(event.Name, filesuffix)
			wfurl := strings.Replace(workflowurl, "%%PATH%%", url.PathEscape(event.Name[:p]), -1)
			if strings.Contains(wfurl, "?") {
				wfurl += "&FILESUFFIX=" + filesuffix
			} else {
				wfurl += "?FILESUFFIX=" + filesuffix
			}
			log.Println("WORKFLOW:", wfurl)

			var client = &http.Client{
				Timeout: time.Second * 5,
			}

			retval, err := client.Get(wfurl)
			if err != nil {
				log.Println("Error:", err)
			}

			defer retval.Body.Close()
			body, err := ioutil.ReadAll(retval.Body)
			if err != nil {
				log.Println("Error:", err)
			}

			log.Println("RESPONSE:", string(body))
		}
	}
}
