package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cosmotek/loguago"
	"github.com/jessevdk/go-flags"
	luajson "github.com/layeh/gopher-json"
	"github.com/rs/zerolog"
	lua "github.com/yuin/gopher-lua"
	"roycetechnology.com/luacsv/assets"
	"roycetechnology.com/luacsv/pkg/utils"
)

var (
	Version = "0.0.0"
	Build   = "-"
)

var opts struct {
	LuaFile string `long:"file" short:"f" description:"Lua file entry" default:"main.lua"`

	LogLevel string `long:"log-level" choice:"debug" choice:"info" choice:"warn" choice:"error" default:"info"`

	Version func() `long:"version" short:"v" description:"Show bulid version"`
}

var logger *loguago.Logger
var parser = flags.NewParser(&opts, flags.Default)

func init() {
	opts.Version = func() {
		fmt.Printf("Version: %v", Version)
		fmt.Printf("\tBuild: %v", Build)
		os.Exit(0)
	}
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
	os.Mkdir("logs", os.ModePerm)
	ts := time.Now().UTC().Format(time.RFC3339)
	ts = strings.Replace(strings.Replace(ts, ":", "", -1), "-", "", -1)
	logger = utils.NewFileLogger("logs/log." + ts + ".log")
	loglevel, _ := zerolog.ParseLevel(opts.LogLevel)
	zerolog.SetGlobalLevel(loglevel)
}

func main() {
	files := make([]string, 0)
	visit := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if filepath.Ext(path) != ".csv" {
			return nil
		}
		files = append(files, path)
		return nil
	}
	filepath.WalkDir(".", visit)

	if len(files) == 0 {
		log.Print("no files")
	}

	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		r := csv.NewReader(f)
		headers, err := r.Read()
		if err == io.EOF {
			log.Fatal(err)
		}

		records, err := r.ReadAll()
		if err == io.EOF {
			log.Fatal(err)
		}

		// Initialize lua vm
		vm := lua.NewState(lua.Options{
			MinimizeStackMemory: true,
		})
		defer vm.Close()
		for _, pair := range []struct {
			n string
			f lua.LGFunction
		}{
			{lua.LoadLibName, lua.OpenPackage}, // Must be first
			{lua.BaseLibName, lua.OpenBase},
			{lua.TabLibName, lua.OpenTable},
		} {
			if err := vm.CallByParam(lua.P{
				Fn:      vm.NewFunction(pair.f),
				NRet:    0,
				Protect: true,
			}, lua.LString(pair.n)); err != nil {
				panic(err)
			}
		}
		vm.PreloadModule("logger", logger.Loader)
		luajson.Preload(vm)

		// Assign DATA
		array := vm.NewTable()
		for _, record := range records {
			table := vm.NewTable()
			for i, value := range record {
				lvalue := lua.LString(value)
				header := strings.TrimSpace(headers[i])
				table.RawSet(lua.LString(header), lvalue)
			}
			array.Append(table)
		}
		vm.SetGlobal("DATA", array)

		script, _ := assets.Asset("assets/main.lua")
		if f, err := os.Open(opts.LuaFile); err == nil {
			if script, err = ioutil.ReadAll(f); err != nil {
				log.Fatal(err)
			}
		}

		// Load lua file
		err = vm.DoString(string(script))
		if err != nil {
			log.Print(path, err)
		}
		log.Printf("%v parsed.", path)
	}
}
