package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dc-dc-dc/cheetah/market/research"
)

var (
	cache = flag.String("cache", "./data/research", "The cache directory for research")
)

func init() {
	flag.Parse()
}

func main() {
	if _, err := os.Stat(*cache); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(*cache, os.ModePerm); err != nil {
				log.Fatalf("[error] failed to create directory=%s, err=%s", *cache, err)
			}
		} else {
			log.Fatalf("[error] got unknown err=%s", err)
		}
	}

	listDir := filepath.Join(*cache, "lists")
	if _, err := os.Stat(listDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(listDir, os.ModePerm); err != nil {
				log.Fatalf("[error] failed to create directory=%s err=%s", listDir, err)
			}
		} else {
			log.Fatalf("[error] got unknown err=%s", err)
		}
	}

	listClient := research.NewStockListResearch()
	items := []research.StockListOption{research.StockListNasdaqExchange, research.StockListNyseExchange}
	for _, item := range items {
		res, err := listClient.GetList(context.Background(), item)
		if err != nil {
			fmt.Printf("[error] failed to fetch option=%s err=%s\n", item, err)
			continue
		}
		fmt.Printf("[info] got %d items from option=%s\n", len(res), item)
		filename := fmt.Sprintf("%s.json", strings.ToLower(strings.ReplaceAll(item.String(), "-", "_")))
		fd, err := os.Create(filepath.Join(listDir, filename))
		if err != nil {
			fmt.Printf("[error] trying to create file=%s, err=%s\n", filename, err)
			continue
		}
		defer fd.Close()
		if err := json.NewEncoder(fd).Encode(res); err != nil {
			fmt.Printf("[error] writing json filename=%s err=%s", filename, err)
		}
	}
}
