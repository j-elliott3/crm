package domain

import "time"

type Stage string

const (
	StageNewLead 	Stage = "New Lead"
	StageQualified 	Stage = "Qualified"
	StageSurvey 	Stage = "Site Visit / Survery"
	StageQuoteSent 	Stage = "Quote Sent"
	StageWon 		Stage = "Won"
	StageLost 		Stage = "Lost"
)

func AllStages() []Stage {
	return []Stage{
		StageNewLead,
		StageQualified,
		StageSurvey,
		StageQuoteSent,
		StageWon,
		StageLost,
	}
}

type Deal struct {
	ID 				int64
	DealName 		string
	CustomerName 	string
	ContactPerson 	string
	Phone 			string
	Email 			string
	EstimatedValue 	float64
	Stage 			Stage
	Source 			string
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
	NextAction 		string
	NextActionDue 	*time.Time
}