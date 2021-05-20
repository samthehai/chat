package resolver

import "github.com/samthehai/chat/internal/interfaces/graph/generated"

type Resolver struct {
	query        generated.QueryResolver
	mutation     generated.MutationResolver
	subscription generated.SubscriptionResolver
	message      generated.MessageResolver
}

func NewResolver(
	query *QueryResolver,
	mutation *MutationResolver,
	subscription *SubscriptionResolver,
	message *MessageResolver,
) Resolver {
	return Resolver{
		query:        query,
		mutation:     mutation,
		subscription: subscription,
		message:      message,
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return r.mutation }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return r.query }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return r.subscription }

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return r.message }
