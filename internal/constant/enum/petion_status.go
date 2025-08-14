package enum

import "fmt"

const (
	PetitionPublished = "published"
	PetitionApproved  = "approved"
	PetitionRejected  = "rejected"
)

var petitionStatus = []string{
	PetitionPublished,
	PetitionApproved,
	PetitionRejected,
}

var ErrorInvalidPetitionStatus = fmt.Errorf("invalid petition status mus be one of: %s", GetAllPetitionStatus())

func ParsePetitionStatus(status string) (string, error) {
	for _, s := range petitionStatus {
		if s == status {
			return s, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", status, ErrorInvalidPetitionStatus)
}

func GetAllPetitionStatus() []string {
	return petitionStatus
}
