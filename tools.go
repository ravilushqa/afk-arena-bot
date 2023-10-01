package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"os/exec"
	"time"

	"github.com/charmbracelet/log"
	"gocv.io/x/gocv"
)

func runADBCommand(args ...string) (string, error) {
	cmd := exec.Command("./platform-tools/adb", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func captureScreen() (image.Image, error) {
	cmd := exec.Command("./platform-tools/adb", "exec-out", "screencap", "-p")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	img, err := png.Decode(&out)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func findInScreen(templateImagePath string, confidence float32) (*image.Point, bool) {
	screen, err := captureScreen()
	if err != nil {
		fmt.Println("Error capturing screen:", err)
		return nil, false
	}
	// Read the main image
	mainImg, err := gocv.ImageToMatRGB(screen)
	if err != nil {
		fmt.Println("Error reading main image")
		return nil, false
	}
	defer mainImg.Close()

	// Read the template image
	templateImg := gocv.IMRead(templateImagePath, gocv.IMReadColor)
	if templateImg.Empty() {
		fmt.Println("Error reading template image")
		return nil, false
	}
	defer templateImg.Close()

	// Perform template matching
	result := gocv.NewMat()
	defer result.Close()

	gocv.MatchTemplate(mainImg, templateImg, &result, gocv.TmCcoeffNormed, gocv.NewMat())

	// Find the maximum value and its location in the result matrix
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)
	log.Info(templateImagePath, "maxVal", maxVal)

	if maxVal >= confidence {
		// Calculate the middle point of the found template
		midX := maxLoc.X + templateImg.Cols()/2
		midY := maxLoc.Y + templateImg.Rows()/2
		return &image.Point{X: midX, Y: midY}, true
	}

	// save screen to debug folder
	gocv.IMWrite(fmt.Sprintf("./debug/%d.png", time.Now().Nanosecond()), mainImg)

	return nil, false
}

func findAllInScreen(templateImagePath string, confidence float32) ([]image.Point, bool) {
	screen, err := captureScreen()
	if err != nil {
		fmt.Println("Error capturing screen:", err)
		return nil, false
	}
	// Read the main image
	mainImg, err := gocv.ImageToMatRGB(screen)
	if err != nil {
		fmt.Println("Error reading main image")
		return nil, false
	}
	defer mainImg.Close()

	// Read the template image
	templateImg := gocv.IMRead(templateImagePath, gocv.IMReadColor)
	if templateImg.Empty() {
		fmt.Println("Error reading template image")
		return nil, false
	}
	defer templateImg.Close()

	// Perform template matching
	result := gocv.NewMat()
	defer result.Close()

	gocv.MatchTemplate(mainImg, templateImg, &result, gocv.TmCcoeffNormed, gocv.NewMat())

	var locations []image.Point

	for {
		_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)

		if maxVal >= confidence {
			// Calculate the middle point of the found template
			midX := maxLoc.X + templateImg.Cols()/2
			midY := maxLoc.Y + templateImg.Rows()/2
			locations = append(locations, image.Point{X: midX, Y: midY})

			for y := maxLoc.Y - 5; y <= maxLoc.Y+templateImg.Rows()+5; y++ {
				for x := maxLoc.X - 5; x <= maxLoc.X+templateImg.Cols()+5; x++ {
					result.SetFloatAt(y, x, 0.0)
				}
			}
		} else {
			break
		}
	}

	if len(locations) > 0 {
		return locations, true
	}

	return nil, false
}

func waitUntilFound(ctx context.Context, templateImagePath string, confidence float32, timeout time.Duration) (*image.Point, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("template %s not found within %s", templateImagePath, timeout)
		default:
			point, found := findInScreen(templateImagePath, confidence)
			if found {
				return point, nil
			}
			log.Info("Template not found, waiting...", "template", templateImagePath)
		}
	}
}

func waitUntilFoundAndClick(ctx context.Context, templateImagePath string, confidence float32, timeout time.Duration) error {
	point, err := waitUntilFound(ctx, templateImagePath, confidence, timeout)
	if err != nil {
		return err
	}

	return clickXY(point.X, point.Y)
}

func clickXY(x, y int) error {
	if x > maxX || y > maxY {
		return fmt.Errorf("x and y must be less than %d and %d respectively", maxX, maxY)
	}
	_, err := runADBCommand("shell", "input", "tap", fmt.Sprintf("%d", x), fmt.Sprintf("%d", y))

	return err
}

func clickImage(imagePath string, confidence float32) error {
	return waitUntilFoundAndClick(context.Background(), fmt.Sprintf("./img/%s.png", imagePath), confidence, 10*time.Second)
}

func openAfkArena() error {
	log.Info("Opening AFK Arena...")
	if _, err := runADBCommand("shell", "monkey", "-p", "com.lilithgame.hgame.gp", "-c", "android.intent.category.LAUNCHER", "1"); err != nil {
		return err
	}

	if _, err := waitUntilFound(context.TODO(), "img/buttons/begin.png", 0.8, 30*time.Second); err != nil {
		return fmt.Errorf("failed to find begin.png: %w", err)
	}

	log.Info("AFK Arena opened successfully!")

	log.Info("Expanding menus...")
	if err := expandMenus(); err != nil {
		return fmt.Errorf("failed to expand menus: %w", err)
	}
	log.Info("Menus expanded successfully!")

	return nil
}

func expandMenus() error {
	arrows, b := findAllInScreen("./img/buttons/downarrow.png", 0.8)
	if !b {
		return nil
	}

	for _, arrow := range arrows {
		if err := clickXY(arrow.X, arrow.Y); err != nil {
			return fmt.Errorf("failed to click down arrow: %w", err)
		}
	}

	return nil
}

func getLowestImagePoint(points []image.Point) image.Point {
	lowest := points[0]
	for _, point := range points {
		if point.Y > lowest.Y {
			lowest = point
		}
	}
	return lowest
}
