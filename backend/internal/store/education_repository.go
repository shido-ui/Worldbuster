package store

import (
 "context"
 "database/sql"
 "errors"
 "time"
)

var ErrCourseNotFound = errors.New("course not found")
var ErrTrainingActive = errors.New("training already active")
var ErrCourseRequirement = errors.New("course requirements not met")

type CourseRecord struct { ID string `json:"id"`; Name string `json:"name"`; DurationMinutes int `json:"durationMinutes"`; EducationGain int `json:"educationGain"`; RequiredLevel int `json:"requiredLevel"` }

type TrainingRecord struct { CourseID string `json:"courseId"`; CourseName string `json:"courseName"`; StartedAt time.Time `json:"startedAt"`; CompletesAt time.Time `json:"completesAt"`; CompletedAt *time.Time `json:"completedAt,omitempty"` }

type EducationRepository struct{DB *DB}

func(r EducationRepository) ListCourses(ctx context.Context)([]CourseRecord,error){
 rows,err:=r.DB.SQL.QueryContext(ctx,`SELECT id,name,duration_minutes,education_gain,required_level FROM education_courses ORDER BY required_level,name`);if err!=nil{return nil,err};defer rows.Close()
 out:=[]CourseRecord{};for rows.Next(){var c CourseRecord;if err:=rows.Scan(&c.ID,&c.Name,&c.DurationMinutes,&c.EducationGain,&c.RequiredLevel);err!=nil{return nil,err};out=append(out,c)};return out,rows.Err()
}

func(r EducationRepository) Status(ctx context.Context,playerID string)(*TrainingRecord,error){
 var t TrainingRecord;var done sql.NullTime
 err:=r.DB.SQL.QueryRowContext(ctx,`SELECT pt.course_id,ec.name,pt.started_at,pt.completes_at,pt.completed_at FROM player_training pt JOIN education_courses ec ON ec.id=pt.course_id WHERE pt.player_id=$1`,playerID).Scan(&t.CourseID,&t.CourseName,&t.StartedAt,&t.CompletesAt,&done)
 if err==sql.ErrNoRows{return nil,nil};if err!=nil{return nil,err};if done.Valid{t.CompletedAt=&done.Time};return &t,nil
}

func(r EducationRepository) Enroll(ctx context.Context,playerID,courseID string,level int)(TrainingRecord,error){
 var t TrainingRecord
 err:=r.DB.WithTx(ctx,func(tx *sql.Tx)error{
  var active bool
  err:=tx.QueryRowContext(ctx,`SELECT EXISTS(SELECT 1 FROM player_training WHERE player_id=$1 AND completed_at IS NULL)`,playerID).Scan(&active);if err!=nil{return err};if active{return ErrTrainingActive}
  var name string;var mins,gain,required int
  err=tx.QueryRowContext(ctx,`SELECT name,duration_minutes,education_gain,required_level FROM education_courses WHERE id=$1`,courseID).Scan(&name,&mins,&gain,&required)
  if err==sql.ErrNoRows{return ErrCourseNotFound};if err!=nil{return err};if level<required{return ErrCourseRequirement}
  now:=time.Now().UTC();done:=now.Add(time.Duration(mins)*time.Minute)
  _,err=tx.ExecContext(ctx,`INSERT INTO player_training(player_id,course_id,started_at,completes_at) VALUES($1,$2,$3,$4)`,playerID,courseID,now,done)
  t=TrainingRecord{CourseID:courseID,CourseName:name,StartedAt:now,CompletesAt:done};return err
 });return t,err
}

func(r EducationRepository) CompleteDue(ctx context.Context)(int,error){
 n:=0
 err:=r.DB.WithTx(ctx,func(tx *sql.Tx)error{
  rows,err:=tx.QueryContext(ctx,`SELECT pt.player_id::text,pt.course_id,ec.education_gain FROM player_training pt JOIN education_courses ec ON ec.id=pt.course_id WHERE pt.completed_at IS NULL AND pt.completes_at<=NOW() FOR UPDATE OF pt SKIP LOCKED LIMIT 100`);if err!=nil{return err};defer rows.Close()
  type due struct{player string;course string;gain int};var ds []due
  for rows.Next(){var d due;if err:=rows.Scan(&d.player,&d.course,&d.gain);err!=nil{return err};ds=append(ds,d)};if err:=rows.Err();err!=nil{return err}
  for _,d:=range ds{if _,err:=tx.ExecContext(ctx,`UPDATE player_profiles SET education=education+$2,xp=xp+50,updated_at=NOW() WHERE id=$1`,d.player,d.gain);err!=nil{return err};if _,err:=tx.ExecContext(ctx,`UPDATE player_training SET completed_at=NOW() WHERE player_id=$1 AND course_id=$2 AND completed_at IS NULL`,d.player,d.course);err!=nil{return err};n++};return nil
 });return n,err
}