package app

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nduyphuong/reverse-registry/config"
	"github.com/nduyphuong/reverse-registry/handler"
	"github.com/nduyphuong/reverse-registry/inject"
	digestfetcher "github.com/nduyphuong/reverse-registry/services/digest-fetcher"
	"github.com/nduyphuong/reverse-registry/utils"
	"github.com/sirupsen/logrus"
)

func RunAPI(conf config.Config) error {
	router := gin.Default()
	log := logrus.New()
	storage, err := inject.GetStorage(conf)
	if err != nil {
		return err
	}
	registryClient, err := inject.GetContainerRegistryClient()
	if err != nil {
		return err
	}
	handlerFactory := handler.New(handler.Options{
		Log:     log,
		Cr:      registryClient,
		Storage: storage,
	})

	router.Use(gin.WrapF(func(resp http.ResponseWriter, req *http.Request) {
		log.WithFields(logrus.Fields{
			"method": req.Method,
			"url":    req.URL.String(),
			"header": utils.Redact(req.Header),
		}).Info("app got request")
	}))
	router.Any("/v2", handlerFactory.V2Handler)
	router.Any("/v2/", handlerFactory.V2Handler)
	router.Any("/token", handlerFactory.TokenHandler)
	router.Any("/token/", handlerFactory.TokenHandler)
	router.Any("/v2/:repo/*rest", handlerFactory.ProxyHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		return err
	}
	return nil
}

func RunFetcher(conf config.Config) error {
	log := logrus.New()
	storage, err := inject.GetStorage(conf)
	if err != nil {
		return err
	}
	registryClient, err := inject.GetContainerRegistryClient()
	if err != nil {
		return err
	}
	d, err := time.ParseDuration(conf.WorkerFetchInterval)
	if err != nil {
		return err
	}
	fetcher := digestfetcher.New(digestfetcher.Options{
		Storage:       storage,
		Registry:      registryClient,
		Log:           log,
		FetchInterval: d,
	})
	return fetcher.Fetch(conf.Images)
}
