package main

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

func claimAfkRewards() error {
	log.Info("Claiming afk rewards...")
	if err := clickXY(maxX/2, int(float64(maxY)*0.8)); err != nil {
		return fmt.Errorf("failed to click on afk rewards: %w", err)
	}

	return waitUntilFoundAndClick(context.TODO(), "img/buttons/collect.png", 0.8, 10*time.Second)
}

func claimFastRewards(times int) error {
	log.Info("Claiming fast rewards...")

	if err := waitUntilFoundAndClick(context.TODO(), "img/buttons/fastrewards.png", 0.8, 10*time.Second); err != nil {
		return fmt.Errorf("failed to click on fast rewards: %w", err)
	}
	for i := 0; i < times; i++ {
		// Submit
		if err := clickXY(700, 1500); err != nil {
			return fmt.Errorf("failed to submit fast rewards: %w", err)
		}
		time.Sleep(defaultWait)
		// Claim fast rewards
		if err := clickXY(932, 2140); err != nil {
			return fmt.Errorf("failed to claim fast rewards: %w", err)
		}
	}

	if err := clickImage("buttons/close", 0.8); err != nil {
		_ = clickXY(safePoint.X, safePoint.Y)
		return clickImage("buttons/close", 0.8)
	}

	return nil
}

func claimMail() error {
	log.Info("Claiming mail...")
	if err := clickImage("buttons/mail", 0.8); err != nil {
		return fmt.Errorf("failed to click on fast rewards: %w", err)
	}
	if err := clickImage("buttons/collect_all", 0.8); err != nil {
		return fmt.Errorf("failed to click on collect: %w", err)
	}

	// click back until we are back at the main screen
	for {
		_ = clickXY(safePoint.X, safePoint.Y)
		if err := clickImage("buttons/back", 0.8); err != nil {
			break
		}
	}

	return nil
}

func collectCompanionPoints() error {
	log.Info("Collecting companion points...")
	if err := clickImage("buttons/friends", 0.8); err != nil {
		return fmt.Errorf("failed to click on friends: %w", err)
	}
	if err := clickImage("buttons/sendandreceive", 0.8); err != nil {
		return fmt.Errorf("failed to click on send and receive: %w", err)
	}
	err := clickImage("buttons/back", 0.8)
	if err != nil {
		log.Info("No back button found")
	}
	err = clickImage("buttons/back", 0.8)
	if err != nil {
		log.Info("No back button found")
	}

	return nil
}

func attemptCampaign() error {
	log.Info("Attempting Campaign battle")
	if err := clickImage("buttons/begin", 0.8); err != nil {
		return fmt.Errorf("failed to click on begin: %w", err)
	}
	err := waitUntilFoundAndClick(context.Background(), "img/buttons/battle.png", 0.8, 10*time.Second)
	if err != nil {
		return err
	}
	//if err = clickImage("buttons/beginbattle", 0.8); err != nil {
	//	return fmt.Errorf("failed to click on beginbattle: %w", err)
	//}
	err = waitUntilFoundAndClick(context.Background(), "img/buttons/pause.png", 0.8, 10*time.Second)
	if err != nil {
		return err
	}
	if err = clickImage("buttons/exitbattle", 0.8); err != nil {
		return fmt.Errorf("failed to click on exitbattle: %w", err)
	}
	//if err = clickImage("buttons/back", 0.8); err != nil {
	//	return fmt.Errorf("failed to click on back: %w", err)
	//}
	return nil
}

func handleBounties() error {
	log.Info("Handling Bounty Board")
	if err := clickImage("buttons/bountyboard", 0.8); err != nil {
		return fmt.Errorf("failed to click on bounty board: %w", err)
	}
	//if err := clickXY(650, 1700); err != nil {
	//	return fmt.Errorf("failed to click on solo tab: %w", err)
	//}
	if err := clickImage("buttons/collect_all", 0.8); err != nil {
		log.Error("No solo bounties to collect")
	}
	if err := clickImage("buttons/dispatch2", 0.8); err != nil {
		log.Error("No solo bounties to dispatch")
	}
	if err := clickImage("buttons/confirm", 0.8); err != nil {
		log.Error("No solo bounties to confirm")
	}
	if err := clickImage("buttons/teambounty", 0.8); err != nil {
		return fmt.Errorf("failed to click on team tab: %w", err)
	}
	if err := clickImage("buttons/collect_all", 0.8); err != nil {
		log.Error("No team bounties to collect")
	}
	if err := clickImage("buttons/dispatch2", 0.8); err != nil {
		log.Error("No team bounties to dispatch")
	}
	if err := clickImage("buttons/confirm", 0.8); err != nil {
		log.Error("No team bounties to confirm")
	}
	if err := clickImage("buttons/back", 0.8); err != nil {
		return fmt.Errorf("failed to click on back: %w", err)
	}

	return nil
}

