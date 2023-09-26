package assets

import (
	"github.com/gobuffalo/packr/v2"
)

var assetsBox *packr.Box

func init() {
	assetsBox = packr.New("app:assets", "../../public")
}

func GetAssetsBox() *packr.Box {
	return assetsBox
}
