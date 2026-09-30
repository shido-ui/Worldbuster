package api

import (
 "encoding/json"
 "errors"
 "net/http"
 "github.com/shido-ui/Worldbuster/backend/internal/store"
)

type EducationAPI struct{Repo store.EducationRepository; Context *PlayerContext}

func(a *EducationAPI) Courses(w http.ResponseWriter,r *http.Request){courses,err:=a.Repo.ListCourses(r.Context());if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"courses unavailable"});return};writeJSON(w,http.StatusOK,map[string]any{"courses":courses})}

func(a *EducationAPI) Status(w http.ResponseWriter,r *http.Request){if a.Context==nil{writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return};_,p,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return};t,err:=a.Repo.Status(r.Context(),p.ID);if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"training unavailable"});return};writeJSON(w,http.StatusOK,map[string]any{"education":p.Education,"training":t})}

func(a *EducationAPI) Enroll(w http.ResponseWriter,r *http.Request){if a.Context==nil{writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return};_,p,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return};var req map[string]string;if json.NewDecoder(r.Body).Decode(&req)!=nil||req["courseId"]==""{writeJSON(w,http.StatusBadRequest,map[string]string{"error":"courseId required"});return};t,err:=a.Repo.Enroll(r.Context(),p.ID,req["courseId"],p.Level);if errors.Is(err,store.ErrTrainingActive){writeJSON(w,http.StatusConflict,map[string]string{"error":"training already active"});return};if errors.Is(err,store.ErrCourseNotFound){writeJSON(w,http.StatusNotFound,map[string]string{"error":"course not found"});return};if errors.Is(err,store.ErrCourseRequirement){writeJSON(w,http.StatusForbidden,map[string]string{"error":"course requirements not met"});return};if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"enrollment failed"});return};writeJSON(w,http.StatusOK,map[string]any{"enrolled":true,"training":t})}

