INSERT INTO jobs (name,department,base_salary,required_level,required_stat)
VALUES ('General Workforce','Operations',100,1,0)
ON CONFLICT (name) DO NOTHING;