package simulation
import "context"
type PersistenceSink interface { EnsureCharacter(context.Context,SimCharacter) error; RecordMemory(context.Context,string,Memory) error; RecordAction(context.Context,string,Action,bool) error }
type RelationshipPersistence interface { RecordRelationship(context.Context,Relationship) error }
type ReputationPersistence interface { PropagateSocialReputation(context.Context,string,int,int,string) error }
