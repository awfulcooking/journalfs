package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/util"

	"github.com/togetherbeer/journalfs/journalcache"
	"github.com/togetherbeer/journalfs/mount"
)

var mountPath = flag.String("p", "/var/log/journalfs", "mount path")
var allowOther = flag.Bool("allowOther", false, "allow other users to access the filesystem. user_allow_other must be enabled in /etc/fuse.conf to use this option without being root")
var follow = flag.Bool("follow", true, "Follow the journal for new entries after startup")
var debug = flag.Bool("debug", false, "print debugging messages")

var buildTime, version = "unknown", "unknown"

func init() {
	flag.Parse()
}

func mountOptions() []mount.MountOption {
	var options []mount.MountOption
	if *allowOther {
		options = append(options, mount.AllowOther)
	}
	return options
}

func main() {
	fmt.Printf("journalfs version %s (build time %s)\n", version, buildTime)
	fmt.Println()

	if runningAsService, _ := util.RunningFromSystemService(); runningAsService {
		if *debug {
			fmt.Println("Running as a systemd service. Cannot enable debug logging as it will cause a journalling loop")
			*debug = false
		}
	}

	journalCache, count, err := loadJournalCache()
	must(err)

	fmt.Printf("Loaded %d entries.\n", count)

	if *follow {
		fmt.Println("Following.")
		journalCache.Follow()
	}

	journalCache.Debug = *debug

	mount := mount.NewMount(*mountPath, journalCache)

	go func() {
		must(mount.Serve(mountOptions()...))
	}()

	fmt.Println()
	fmt.Println("Serving", *mountPath)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	fmt.Println()
	fmt.Println("Got stop signal. Unmounting..")

	unmountErr := mount.Unmount()
	fmt.Println(unmountErr)
}

func loadJournalCache() (*journalcache.JournalCache, int, error) {
	jc, err := journalcache.NewJournalCache()

	if err != nil {
		return nil, 0, err
	}

	count, err := jc.Load()
	return jc, count, err
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func printJSON(v interface{}) error {
	if s, err := toJSONString(v); err != nil {
		return err
	} else {
		fmt.Println(s)
		return nil
	}
}

func toJSONString(v interface{}) (string, error) {
	if bytes, err := json.MarshalIndent(v, "", "  "); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}
