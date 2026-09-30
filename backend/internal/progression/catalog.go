package progression

type SkillCatalog struct {
 Registry *SkillRegistry
}

func NewDefaultSkillCatalog()*SkillCatalog{
 defs:=[]SkillDefinition{
  {ID:"focus",Name:"Focus",Description:"Improves sustained task performance.",Stat:"intelligence",MaxLevel:100},
  {ID:"fitness",Name:"Fitness",Description:"Represents general physical conditioning.",Stat:"endurance",MaxLevel:100},
  {ID:"mobility",Name:"Mobility",Description:"Represents movement efficiency.",Stat:"speed",MaxLevel:100},
  {ID:"discipline",Name:"Discipline",Description:"Represents consistency under pressure.",Stat:"defense",MaxLevel:100},
 }
 return &SkillCatalog{Registry:NewSkillRegistry(defs)}
}

func(c *SkillCatalog) List()[]SkillDefinition{
 if c==nil||c.Registry==nil{return nil}
 out:=make([]SkillDefinition,0,len(c.Registry.definitions))
 for _,d:=range c.Registry.definitions{out=append(out,d)}
 return out
}
