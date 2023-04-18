package dummy

import (
	"time"

	"github.com/flanksource/duty/models"
	"github.com/google/uuid"
)

var LogisticsAPIDownIncident = models.Incident{
	ID:          uuid.MustParse("7c05a739-8a1c-4999-85f7-d93d03f32044"),
	Title:       "Logistics API is down",
	CreatedBy:   JohnDoe.ID,
	Type:        models.IncidentTypeAvailability,
	Status:      models.IncidentStatusOpen,
	Severity:    "Blocker",
	CommanderID: &JohnDoe.ID,
}

var AllDummyIncidents = []models.Incident{LogisticsAPIDownIncident}

var FirstComment = models.Comment{
	ID:         uuid.New(),
	CreatedBy:  JohnWick.ID,
	Comment:    "This is a comment",
	IncidentID: LogisticsAPIDownIncident.ID,
	CreatedAt:  time.Now().Add(-time.Hour),
	UpdatedAt:  time.Now().Add(-time.Hour),
}

var SecondComment = models.Comment{
	ID:         uuid.New(),
	CreatedBy:  JohnDoe.ID,
	Comment:    "A comment by John Doe",
	IncidentID: LogisticsAPIDownIncident.ID,
	CreatedAt:  time.Now().Add(-time.Minute),
	UpdatedAt:  time.Now().Add(-time.Minute),
}

var ThirdComment = models.Comment{
	ID:         uuid.New(),
	CreatedBy:  JohnDoe.ID,
	Comment:    "Another comment by John Doe",
	IncidentID: LogisticsAPIDownIncident.ID,
	CreatedAt:  time.Now(),
	UpdatedAt:  time.Now(),
}

var AllDummyComments = []models.Comment{FirstComment, SecondComment, ThirdComment}

var FirstResponder = models.Responder{
	ID:         uuid.New(),
	IncidentID: LogisticsAPIDownIncident.ID,
	Type:       "whattype",
	PersonID:   &JohnWick.ID,
	CreatedBy:  JohnWick.ID,
	CreatedAt:  time.Now().Add(-time.Hour),
	UpdatedAt:  time.Now().Add(-time.Hour),
}

var SecondResponder = models.Responder{
	ID:         uuid.New(),
	IncidentID: LogisticsAPIDownIncident.ID,
	Type:       "whattype",
	PersonID:   &JohnDoe.ID,
	CreatedBy:  JohnDoe.ID,
	CreatedAt:  time.Now().Add(-time.Minute),
	UpdatedAt:  time.Now().Add(-time.Minute),
}

var ThirdResponder = models.Responder{
	ID:         uuid.New(),
	IncidentID: LogisticsAPIDownIncident.ID,
	Type:       "whattype",
	PersonID:   &JohnDoe.ID,
	CreatedBy:  JohnDoe.ID,
	CreatedAt:  time.Now(),
	UpdatedAt:  time.Now(),
}

var AllDummyResponders = []models.Responder{FirstResponder, SecondResponder, ThirdResponder}
