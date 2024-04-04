package converter

func TopicName2ConsumerName(prefix, topicName string) string {
	return prefix + "-" + topicName
}
