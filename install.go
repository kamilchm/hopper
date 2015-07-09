// Hop installation package
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Installs hop in given workspace
func (w *workspace) installHop(name string, force bool) error {
	hopFile := filepath.Join(w.BinDir, name)
	_, err := os.Lstat(hopFile)
	if err == nil {
		hopTarget, err := os.Readlink(hopFile)
		if err == nil && hopTarget == w.HopperPath {
			log.Notice("%v already installed, nothing to do", name)
			return nil
		} else if force {
			log.Info("%v file will be replaced by %v hop", hopFile, name)
			if err := os.Remove(hopFile); err != nil {
				return fmt.Errorf("Couldn't replace %v file with hop: %v",
					hopFile, err)
			}
		} else {
			return fmt.Errorf("Couldn't install %v, because the file %v "+
				"already exist and it's not a hop. You could try "+
				"to use the --force flag to overwrite it.", name, hopFile)
		}
	}

	if err := os.Symlink(w.HopperPath, hopFile); err == nil {
		log.Notice("%v successfully installed in %v", name, w.BinDir)
		return nil
	} else {
		return fmt.Errorf("Couldn't install %v: %v", name, err)
	}
}
