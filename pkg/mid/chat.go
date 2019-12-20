package mid

// Chat is the chat object that represents a chat where the bot has been communicated with
type Chat struct {
	IsShutup bool
}

// Chats is the map of chat ID to Chat object
type Chats map[int64]Chat


