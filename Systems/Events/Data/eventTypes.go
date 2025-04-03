package events_data

type EventName string

const (
	ChangeGameState         EventName = "changeState"
	GameSettingsUpdated     EventName = "settingsUpdate"
	ReturnGameState         EventName = "returnState"
	SubmitHighScore         EventName = "submitHighScore"
	AddMessage              EventName = "addMessage"
	ScoreSubmissionComplete EventName = "submissionComplete"
	StatUpgradeEvent        EventName = "statUpgrade"
)
