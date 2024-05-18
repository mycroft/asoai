package internal

import "fmt"

func GetCurrentSession() (string, error) {
	uuid, err := DbGetCurrentSession()
	if err != nil {
		return "", fmt.Errorf("could not get current session uuid: %v", err)
	}

	return uuid, nil
}

func SetCurrentSession(uuid string) error {
	if err := DbSetCurrentSession(uuid); err != nil {
		return fmt.Errorf("could not set current session: %v", err)
	}

	return nil
}
