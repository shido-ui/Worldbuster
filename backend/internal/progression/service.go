package progression

import("errors";"sync")
var ErrUnknownCourse=errors.New("unknown course")
var ErrCourseRequirement=errors.New("course requirements not met")
type Service struct{mu sync.RWMutex;characters map[string]*CharacterProgression;courses map[string]Course;enrollments map[string]Enrollment}
func NewService()*Service{return &Service{characters:map[string]*CharacterProgression{},courses:map[string]Course{},enrollments:map[string]Enrollment{}}}
func(s *Service)Get(id string)CharacterProgression{s.mu.Lock();defer s.mu.Unlock();if c,ok:=s.characters[id];ok{return *c};c:=&CharacterProgression{CharacterID:id,Level:1};s.characters[id]=c;return *c}
func(s *Service)AddXP(id string,xp int64)CharacterProgression{s.mu.Lock();defer s.mu.Unlock();c:=s.getLocked(id);c.XP+=xp;for c.XP>=int64(c.Level)*1000{c.XP-=int64(c.Level)*1000;c.Level++};return *c}
func(s *Service)RegisterCourse(c Course)error{if c.ID==""||c.DurationHours<=0||c.EducationGain<=0{return ErrUnknownCourse};s.mu.Lock();defer s.mu.Unlock();s.courses[c.ID]=c;return nil}
func(s *Service)Enroll(charID,courseID string)error{s.mu.Lock();defer s.mu.Unlock();course,ok:=s.courses[courseID];if !ok{return ErrUnknownCourse};c:=s.getLocked(charID);if c.Level<course.RequiredLevel{return ErrCourseRequirement};s.enrollments[charID]=Enrollment{CharacterID:charID,CourseID:courseID};return nil}
func(s *Service)ListCourses()[]Course{s.mu.RLock();defer s.mu.RUnlock();out:=make([]Course,0,len(s.courses));for _,c:=range s.courses{out=append(out,c)};return out}
func(s *Service)getLocked(id string)*CharacterProgression{if c,ok:=s.characters[id];ok{return c};c:=&CharacterProgression{CharacterID:id,Level:1};s.characters[id]=c;return c}
