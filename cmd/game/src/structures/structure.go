package structures

import (
	"gorm.io/gorm"
)

const (
	UserTG   Source = "tg"
	UserVK   Source = "vk"
	UserSite Source = "site"

	SourceStateSelectingStory SourceState = 0
	SourceStateReadingStory   SourceState = 1
)

type Source string
type SourceState int
type ProfileSource struct {
	SourceID    int
	Source      Source
	SourceState SourceState
}

type Story struct {
	gorm.Model
	Name          string
	Desc          string
	IsByCommunity bool
	Rating        float32
}

type StoryLine struct {
	gorm.Model
	Text             string
	Step             int
	IsSide           bool
	StoryId          int
	Story            Story
	StoryLineChoices []StoryLineChoice
}

type StoryLineChoice struct {
	gorm.Model
	Text            string
	StoryLineID     int
	StoryLine       StoryLine
	NextStoryLineID *int
	NextStoryLine   *StoryLine
}

type Profile struct {
	gorm.Model
	ProfileSource
	SelectedStoryID *int
	SelectedStory   *Story
}

type ProfileProgress struct {
	gorm.Model
	ProfileID   int
	Profile     Profile
	StoryID     int
	Story       Story
	StoryLineID int
	StoryLine   StoryLine
	ChoiceID    *int
	Choice      *StoryLineChoice
}
