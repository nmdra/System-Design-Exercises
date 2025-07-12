package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/lmittmann/tint"
)

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelInfo})))

	idStr := os.Getenv("NODE_ID")
	peerStr := os.Getenv("PEERS")
	redisAddr := os.Getenv("REDIS_ADDR")

	if idStr == "" || peerStr == "" || redisAddr == "" {
		slog.Error("Missing NODE_ID, PEERS, or REDIS_ADDR")
		os.Exit(1)
	}

	id, _ := strconv.Atoi(idStr)
	var allIDs []int
	for _, s := range strings.Split(peerStr, ",") {
		n, _ := strconv.Atoi(s)
		allIDs = append(allIDs, n)
	}

	node := NewNode(id, allIDs, redisAddr)
	node.Start()

	// https://stackoverflow.com/a/48769120
	select {}
}
