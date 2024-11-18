package main

func saveUserInput(chatID int64, key, value string) {
	if userInputs[chatID] == nil {
		userInputs[chatID] = make(map[string]string)
	}
	userInputs[chatID][key] = value
}
