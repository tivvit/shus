package main

import (
	"errors"
	"flag"
	"github.com/dgraph-io/badger"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"github.com/tivvit/shush/shush-api"
	"github.com/tivvit/shush/shush/backend"
	"github.com/tivvit/shush/shush/cache"
	"github.com/tivvit/shush/shush/config"
	backendConf "github.com/tivvit/shush/shush/config/backend"
	cacheConf "github.com/tivvit/shush/shush/config/cache"
	"github.com/valyala/fasthttp"
	"net/http"
	"sync"
)

var b cache.Cache

func main() {
	confFile := flag.String("confFile", "conf.yml", "Configuration file path")
	flag.Parse()
	c, err := config.NewConf(*confFile)
	if err != nil {
		log.Fatal(err)
	}
	setupLogger(c.Log)
	bck, err := initBackend(c.Backend)
	if err != nil {
		log.Fatal(err)
	}
	b, err = initCache(bck, c.Cache)
	if err != nil {
		log.Warn(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Infof("starting server at %s", c.Server.Address)
		err = fasthttp.ListenAndServe(c.Server.Address, fastHTTPHandler)
		if err != nil {
			log.Error(err)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		log.Printf("API Server starting at %s", c.Api.Address)
		log.Fatal(http.ListenAndServe(c.Api.Address, shush_api.NewRouter()))
		wg.Done()
	}()
	wg.Wait()
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	short := string(ctx.Path()[1:])
	url, err := b.Get(short)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	} else {
		ctx.Redirect(url, fasthttp.StatusFound)
	}
}

func setupLogger(c config.Log) {
	lvl, err := log.ParseLevel(c.Level)
	if err != nil {
		log.Error(err)
	}
	log.SetLevel(lvl)
}

func initBackend(bc backendConf.Conf) (backend.Backend, error) {
	if bc.InMem != nil {
		return backend.NewInMem(), nil
	}
	if bc.JsonFile != nil {
		return backend.NewJsonFile(bc.JsonFile.Path), nil
	}
	if bc.Redis != nil {
		return backend.NewRedis(&redis.Options{
			Addr: bc.Redis.Address,
		}), nil
	}
	if bc.Redis != nil {
		return backend.NewRedis(&redis.Options{
			Addr: bc.Redis.Address,
		}), nil
	}
	if bc.Badger != nil {
		return backend.NewBadger(badger.DefaultOptions(bc.Badger.Path)), nil
	}
	return nil, errors.New("unknown backend")
}

func initCache(b backend.Backend, cc *cacheConf.Conf) (cache.Cache, error) {
	if cc == nil {
		return b, nil
	}
	if cc.BigCache != nil {
		return cache.NewBigCache(b, cc.BigCache), nil
	}
	if cc.FreeCache != nil {
		return cache.NewFreeCache(b, cc.FreeCache), nil
	}
	if cc.LruCache != nil {
		return cache.NewLru(b, cc.LruCache), nil
	}
	if cc.FastCache != nil {
		return cache.NewFastCache(b, cc.FastCache), nil
	}
	if cc.RistrettoCache != nil {
		return cache.NewRistrettoCache(b, cc.RistrettoCache), nil
	}
	return b, errors.New("unknown cache")
}
