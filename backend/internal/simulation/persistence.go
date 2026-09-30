package simulation

import "context"

type PersistenceSink interface {
 EnsureCharacter(context.Context,SimCharacter) error
 RecordMemory(context.Context,string,Memory) error
 RecordAction(context.Context,string,Action,int,bool) error
}
