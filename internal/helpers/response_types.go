package helpers

type Empty struct {
}

type ErrorsArray struct {
	Errors []string `json:"errors" example:"cannot ping database,scheduler offline"`
}

type Error struct {
	Error string `json:"error" example:"a server error was encountered"`
}

type Message struct {
	Message string `json:"message" example:"i just wanted to say hi"`
}