func handleArenaOfHeroes(count int) error {
	counter := 0
	log.Info("Battling Arena of Heroes ", count, " times")
	if err := clickXY(740, 1312); err != nil {
		return fmt.Errorf("failed to click on arena of heroes: %w", err)
	}
	if err := clickXY(550, 1050); err != nil {
		return fmt.Errorf("failed to click on arena of heroes: %w", err)
	}

	if err := waitUntilFoundAndClick(context.Background(), "img/labels/rewards.png", 0.7, 5*time.Second); err != nil {
		log.Error("No rewards labels found")
	}

	if err := clickImage("labels/arenaofheroes_new", 0.8); err != nil {
		log.Error("No arena of heroes found")
	}
	if err := clickImage("buttons/challenge", 0.8); err != nil {
		return fmt.Errorf("failed to click on challenge: %w", err)
	}
	for counter < count {
		_, err := waitUntilFound(context.TODO(), "img/buttons/arenafight.png", 0.8, 5*time.Second)
		if err != nil {
			return fmt.Errorf("failed to find arena fight: %w", err)
		}
		images, b := findAllInScreen("img/buttons/arenafight.png", 0.8)
		if !b {
			log.Info("No arena fight buttons found")
			return nil
		}
		if len(images) == 1 {
			break
		}
		weakestEnemy := getLowestImagePoint(images)
		err = clickXY(weakestEnemy.X, weakestEnemy.Y)
		if err != nil {
			return err
		}

		if err := clickImage("buttons/battle", 0.59); err != nil {
			return fmt.Errorf("failed to click on battle: %w", err)
		}
		time.Sleep(2 * time.Second)
		if err := clickImage("buttons/skip", 0.8); err != nil {
			log.Error("No skip button found")
		}
		if err = waitUntilFoundAndClick(context.Background(), "img/labels/rewards.png", 0.7, 10*time.Second); err != nil {
			log.Error("No rewards labels found")
		}
		if err = waitUntilFoundAndClick(context.Background(), "img/labels/taptocontinue.png", 0.6, 10*time.Second); err != nil {
			log.Error("No tap to continue labels found")
		}

		counter = counter + 1
	}
	if err := clickImage("buttons/exitmenu", 0.8); err != nil {
		return fmt.Errorf("failed to click on exitmenu: %w", err)
	}
	if err := clickImage("buttons/back", 0.8); err != nil {
		return fmt.Errorf("failed to click on back: %w", err)
	}
	if err := clickImage("buttons/back", 0.8); err != nil {
		return fmt.Errorf("failed to click on back: %w", err)
	}

	return nil
}

func collectGladiatorCoins() error {
	log.Info("Collecting Gladiator Coins")
	if err := clickXY(740, 1312); err != nil {
		return fmt.Errorf("failed to click on arena of heroes: %w", err)
	}
	//if err := clickXY(550, 1050); err != nil {
	//	return fmt.Errorf("failed to click on arena of heroes: %w", err)
	//}
	if err := clickImage("labels/legendstournament_new", 0.8); err != nil {
		return fmt.Errorf("failed to click on legends tournament: %w", err)
	}
	time.Sleep(2 * time.Second)
	if err := clickXY(550, 375); err != nil {
		return fmt.Errorf("failed to click on legends tournament: %w", err)
	}
	if err := clickXY(safePoint.X, safePoint.Y); err != nil { // skip rewards
		return err
	}
	if err := clickImage("buttons/back", 0.8); err != nil {
		return fmt.Errorf("failed to click on back: %w", err)
	}
	if err := clickImage("buttons/back", 0.8); err != nil {
		return fmt.Errorf("failed to click on back: %w", err)
	}
	if err := clickImage("buttons/back", 0.8); err != nil {
		return fmt.Errorf("failed to click on back: %w", err)
	}

	return nil
}

func collectInnGifts() error {
	clicks := 0
	xAxis := 250
	log.Info("Attempting daily Inn gift collection")
	time.Sleep(2 * time.Second)
	if err := clickXY(800, 550); err != nil {
		return fmt.Errorf("failed to click on inn: %w", err)
	}
	_, err := waitUntilFound(context.Background(), "img/buttons/manage.png", 0.8, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to find manage: %w", err)
	}

	for clicks < 10 {
		time.Sleep(500 * time.Millisecond)
		if err := clickXY(xAxis, 1800); err != nil {
			return fmt.Errorf("failed to click on inn: %w", err)
		}
		xAxis = xAxis + 50
		clicks = clicks + 1
		time.Sleep(500 * time.Millisecond)
		if err := clickXY(550, 1750); err != nil {
			return fmt.Errorf("failed to click on inn: %w", err)
		}
	}
	if err := clickImage("buttons/back", 0.8); err != nil {
		return fmt.Errorf("failed to click on back: %w", err)
	}

	return nil
}

