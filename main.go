package main

import (
	"fmt"
	_ "image/png"
	"time"

	"github.com/charmbracelet/log"
)

// config
var (
	maxX            = 1080
	maxY            = 2400
	fastRewardTimes = 1
	defaultWait     = 2 * time.Second
	arenaCount      = 5
)

func main() {
	log.Info("Starting Afk Arena Bot")
	err := openAfkArena()
	if err != nil {
		log.Error("Failed to open Afk Arena", err)
		return
	}

	if err = claimAfkRewards(); err != nil {
		fmt.Println("Failed to claim afk rewards:", err)
		return
	}

	if err = claimFastRewards(fastRewardTimes); err != nil {
		fmt.Println("Failed to claim fast rewards:", err)
		return
	}

	if err = claimMail(); err != nil {
		fmt.Println("Failed to claim mail:", err)
		return
	}

	if err := collectCompanionPoints(); err != nil {
		fmt.Println("Failed to collect companion points:", err)
		return
	}

	if err = attemptCampaign(); err != nil {
		fmt.Println("Failed to attempt campaign:", err)
		return
	}

	if err = handleBounties(); err != nil {
		fmt.Println("Failed to handle bounties:", err)
		return
	}

	if err = handleArenaOfHeroes(arenaCount); err != nil {
		log.Error("Failed to handle arena of heroes:", err)
		return
	}

	if err = collectGladiatorCoins(); err != nil {
		log.Error("Failed to collect gladiator coins:", err)
		return
	}

	if err = collectInnGifts(); err != nil {
		log.Error("Failed to collect inn gifts:", err)
		return
	}

	log.Info("Afk Arena Bot finished successfully")
}
