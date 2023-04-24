package db

import (
	structures "wer/cmd/game/src/structures"
)

func (r Repository) Story(id uint) (structures.Story, error) {
	var story structures.Story

	r.db.Table("stories").First(&story, id)

	return story, nil
}

func (r Repository) StoriesList() ([]structures.Story, error) {
	var stories []structures.Story

	r.db.Table("stories").Limit(3).Find(&stories)

	return stories, nil
}

func (r Repository) CreateNewProfile(p structures.Profile) (structures.Profile, error) {

	r.db.Create(&p)

	return p, nil
}

func (r Repository) ProfileBySource(ps structures.ProfileSource) (structures.Profile, error) {
	var profile *structures.Profile

	r.db.Table("profiles").Where("source_id=? AND source=?", ps.SourceID, ps.Source).First(&profile)

	return *profile, nil
}

func (r Repository) SelectStoryForProfile(s structures.Story, p structures.Profile) error {
	p.SelectedStory = &s
	r.db.Save(&p)

	return nil
}

func (r Repository) HasProfileProgressForStory(p structures.Profile, s structures.Story) (bool, error) {
	var pp *structures.ProfileProgress

	r.db.Table("profile_progresses").Where("profile_id=? AND story_id=?", p.ID, s.ID).First(&pp)

	if pp.ID == 0 {
		return false, nil
	}

	return true, nil
}

func (r Repository) CreateProfileProgress(p structures.Profile, s structures.Story) error {
	var sl *structures.StoryLine
	r.db.Table("story_lines").Where("story_id=? AND step=1", s.ID).First(&sl)

	pp := structures.ProfileProgress{
		Profile:   p,
		Story:     s,
		StoryLine: *sl,
	}
	r.db.Create(&pp)

	return nil
}

func (r Repository) AddProfileProgress(pp structures.ProfileProgress) error {

	r.db.Create(&pp)

	return nil
}

func (r Repository) LastStoryLineByProfile(p *structures.Profile) (structures.StoryLine, error) {
	var profileProgress *structures.ProfileProgress

	r.db.Table("profile_progresses").
		Where("story_id=? AND profile_id=?", p.SelectedStoryID, p.ID).
		Preload("StoryLine").
		Preload("StoryLine.StoryLineChoices").
		Last(&profileProgress)

	return profileProgress.StoryLine, nil
}

func (r Repository) NextStoryLineByChoice(slc *structures.StoryLineChoice) (structures.StoryLine, error) {
	var sl *structures.StoryLine

	r.db.Table("story_lines").Where("id=? AND story_id=?", slc.NextStoryLineID, slc.StoryLineID).Preload("Story").Preload("StoryLineChoices").First(&sl)

	return *sl, nil
}

func (r Repository) LastProfileProgress(p *structures.Profile) (structures.ProfileProgress, error) {
	var profileProgress *structures.ProfileProgress

	r.db.Table("profile_progresses").Where("profile_id=? AND story_id=?", p.ID, p.SelectedStoryID).Preload("StoryLine").Preload("StoryLine.StoryLineChoices").First(&profileProgress)

	return *profileProgress, nil
}

func (r Repository) UpdateLastProfileProgress(pp *structures.ProfileProgress) error {

	r.db.Model(&pp).Update("choice_id", pp.ChoiceID)

	return nil
}

func (r Repository) UpdateProfileState(p structures.Profile, s structures.SourceState) error {

	r.db.Model(&p).Update("source_state", s)

	return nil
}
