package main

func next(chatID int64) string {
	users, err := getUsers(chatID)
	if err != nil {
		return ""
	}

	ind, err := getActiveIndex(chatID)
	if err != nil {
		return ""
	}

	ind = (ind + 1) % len(users)

	err = setActiveIndex(chatID, ind)
	if err != nil {
		return ""
	}

	return users[ind]
}

func prev(chatID int64) string {
	users, err := getUsers(chatID)
	if err != nil {
		return ""
	}

	ind, err := getActiveIndex(chatID)
	if err != nil {
		return ""
	}

	if ind == 0 {
		ind = len(users) - 1
	} else {
		ind = ind - 1
	}

	if ind < 0 {
		ind = 0
	}

	err = setActiveIndex(chatID, ind)
	if err != nil {
		return ""
	}

	return users[ind]
}

func who(chatID int64) string {
	users, err := getUsers(chatID)
	if err != nil {
		return ""
	}

	ind, err := getActiveIndex(chatID)
	if err != nil {
		return ""
	}

	return users[ind]
}

func setEstablish(chatID int64, users []string) {
	err := clearUserList(chatID)
	if err != nil {
		panic(err)
	}
	err = setActiveIndex(chatID, 0)
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		addUser(chatID, user)
	}
}
