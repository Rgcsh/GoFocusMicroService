package consumers

func SetUp() {
	go QueueConsumerSms()
	go QueueConsumerEmail()
}
