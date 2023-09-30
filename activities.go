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
		return fmt.Errorf("failed to close fast rewards: %w", err)
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
	err := clickImage("buttons/back", 0.8)
	if err == nil {
		return nil
	}
	if err = clickXY(safePoint.X, safePoint.Y); err != nil { // skip rewards
		return err
	}

	return clickImage("buttons/back", 0.8)
}

func collectCompanionPoints() error {
	log.Info("Collecting companion points...")
	if err := clickImage("buttons/friends", 0.8); err != nil {
		return fmt.Errorf("failed to click on friends: %w", err)
	}
	if err := clickImage("buttons/sendandreceive", 0.8); err != nil {
		return fmt.Errorf("failed to click on send and receive: %w", err)
	}

	return clickImage("buttons/back", 0.8)
}

// def attemptCampaign():
//
//	printBlue('Attempting Campaign battle')
//	confirmLocation('campaign')
//	click('buttons/begin', seconds=2)
//	if (isVisible('buttons/begin', 0.7)): # If we see second Begin it's a multi so we take different actions
//	    click('buttons/begin', 0.7, seconds=2)
//	    click('buttons/beginbattle', seconds=4)
//	    click('buttons/pause', retry=3) # 3 retries as ulting heroes can cover the button
//	    click('buttons/exitbattle')
//	    click('buttons/back')
//	else: # else it's a single battle
//	    click('buttons/battle', 0.8, retry=3, seconds=3)
//	    click('buttons/battle_large', 0.8, suppress=True) #If you have no autobattle button its larger
//	    click('buttons/pause', 0.8, retry=3) # 3 retries as ulting heroes can cover the button
//	    click('buttons/exitbattle')
//	if confirmLocation('campaign', bool=True):
//	    printGreen('    Campaign attempted successfully')
//	else:
//	    printError('    Something went wrong, attempting to recover')
//	    recover()
func attemptCampaign() error {
	log.Info("Attempting Campaign battle")
	if err := clickImage("buttons/begin", 0.8); err != nil {
		return fmt.Errorf("failed to click on begin: %w", err)
	}
	err := waitUntilFoundAndClick(context.Background(), "img/buttons/begin_plain.png", 0.8, 10*time.Second)
	if err != nil {
		return err
	}
	if err = clickImage("buttons/beginbattle", 0.8); err != nil {
		return fmt.Errorf("failed to click on beginbattle: %w", err)
	}
	err = waitUntilFoundAndClick(context.Background(), "img/buttons/pause.png", 0.8, 10*time.Second)
	if err != nil {
		return err
	}
	if err = clickImage("buttons/exitbattle", 0.8); err != nil {
		return fmt.Errorf("failed to click on exitbattle: %w", err)
	}
	if err = clickImage("buttons/back", 0.8); err != nil {
		return fmt.Errorf("failed to click on back: %w", err)
	}
	return nil
}

// def handleBounties():
//
//	printBlue('Handling Bounty Board')
//	confirmLocation('darkforest')
//	clickXY(600, 1320)
//	if (isVisible('labels/bountyboard')):
//	    clickXY(650, 1700) # Solo tab
//	    click('buttons/collect_all', seconds=2, suppress=True)
//	    if config.getboolean('DAILIES', 'solobounties') is True:
//	        wait()
//	        click('buttons/dispatch', confidence=0.8, suppress=True, grayscale=True)
//	        click('buttons/confirm', suppress=True)
//	    clickXY(950,1700) # Team tab
//	    click('buttons/collect_all', seconds=2, suppress=True)
//	    click('buttons/dispatch', confidence=0.8, suppress=True, grayscale=True)
//	    click('buttons/confirm', suppress=True)
//	    click('buttons/back')
//	    printGreen('    Bounties attempted successfully')
//	else:
//	    printError('    Bounty Board not found, attempting to recover')
//	    recover()
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

