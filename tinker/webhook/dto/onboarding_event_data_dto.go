package dto

type OnboardingEventDataDto struct {
	Message string
}

func NewOnboardingEventDataDto(data map[string]interface{}) *OnboardingEventDataDto {
	message, _ := data["message"].(string)
	return &OnboardingEventDataDto{Message: message}
}

func (dto *OnboardingEventDataDto) ToMap() map[string]interface{} {
	return map[string]interface{}{"message": dto.Message}
}
