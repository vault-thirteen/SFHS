package server

import (
	"fmt"
	"runtime"
)

const (
	LibName    = "SFHS"
	LibVersion = "0.11.0"

	ProductServer = "Server"

	ComponentSFRODB = "SFRODB"
)

const (
	StartupText       = "%s %s, ver. %s, %s."
	ComponentInfoText = "%s Version: %s."
)

func ShowIntroText(product string) {
	fmt.Println(
		fmt.Sprintf(StartupText, LibName, product, LibVersion, runtime.Version()),
	)
}

func ShowComponentInfoText(componentName string, componentVersion string) {
	fmt.Println(fmt.Sprintf(ComponentInfoText, componentName, componentVersion))
}
