package main

import (
	"context"
	"image"
	_ "image/png"
	"time"

	"github.com/charmbracelet/log"
)

// config
var (
	maxX            = 1080
	maxY            = 2400
	fastRewardTimes = 5
	defaultWait     = 2 * time.Second
	arenaCount      = 1
	safePoint       = image.Point{X: maxX / 2, Y: 0}
	shopRefresh     = 3
)

func main() {
	//panic(buyTempleOfTime())
	var err error
	log.Info("Starting Afk Arena Bot")
	err = openAfkArena()
	if err != nil {
		log.Error("Failed to open Afk Arena", err)
		return
	}

	if err = claimAfkRewards(); err != nil {
		log.Error("Failed to claim afk rewards:", err)
		return
	}

	if err = claimFastRewards(fastRewardTimes); err != nil {
		log.Error("Failed to claim fast rewards:", err)
		return
	}

	if err = claimMail(); err != nil {
		log.Error("Failed to claim mail:", err)
		return
	}

	if err := collectCompanionPoints(); err != nil {
		log.Error("Failed to collect companion points:", err)
		return
	}

	if err = attemptCampaign(); err != nil {
		log.Error("Failed to attempt campaign:", err)
		return
	}

	err = waitUntilFoundAndClick(context.TODO(), "./img/buttons/darkforest_unselected.png", 0.8, 10*time.Second)
	if err != nil {
		log.Error("Failed to find darkforest_unselected:", err)
		return
	}

	if err = handleBounties(); err != nil {
		log.Error("Failed to handle bounties:", err)
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

	err = waitUntilFoundAndClick(context.TODO(), "./img/buttons/ranhorn_unselected.png", 0.8, 10*time.Second)
	if err != nil {
		log.Error("Failed to find ranhorn_unselected:", err)
		return
	}

	if err = collectInnGifts(); err != nil {
		log.Error("Failed to collect inn gifts:", err)
		return
	}

	if err = handleShopPurchase(shopRefresh); err != nil {
		log.Error("Failed to handle shop purchase:", err)
		return
	}

	if err = handleGuildHunts(); err != nil {
		log.Error("Failed to handle guild hunts:", err)
		return
	}

	if err = collectQuests(); err != nil {
		log.Error("Failed to collect quests:", err)
		return
	}

	log.Info("Afk Arena Bot finished successfully")
}
