package main

import (
	"fmt"
	"image"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/log"
	"gocv.io/x/gocv"
)

func takeScreenshot() (string, error) {
	_, err := runADBCommand("shell", "screencap", "-p", "/sdcard/screenshot.png")
	if err != nil {
		return "", fmt.Errorf("failed to take screenshot: %w", err)
	}

	temp, err := os.CreateTemp(os.TempDir(), "afkarena-screenshot-*.png")
	if err != nil {
		return "", err
	}
	log.Debug("Created temp file", "path", temp.Name())
	defer temp.Close()

	// Pull the screenshot to local machine
	_, err = runADBCommand("pull", "/sdcard/screenshot.png", temp.Name())
	if err != nil {
		return "", fmt.Errorf("failed to pull screenshot: %w", err)
	}

	return temp.Name(), nil
}

func findTemplateInImage(mainImagePath, templateImagePath string, confidence float32) (*image.Point, bool) {
	// Read the main image
	mainImg := gocv.IMRead(mainImagePath, gocv.IMReadColor)
	if mainImg.Empty() {
		log.Error("Error reading main image")
		return nil, false
	}
	defer mainImg.Close()

	// Read the template image
	templateImg := gocv.IMRead(templateImagePath, gocv.IMReadColor)
	if templateImg.Empty() {
		log.Error("Error reading template image")
		return nil, false
	}
	defer templateImg.Close()

	// Perform template matching
	result := gocv.NewMat()
	defer result.Close()

	gocv.MatchTemplate(mainImg, templateImg, &result, gocv.TmCcoeffNormed, gocv.NewMat())

	// Find the maximum value and its location in the result matrix
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)

	log.Info("Template matching result", "maxVal", maxVal, "maxLoc", maxLoc, "template", templateImagePath)
	if maxVal >= confidence {
		// Calculate the middle point of the found template
		midX := maxLoc.X + templateImg.Cols()/2
		midY := maxLoc.Y + templateImg.Rows()/2
		return &image.Point{X: midX, Y: midY}, true
	}

	return nil, false
}

func clickXY(x, y int, wait time.Duration) error {
	if x > maxX || y > maxY {
		return fmt.Errorf("x and y must be less than %d and %d respectively", maxX, maxY)
	}
	_, err := runADBCommand("shell", "input", "tap", fmt.Sprintf("%d", x), fmt.Sprintf("%d", y))

	time.Sleep(wait)
	return err
}

func clickXYDefault(x, y int) error {
	return clickXY(x, y, defaultWait)
}

func clickImage(imagePath string, confidence float32) error {
	screenshotPath, err := takeScreenshot()
	if err != nil {
		return err
	}
	defer os.Remove(screenshotPath)

	point, found := findTemplateInImage(screenshotPath, fmt.Sprintf("./img/%s.png", imagePath), confidence)
	if !found {
		return fmt.Errorf("template not found in image: %s", imagePath)
	}

	return clickXY(point.X, point.Y, defaultWait)
}

func clickImageWithRetry(imagePath string, confidence float32, retry int) error {
	for i := 0; i < retry; i++ {
		if err := clickImage(imagePath, confidence); err != nil {
			time.Sleep(time.Second)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to click image: %s", imagePath)
}

func isVisible(imagePath string) (bool, error) {
	screenshotPath, err := takeScreenshot()
	if err != nil {
		return false, err
	}
	defer os.Remove(screenshotPath)

	_, found := findTemplateInImage(screenshotPath, fmt.Sprintf("./img/%s.png", imagePath), 0.80)
	return found, nil
}

func runADBCommand(args ...string) (string, error) {
	cmd := exec.Command("./platform-tools/adb", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func openAfkArena() error {
	log.Info("Opening AFK Arena...")
	if _, err := runADBCommand("shell", "monkey", "-p", "com.lilithgame.hgame.gp", "-c", "android.intent.category.LAUNCHER", "1"); err != nil {
		return err
	}

	if err := waitUntilGameLoaded(); err != nil {
		return fmt.Errorf("failed to wait until game loaded: %w", err)
	}

	log.Info("AFK Arena opened successfully!")

	log.Info("Expanding menus...")
	if err := expandMenus(); err != nil {
		return fmt.Errorf("failed to expand menus: %w", err)
	}
	log.Info("Menus expanded successfully!")

	return nil
}

func waitUntilGameLoaded() error {
	log.Info("Waiting for game to load...")
	// wait 1m with check every 2 seconds
	for i := 0; i < 30; i++ {
		visible, err := isVisible("buttons/campaign_selected")
		if err != nil {
			return err
		}
		if visible {
			return nil
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("game failed to load")
}

func clickButton(button string) error {
	err := clickImage(fmt.Sprintf("buttons/%s_selected", button), 0.8)
	if err != nil {
		return clickImage(fmt.Sprintf("buttons/%s_unselected", button), 0.8)
	}
	return nil
}

func confirmLocation(location string) error {
	log.Info("Confirming location...", "location", location)
	inLocation, err := isVisible(fmt.Sprintf("buttons/%s_selected", location))
	if err != nil {
		return fmt.Errorf("failed to check if in campaign: %w", err)
	}

	//click campaign if not in campaign until in campaign, max 5 times
	for i := 0; i < 5 && !inLocation; i++ {
		if err = clickButton(location); err != nil {
			return fmt.Errorf("failed to click button: %w", err)
		}
		inLocation, err = isVisible(fmt.Sprintf("buttons/%s_selected", location))
		if err != nil {
			return fmt.Errorf("failed to check if in campaign: %w", err)
		}
		time.Sleep(defaultWait)
	}

	if !inLocation {
		return fmt.Errorf("not in %s", location)
	}

	return nil
}

func expandMenus() error {
	for {
		visible, err := isVisible("buttons/downarrow")
		if err != nil {
			return err
		}
		if !visible {
			return nil
		}
		if err := clickImage("buttons/downarrow", 0.8); err != nil {
			return err
		}
		time.Sleep(defaultWait)
	}
}
