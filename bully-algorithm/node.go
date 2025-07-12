package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	ID                 int
	AllNodeIDs         []int
	LeaderID           int
	ElectionInProgress bool

	receivedOK     bool
	sentOK         bool
	lastElectionAt time.Time
	electionTimer  *time.Timer

	Redis *RedisHelper
}

func NewNode(id int, allIDs []int, redisAddr string) *Node {
	return &Node{
		ID:         id,
		AllNodeIDs: allIDs,
		Redis:      NewRedisHelper("election", redisAddr, id),
	}
}

func (n *Node) Start() {
	go n.Redis.Subscribe("election", n.HandleMessage)
	go n.Redis.Subscribe("heartbeat", n.HandleHeartbeat)

	slog.Info("Node started", slog.Int("id", n.ID))

	go n.StartHeartbeat()
	go n.MonitorLeader()
	go n.PrintLeader()

	time.Sleep(time.Duration(n.ID) * time.Second)
	n.TriggerElection()
}

func (n *Node) TriggerElection() {
	if n.ElectionInProgress {
		return
	}

	if time.Since(n.lastElectionAt) < 2*time.Second {
		// Avoid too frequent elections
		return
	}
	n.lastElectionAt = time.Now()

	n.ElectionInProgress = true
	n.receivedOK = false
	n.sentOK = false

	slog.Info("Starting election", slog.Int("node", n.ID))

	hasHigher := false
	for _, peerID := range n.AllNodeIDs {
		if peerID > n.ID {
			hasHigher = true
			n.Redis.Publish("election", fmt.Sprintf("ELECTION:%d", n.ID))
		}
	}

	if !hasHigher {
		n.becomeLeaderAndAnnounce()
		return
	}

	if n.electionTimer != nil {
		n.electionTimer.Stop()
	}
	n.electionTimer = time.AfterFunc(3*time.Second, func() {
		if !n.receivedOK {
			n.becomeLeaderAndAnnounce()
		} else {
			// Wait for coordinator announcement after receiving OK
			time.AfterFunc(5*time.Second, func() {
				if n.ElectionInProgress {
					slog.Warn("Coordinator timeout, retrying election", slog.Int("node", n.ID))
					n.ElectionInProgress = false
					n.TriggerElection()
				}
			})
		}
	})
}

func (n *Node) HandleMessage(msg string) {
	switch {
	case strings.HasPrefix(msg, "ELECTION:"):
		sender := parseID(msg)
		if sender < n.ID {
			slog.Info("Received ELECTION", slog.Int("me", n.ID), slog.Int("from", sender))
			if !n.sentOK {
				n.Redis.Publish("election", fmt.Sprintf("OK:%d", n.ID))
				n.sentOK = true
			}
			n.TriggerElection()
		}
	case strings.HasPrefix(msg, "OK:"):
		slog.Info("Received OK", slog.Int("node", n.ID))
		n.receivedOK = true
	case strings.HasPrefix(msg, "COORDINATOR:"):
		leader := parseID(msg)
		slog.Info("New leader announced", slog.Int("me", n.ID), slog.Int("leader", leader))
		n.LeaderID = leader
		n.ElectionInProgress = false
		n.sentOK = false
		if n.electionTimer != nil {
			n.electionTimer.Stop()
		}
		n.Redis.RefreshHeartbeatTimer()
	}
}

func (n *Node) HandleHeartbeat(msg string) {
	if strings.HasPrefix(msg, "HEARTBEAT:") {
		leader := parseID(msg)
		if leader == n.LeaderID {
			n.Redis.RefreshHeartbeatTimer()
		}
	}
}

func (n *Node) StartHeartbeat() {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		if n.LeaderID == n.ID {
			n.Redis.Publish("heartbeat", fmt.Sprintf("HEARTBEAT:%d", n.ID))
		}
	}
}

func (n *Node) MonitorLeader() {
	for range n.Redis.HeartbeatTimeout() {
		slog.Warn("Heartbeat timeout – assuming leader is down", slog.Int("me", n.ID))

		// Wait a short randomized duration before triggering election
		delay := time.Duration(500+rand.Intn(1000)) * time.Millisecond
		time.AfterFunc(delay, func() {
			n.TriggerElection()
		})
	}
}

func (n *Node) becomeLeaderAndAnnounce() {
	slog.Info("🎉 Became Leader", slog.Int("node", n.ID))
	n.LeaderID = n.ID
	n.ElectionInProgress = false
	n.Redis.Publish("election", fmt.Sprintf("COORDINATOR:%d", n.ID))
	n.Redis.RefreshHeartbeatTimer()
}

func (n *Node) PrintLeader() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if n.LeaderID == n.ID {
			slog.Info("I'm leader", slog.Int("node", n.ID))
		}
	}
}

func parseID(msg string) int {
	parts := strings.Split(msg, ":")
	id, _ := strconv.Atoi(parts[1])
	return id
}
