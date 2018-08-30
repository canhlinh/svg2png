# svg2png
Using chrome headless to convert svg to png

### Example:
```
package main

import (
	"os/exec"

	"github.com/canhlinh/svg2png"
	"github.com/sirupsen/logrus"
)

func main() {
	chrome := svg2png.NewChrome().SetHeight(600).SetWith(600)
	filepath := "Soccerball_mask_transparent_background.png"
	if err := chrome.Screenshoot("https://upload.wikimedia.org/wikipedia/commons/8/84/Example.svg", filepath); err != nil {
		logrus.Panic(err)
	}

	exec.Command("open", filepath).Run()
}

```