func handleShopPurchase(refreshCount int) error {
	log.Info("Handling shop purchase")
	if err := clickImage("buttons/shop/shop", 0.8); err != nil {
		return fmt.Errorf("failed to click on shop: %w", err)
	}

	for i := 0; i < refreshCount; i++ {
		if err := waitUntilFoundAndClick(context.TODO(), "img/buttons/shop/quickbuy.jpg", 0.8, 5*time.Second); err != nil {
			return fmt.Errorf("failed to click on quick buy: %w", err)
		}
		if err := waitUntilFoundAndClick(context.TODO(), "img/buttons/shop/purchase.png", 0.8, 5*time.Second); err != nil {
			return fmt.Errorf("failed to click on purchase: %w", err)
		}
		if err := waitUntilFoundAndClick(context.Background(), "img/labels/rewards.png", 0.8, 5*time.Second); err != nil {
			log.Error("No rewards labels found")
		}
		_ = clickXY(safePoint.X, safePoint.Y)
		if i == refreshCount-1 {
			break
		}
		if err := waitUntilFoundAndClick(context.Background(), "img/buttons/shop/refresh.jpg", 0.8, 5*time.Second); err != nil {
			log.Error("No rewards labels found")
		}
		if err := clickImage("buttons/confirm", 0.8); err != nil {
			return fmt.Errorf("failed to click on shop: %w", err)
		}
	}

	if err := clickImage("buttons/back", 0.8); err != nil {
		return fmt.Errorf("failed to click on back: %w", err)
	}

	return nil
}

func handleGuildHunts() error {
	log.Info("Attempting to run Guild Hunts")
	if err := clickImage("buttons/guild", 0.8); err != nil {
		return fmt.Errorf("failed to click on guild: %w", err)
	}

	if err := waitUntilFoundAndClick(context.Background(), "img/buttons/fortune_chest.png", 0.7, 10*time.Second); err != nil {
		log.Error("No fortune chest found")
	}
	if err := clickImage("buttons/exitmenu", 0.8); err != nil {
		fmt.Println("No exit menu found")
	}

	if err := clickImage("buttons/guildhunting", 0.8); err != nil {
		return fmt.Errorf("failed to click on guild: %w", err)
	}

	if err := clickImage("buttons/quickbattle", 0.8); err != nil {
		log.Error("No quick battle found")
	}
	if err := clickImage("buttons/sweep", 0.8); err != nil {
		log.Error("No quick battle found")
	}
	if err := clickImage("buttons/confirm", 0.8); err != nil {
		log.Error("No confirm button found")
	}

	if err := clickXY(safePoint.X, safePoint.Y); err != nil {
		log.Error("No quick battle found")
	}

	if err := clickXY(safePoint.X, safePoint.Y); err != nil {
		log.Error("No quick battle found")
	}

	if err := clickImage("buttons/arrow_right", 0.8); err != nil {
		log.Error("No quick battle found")
	}

	if err := clickImage("buttons/quickbattle", 0.8); err != nil {
		log.Error("No quick battle found")
	}
	if err := clickImage("buttons/sweep", 0.8); err != nil {
		log.Error("No quick battle found")
	}

	if err := clickImage("buttons/confirm", 0.8); err != nil {
		log.Error("No confirm button found")
	}

	for i := 0; i < 5; i++ {
		_ = clickXY(safePoint.X, safePoint.Y)
		_, found := findInScreen("img/buttons/back.png", 0.8)
		if found {
			break
		}

	}

	for {
		if err := clickImage("buttons/back", 0.8); err != nil {
			break
		}
	}
	return nil
}

func collectQuests() error {
	log.Info("Attempting to collect quests")
	if err := clickImage("buttons/quest", 0.8); err != nil {
		return fmt.Errorf("failed to click on quests: %w", err)
	}
	if err := clickImage("buttons/collect", 0.7); err != nil {
		return fmt.Errorf("failed to click on collect: %w", err)
	}
	if err := clickImage("buttons/fullquestchest", 0.8); err != nil {
		return fmt.Errorf("failed to click on fullquestchest: %w", err)
	}
	time.Sleep(2 * time.Second)
	_ = clickXY(safePoint.X, safePoint.Y)

	return clickImage("buttons/back", 0.8)
}

func buyTempleOfTime() error {
	log.Info("Attempting to buy Temple of Time")

	time.Sleep(5 * time.Second)

	for {
		if err := clickImage("buttons/summon", 0.8); err != nil {
			return fmt.Errorf("failed to click on summon: %w", err)
		}
		if err := clickImage("buttons/card", 0.8); err != nil {
			return fmt.Errorf("failed to click on card: %w", err)
		}
		if err := clickImage("buttons/back", 0.8); err != nil {
			return fmt.Errorf("failed to click on card: %w", err)
		}
	}
}
