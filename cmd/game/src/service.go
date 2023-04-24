package src

import (
	"wer/cmd/game/db"
	"wer/cmd/game/src/structures"
)

type GameService struct {
	repo    *db.Repository
	profile *structures.Profile
}

func NewGameService() *GameService {
	repo := db.Connect()
	return &GameService{repo: repo}
}

func (g *GameService) SetProfile(p *structures.Profile) {
	g.profile = p
}

func (g *GameService) SelectStory(storyId uint, ps structures.ProfileSource) error {
	s, err := g.repo.Story(storyId)
	if err != nil {
		return err
	}
	err = g.repo.SelectStoryForProfile(s, *g.profile)
	if err != nil {
		return err
	}
	hasPP, err := g.repo.HasProfileProgressForStory(*g.profile, s)
	if err != nil {
		return err
	}
	if !hasPP {
		err = g.repo.CreateProfileProgress(*g.profile, s)
		if err != nil {
			return err
		}
	}
	err = g.repo.UpdateProfileState(*g.profile, structures.SourceStateReadingStory)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameService) StoriesList() ([]structures.Story, error) {
	s, err := g.repo.StoriesList()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (g *GameService) UpdateProfileState(s structures.SourceState) error {
	err := g.repo.UpdateProfileState(*g.profile, s)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameService) CreateProfileIfNotExist(ps structures.ProfileSource) (*structures.Profile, error) {
	p, err := g.repo.ProfileBySource(ps)
	if err != nil {
		return nil, err
	}
	if p.ID == 0 {
		p = structures.Profile{
			ProfileSource: ps,
		}
		p, err = g.repo.CreateNewProfile(p)
		if err != nil {
			return nil, err
		}
	}
	g.SetProfile(&p)
	return &p, nil
}

func (g *GameService) LastStoryLine(p structures.Profile) (*structures.StoryLine, error) {
	sl, err := g.repo.LastStoryLineByProfile(&p)
	if err != nil {
		return nil, err
	}
	return &sl, nil
}

func (g *GameService) NextStoryLineByChoice(slc structures.StoryLineChoice) (*structures.StoryLine, error) {
	sl, err := g.repo.NextStoryLineByChoice(&slc)
	if err != nil {
		return nil, err
	}
	return &sl, nil
}

func (g *GameService) ProfileBySource(ps structures.ProfileSource) (*structures.Profile, error) {
	p, err := g.repo.ProfileBySource(ps)
	if err != nil {
		return nil, err
	}
	g.SetProfile(&p)
	return &p, nil
}

func (g *GameService) LastProfileProgress() (*structures.ProfileProgress, error) {
	pp, err := g.repo.LastProfileProgress(g.profile)
	if err != nil {
		return nil, err
	}
	return &pp, nil
}

func (g *GameService) UpdateLastProfileProgress(pp structures.ProfileProgress) error {
	err := g.repo.UpdateLastProfileProgress(&pp)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameService) AddProfileProgress(pp structures.ProfileProgress) error {
	err := g.repo.AddProfileProgress(pp)
	if err != nil {
		return err
	}
	return nil
}
