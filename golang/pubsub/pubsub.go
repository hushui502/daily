package main

// Pubsub is prove of concept implement for Redis "Pub/Sub" messaging management feature.
// SUBSCRIBE, UNSUBSCRIBE and PUBLISH implement the Publish/Subscribe messaging paradigm
// where (citing Wikipedia) senders (publishers) are not programmed to send their messages
// to specific receivers (subscribers). (sited from here)

type chanMapStringList map[chan interface{}][]string
type stringMapChanList map[string][]chan interface{}

type Pubsub struct {
	capacity int
	clientMapTopics chanMapStringList
	topicMapClients stringMapChanList
}

func NewPubSub(initChanCapacity int) *Pubsub {
	initClientMapTopics := make(chanMapStringList)
	initTopicMapClients := make(stringMapChanList)

	server := Pubsub{
		capacity:        initChanCapacity,
		clientMapTopics: initClientMapTopics,
		topicMapClients: initTopicMapClients,
	}
	server.capacity = initChanCapacity

	return &server
}

func (p *Pubsub) Publish(content interface{}, topics ...string) {
	for _, topic := range topics {
		if chanList, ok := p.topicMapClients[topic]; ok {
			for _, channel := range chanList {
				channel <- content
			}
		}
	}
}

func (p *Pubsub) Subscribe(topics ...string) chan interface{} {
	workChan := make(chan interface{}, p.capacity)
	p.updateTopicMapClient(workChan, topics)

	return workChan
}

func (p *Pubsub) updateTopicMapClient(clientChan chan interface{}, topics []string) {
	var updateChanList []chan interface{}

	for _, topic := range topics {
		updateChanList = p.topicMapClients[topic]
		updateChanList = append(updateChanList, clientChan)
		p.topicMapClients[topic] = updateChanList
	}

	p.clientMapTopics[clientChan] = topics
}

func (p *Pubsub) AddSubscription(clientChan chan interface{}, topics ...string) {
	p.updateTopicMapClient(clientChan, topics)
}

func (p *Pubsub) RemoveSubscription(clientChan chan interface{}, topics ...string) {
	for _, topic := range topics {
		if chanList, ok := p.topicMapClients[topic]; ok {
			var updateChanList []chan interface{}
			for _, client := range chanList {
				if client != clientChan {
					updateChanList = append(updateChanList, client)
				}
			}
			p.topicMapClients[topic] = updateChanList
		}

		if topicList, ok := p.clientMapTopics[clientChan]; ok {
			var updateTopicList []string
			for _, updateTopic := range topicList {
				if updateTopic != topic {
					updateTopicList = append(updateTopicList, updateTopic)
				}
			}
			p.clientMapTopics[clientChan] = updateTopicList
		}
	}
}



















