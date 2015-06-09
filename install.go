package main

import (
	"fmt"
	"log"
	"os"
)

func installHop(name, hopDir, hopperBin string, force bool) error {
	hopFile := hopDir + "/" + name
	_, err := os.Lstat(hopFile)
	if err == nil {
		hopTarget, err := os.Readlink(hopFile)
		if err == nil && hopTarget == hopperBin {
			log.Printf("%v already installed, nothing to do", name)
			return nil
		} else if force {
			log.Printf("%v file will be replaced by %v hop", hopFile, name)
			if err := os.Remove(hopFile); err != nil {
				return fmt.Errorf("Couldn't replace %v file with hop: %v",
					hopFile, err)
			}
		} else {
			return fmt.Errorf("Couldn't install %v, because the file %v "+
				"already exist and it's not a hop. You could try "+
				"to use the force flag to overwrite it.", name, hopFile)
		}
	}

	if err := os.Symlink(hopperBin, hopFile); err == nil {
		log.Println(name, "successfully installed in", hopDir)
		return nil
	} else {
		return fmt.Errorf("Couldn't install %v: %v", name, err)
	}
}
