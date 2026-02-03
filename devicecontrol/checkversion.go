package devicecontrol

import (
	"strconv"
	"strings"
)

// Предполагаем, что версии имеют формат "x.x.x"
func (dev *Smartphone) canUpdateCurrentVersion(maxVersion string) error {
	maxVerParts := strings.Split(maxVersion, ".")
	currentVerParts := strings.Split(dev.info.OSVersion, ".")

	if len(maxVerParts) > len(currentVerParts) {
		for _ = range len(maxVerParts) - len(currentVerParts) {
			currentVerParts = append(currentVerParts, "0")
		}
	}

	canUpdate := false
	for i := 0; i < len(currentVerParts); i++ {
		if i >= len(maxVerParts) {
			return ErrUnsupported
		}

		curV, err := strconv.Atoi(currentVerParts[i])
		if err != nil {
			return ErrUnsupported
		}

		curMaxV, err := strconv.Atoi(maxVerParts[i])
		if err != nil {
			return ErrUnsupported
		}

		if curV > curMaxV {
			break
		}

		if curV < curMaxV {
			canUpdate = true
			break
		}
	}

	if !canUpdate {
		return ErrUnsupported
	}

	return nil
}
