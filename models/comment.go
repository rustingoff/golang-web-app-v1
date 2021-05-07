package models

func GetComment() ([]string, error) {
	return client.LRange("comments", 0, 10).Result()
}

func PostComment(comment string) error {
	return client.LPush("comments", comment).Err()
}
