package main

import (
	"log"
	"os"

	"github.com/mschneider82/sharecmd/clipboard"
	"github.com/mschneider82/sharecmd/config"
	"github.com/mschneider82/sharecmd/provider"
	"github.com/mschneider82/sharecmd/provider/dropbox"
	"github.com/mschneider82/sharecmd/provider/googledrive"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFile = kingpin.Flag("config", "Client configuration file").Default(config.UserHomeDir() + "/.config/sharecmd/config.json").String()
	setup      = kingpin.Flag("setup", "Setup client configuration").Bool()
	file       = kingpin.Arg("file", "filename to upload").File()
)

// ShareCmd cli app
type ShareCmd struct {
	config   *config.Config
	provider provider.Provider
}

func main() {
	kingpin.Parse()

	if *setup {
		config.Setup(*configFile)
		os.Exit(0)
	}
	if file != nil {
		sharecmd := ShareCmd{}
		cfg, err := config.LookupConfig(*configFile)
		if err != nil {
			log.Fatalf("lookupConfig: %v \n", err)
		}
		sharecmd.config = &cfg

		switch sharecmd.config.Provider {
		case "googledrive":
			sharecmd.provider = googledrive.NewGoogleDriveProvider(sharecmd.config.ProviderSettings["googletoken"])
		case "dropbox":
			sharecmd.provider = dropbox.NewDropboxProvider(cfg.ProviderSettings["token"])
		default:
			config.Setup(*configFile)
			os.Exit(0)
		}

		fileid, err := sharecmd.provider.Upload(*file, "")
		if err != nil {
			log.Fatalf("Can't upload file: %s", err.Error())
		}
		link, err := sharecmd.provider.GetLink(fileid)
		if err != nil {
			log.Fatalf("Can't get link for file: %s", err.Error())
		}
		log.Printf("URL: %s", link)
		clipboard.ToClip(link)
	}
}
