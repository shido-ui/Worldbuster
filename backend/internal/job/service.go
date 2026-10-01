package job

import (
	"errors"
	"sync"
)

var ErrUnknownJob = errors.New("unknown job")
var ErrUnknownPosition = errors.New("unknown position")
var ErrRequirements = errors.New("requirements not met")

type Service struct {
	mu         sync.RWMutex
	jobs       map[string]Job
	employment map[string]Employment
}

func NewService() *Service {
	return &Service{jobs: map[string]Job{}, employment: map[string]Employment{}}
}
func (s *Service) Register(j Job) error {
	if j.ID == "" || len(j.Positions) == 0 {
		return ErrUnknownJob
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[j.ID] = j
	return nil
}
func (s *Service) List() []Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		out = append(out, j)
	}
	return out
}
func (s *Service) Employ(characterID, jobID, positionID string, education, stat int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[jobID]
	if !ok {
		return ErrUnknownJob
	}
	var p *Position
	for i := range j.Positions {
		if j.Positions[i].ID == positionID {
			p = &j.Positions[i]
			break
		}
	}
	if p == nil {
		return ErrUnknownPosition
	}
	if education < p.RequiredEducation || stat < p.RequiredStat {
		return ErrRequirements
	}
	s.employment[characterID] = Employment{CharacterID: characterID, JobID: jobID, PositionID: positionID}
	return nil
}
func (s *Service) GetEmployment(id string) (Employment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.employment[id]
	return e, ok
}