// def handleArenaOfHeroes(count):
//
//	counter = 0
//	printBlue('Battling Arena of Heroes ' + str(count) + ' times')
//	confirmLocation('darkforest')
//	clickXY(740, 1050)
//	clickXY(550, 50)
//	if isVisible('labels/arenaofheroes_new'): # The label font changes for reasons
//	    click('labels/arenaofheroes_new', suppress=True)
//	    click('buttons/challenge', retry=3) # retries for animated button
//	    while counter < count:
//	        wait(1) # To avoid error when clickMultipleChoice returns no results
//	        clickMultipleChoice('buttons/arenafight', 4, confidence=0.98) # Select 4th opponent
//	        click('buttons/battle', 0.6, retry=3, suppress=True) # lower confidence as it's an animated button
//	        wait(2)
//	        click('buttons/skip', retry=5, suppress=True) # Retries as ulting heros can cover the button
//	        if (isVisible('labels/defeat')):
//	            printError('    Battle #' + str(counter+1) + ' Defeat!')
//	        else:
//	            printGreen('    Battle #' + str(counter+1) + ' Victory!')
//	            clickXY(600, 550) # Clear loot popup
//	        clickXY(600, 550)
//	        counter = counter+1
//	    click('buttons/exitmenu')
//	    click('buttons/back')
//	    click('buttons/back')
//	    printGreen('    Arena battles complete')
//	else:
//	    printError('Arena of Heroes not found, attempting to recover')
//	    recover()
func handleArenaOfHeroes(count int) error {
	counter := 0
	log.Info("Battling Arena of Heroes ", count, " times")
	if err := clickXY(740, 1312); err != nil {
		return fmt.Errorf("failed to click on arena of heroes: %w", err)
	}
	if err := clickXY(550, 1050); err != nil {
		return fmt.Errorf("failed to click on arena of heroes: %w", err)
	}
	if err := clickImage("labels/arenaofheroes_new", 0.8); err != nil {
		log.Error("No arena of heroes found")
	}
	if err := clickImage("buttons/challenge", 0.8); err != nil {
		return fmt.Errorf("failed to click on challenge: %w", err)
	}
	for counter < count {
		images, b := findAllInScreen("img/buttons/arenafight.png", 0.8)
		if !b {
			log.Info("No arena fight buttons found")
			return nil
		}
		if len(images) == 1 {
			break
		}
		weakestEnemy := getLowestImagePoint(images)
		err := clickXY(weakestEnemy.X, weakestEnemy.Y)
		if err != nil {
			return err
		}

		if err := clickImage("buttons/battle", 0.6); err != nil {
			return fmt.Errorf("failed to click on battle: %w", err)
		}
		time.Sleep(2 * time.Second)
		if err := clickImage("buttons/skip", 0.8); err != nil {
			log.Error("No skip button found")
		}
		time.Sleep(2 * time.Second)
		if err := clickXY(600, 687); err != nil {
			return fmt.Errorf("failed to click on loot popup: %w", err)
		}
		time.Sleep(2 * time.Second)
		if err := clickXY(600, 687); err != nil {
			return fmt.Errorf("failed to click on loot popup: %w", err)
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

// def collectGladiatorCoins():
//
//	printBlue('Collecting Gladiator Coins')
//	confirmLocation('darkforest')
//	clickXY(740, 1050)
//	clickXY(550, 50)
//	if isVisible('labels/legendstournament_new'): # The label font changes for reasons
//	    click('labels/legendstournament_new', suppress=True)
//	    clickXY(550, 300, seconds=2)
//	    clickXY(50, 1850)
//	    click('buttons/back')
//	    click('buttons/back')
//	    printGreen('    Gladiator Coins collected')
//	else:
//	    printError('    Legends Tournament not found, attempting to recover')
//	    recover()
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

// def collectInnGifts():
//
//	clicks = 0
//	x_axis = 250
//	printBlue('Attempting daily Inn gift collection')
//	confirmLocation('ranhorn')
//	clickXY(800,290, seconds=4)
//	if isVisible('buttons/manage'):
//	    while clicks < 10: # We spam clicks in the right area an d pray
//	        clickXY(x_axis, 1300, seconds=0.5)
//	        x_axis = x_axis + 50
//	        clicks = clicks + 1
//	        clickXY(550, 1400, seconds=0.5) # Clear loot
//	    click('buttons/back')
//	    printGreen('    Inn Gifts collected.')
//	else:
//	    printError('    Inn not found, attempting to recover')
//	    recover()
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
