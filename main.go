package main

import (
	"flag"
	"fmt"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"log"
	"os"
)

var (
	opts *Opts
)

type Opts struct {
	Scene     string
	Item      string
	StartRec  bool
	StopRec   bool
	ToggleRec bool
	Host      string
	Port      string
	Password  string
}

func NewOpts() *Opts {
	opts := &Opts{}

	flag.StringVar(&opts.Scene, "scene", lookupEnv("SCENE"), "Set/change to scene")
	flag.StringVar(&opts.Item, "item", lookupEnv("ITEM"), "Toggle scene item")
	flag.BoolVar(&opts.StartRec, "start-rec", lookupEnv("START_RECORD") != "", "Start recording")
	flag.BoolVar(&opts.StopRec, "stop-rec", lookupEnv("STOP_RECORD") != "", "Stop recording")
	flag.BoolVar(&opts.ToggleRec, "rec", lookupEnv("RECORD") != "", "Toggle recording")
	flag.StringVar(&opts.Host, "host", lookupEnv("HOST", "localhost"), "Host to connect to")
	flag.StringVar(&opts.Port, "port", lookupEnv("PORT", "4455"), "Port to connect to")
	flag.StringVar(&opts.Password, "password", lookupEnv("PASSWORD"), "Password")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if opts.Password == "" {
		fmt.Println("Password is mandatory")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return opts
}

func lookupEnv(key string, defaultValues ...string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	for _, v := range defaultValues {
		if v != "" {
			return v
		}
	}
	return ""
}

func main() {
	opts = NewOpts()
	client, err := goobs.New(opts.Host+":"+opts.Port, goobs.WithPassword(opts.Password))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	//version, err := client.General.GetVersion()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("OBS Studio version: %s\n", version.ObsVersion)
	//fmt.Printf("Websocket server version: %s\n", version.ObsWebSocketVersion)

	if opts.Scene != "" && opts.Item == "" {
		p := &scenes.SetCurrentProgramSceneParams{SceneName: opts.Scene}
		_, err := client.Scenes.SetCurrentProgramScene(p)
		if err != nil {
			fmt.Printf("cannot set current scene: %v", err)
		}
	} else if opts.StartRec {
		_, err := client.Record.StartRecord()
		if err != nil {
			fmt.Printf("cannot start recording: %v", err)
		}
	} else if opts.StopRec {
		_, err := client.Record.StopRecord()
		if err != nil {
			fmt.Printf("cannot stop recording: %v", err)
		}
	} else if opts.ToggleRec {
		_, err := client.Record.ToggleRecord()
		if err != nil {
			fmt.Printf("cannot toggle recording: %v", err)
		}
	} else if opts.Scene != "" && opts.Item != "" {
		p := &sceneitems.GetSceneItemIdParams{
			SceneName:    opts.Scene,
			SourceName:   opts.Item,
			SearchOffset: 0,
		}
		r, err := client.SceneItems.GetSceneItemId(p)
		if err != nil {
			fmt.Printf("cannot get scene %v", err)
		}
		p1 := &sceneitems.GetSceneItemEnabledParams{
			SceneItemId: r.SceneItemId,
			SceneName:   opts.Scene,
		}
		r1, err := client.SceneItems.GetSceneItemEnabled(p1)
		if err != nil {
			fmt.Printf("cannot get mode", err)
		}
		b2 := !r1.SceneItemEnabled
		p2 := &sceneitems.SetSceneItemEnabledParams{
			SceneItemEnabled: &b2,
			SceneName:        opts.Scene,
			SceneItemId:      r.SceneItemId,
		}
		_, err = client.SceneItems.SetSceneItemEnabled(p2)
		if err != nil {
			fmt.Printf("cannot toggle item: %v", err)
		}
	} else {
		fmt.Println("not sure what to do...")
		flag.PrintDefaults()
	}
	client.Disconnect()
}
