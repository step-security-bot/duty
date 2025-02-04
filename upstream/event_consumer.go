package upstream

import (
	"fmt"

	"github.com/flanksource/duty/context"
	"github.com/flanksource/duty/models"
	"github.com/flanksource/postq"
	"github.com/google/uuid"
)

const (
	EventPushQueueCreate = "push_queue.create"

	// EventPushQueueDelete is fired when a record, on one of the the tables we're tracking,
	// is hard deleted.
	EventPushQueueDelete = "push_queue.delete"
)

// NewDeleteFromUpstreamConsumer acts as an adapter to supply DeleteFromUpstream event consumer.
func NewDeleteFromUpstreamConsumer(config UpstreamConfig) func(ctx context.Context, events postq.Events) postq.Events {
	return func(ctx context.Context, events postq.Events) postq.Events {
		return DeleteFromUpstream(ctx, config, events)
	}
}

// DeleteFromUpstream sends a delete request to the upstream server for the given events.
func DeleteFromUpstream(ctx context.Context, config UpstreamConfig, events []postq.Event) []postq.Event {
	upstreamMsg := &PushData{
		AgentName: config.AgentName,
	}

	var failedEvents []postq.Event
	for _, cl := range GroupChangelogsByTables(events) {
		switch cl.TableName {
		case "topologies":
			for i := range cl.ItemIDs {
				upstreamMsg.Topologies = append(upstreamMsg.Topologies, models.Topology{ID: uuid.MustParse(cl.ItemIDs[i][0])})
			}

		case "components":
			for i := range cl.ItemIDs {
				upstreamMsg.Components = append(upstreamMsg.Components, models.Component{ID: uuid.MustParse(cl.ItemIDs[i][0])})
			}

		case "canaries":
			for i := range cl.ItemIDs {
				upstreamMsg.Canaries = append(upstreamMsg.Canaries, models.Canary{ID: uuid.MustParse(cl.ItemIDs[i][0])})
			}

		case "checks":
			for i := range cl.ItemIDs {
				upstreamMsg.Checks = append(upstreamMsg.Checks, models.Check{ID: uuid.MustParse(cl.ItemIDs[i][0])})
			}

		case "config_scrapers":
			for i := range cl.ItemIDs {
				upstreamMsg.ConfigScrapers = append(upstreamMsg.ConfigScrapers, models.ConfigScraper{ID: uuid.MustParse(cl.ItemIDs[i][0])})
			}

		case "config_items":
			for i := range cl.ItemIDs {
				upstreamMsg.ConfigItems = append(upstreamMsg.ConfigItems, models.ConfigItem{ID: uuid.MustParse(cl.ItemIDs[i][0])})
			}

		case "config_component_relationships":
			for i := range cl.ItemIDs {
				upstreamMsg.ConfigComponentRelationships = append(upstreamMsg.ConfigComponentRelationships, models.ConfigComponentRelationship{
					ComponentID: uuid.MustParse(cl.ItemIDs[i][0]),
					ConfigID:    uuid.MustParse(cl.ItemIDs[i][1]),
				})
			}

		case "component_relationships":
			for i := range cl.ItemIDs {
				upstreamMsg.ComponentRelationships = append(upstreamMsg.ComponentRelationships, models.ComponentRelationship{
					ComponentID:    uuid.MustParse(cl.ItemIDs[i][0]),
					RelationshipID: uuid.MustParse(cl.ItemIDs[i][1]),
					SelectorID:     cl.ItemIDs[i][2],
				})
			}

		case "config_relationships":
			for i := range cl.ItemIDs {
				upstreamMsg.ConfigRelationships = append(upstreamMsg.ConfigRelationships, models.ConfigRelationship{
					RelatedID:  cl.ItemIDs[i][0],
					ConfigID:   cl.ItemIDs[i][1],
					SelectorID: cl.ItemIDs[i][2],
				})
			}
		}
	}

	upstreamClient := NewUpstreamClient(config)
	err := upstreamClient.Delete(ctx, upstreamMsg)
	if err != nil {
		if len(events) == 1 {
			errMsg := fmt.Errorf("failed to push delete items to upstream: %w", err)
			failedEvents = append(failedEvents, addErrorToFailedEvents(events, errMsg)...)
		} else {
			// Error encountered while pushing could be an SQL or Application error
			// Since we do not know which event in the bulk is failing
			// Process each event individually since upsteam.Push is idempotent
			for _, e := range events {
				failedEvents = append(failedEvents, DeleteFromUpstream(ctx, config, []postq.Event{e})...)
			}
		}
	}

	return failedEvents
}

// getPushUpstreamConsumer acts as an adapter to supply PushToUpstream event consumer.
func NewPushUpstreamConsumer(config UpstreamConfig) func(ctx context.Context, events postq.Events) postq.Events {
	return func(ctx context.Context, events postq.Events) postq.Events {
		return PushToUpstream(ctx, config, events)
	}
}

// PushToUpstream fetches records specified in events from this instance and sends them to the upstream instance.
func PushToUpstream(ctx context.Context, config UpstreamConfig, events []postq.Event) []postq.Event {
	upstreamMsg := &PushData{
		AgentName: config.AgentName,
	}

	var failedEvents []postq.Event
	for _, cl := range GroupChangelogsByTables(events) {
		switch cl.TableName {
		case "topologies":
			if err := ctx.DB().Where("id IN ?", cl.ItemIDs).Find(&upstreamMsg.Topologies).Error; err != nil {
				errMsg := fmt.Errorf("error fetching topologies: %w", err)
				failedEvents = append(failedEvents, addErrorToFailedEvents(cl.Events, errMsg)...)
			}

		case "components":
			if err := ctx.DB().Where("id IN ?", cl.ItemIDs).Find(&upstreamMsg.Components).Error; err != nil {
				errMsg := fmt.Errorf("error fetching components: %w", err)
				failedEvents = append(failedEvents, addErrorToFailedEvents(cl.Events, errMsg)...)
			}

		case "canaries":
			if err := ctx.DB().Where("id IN ?", cl.ItemIDs).Find(&upstreamMsg.Canaries).Error; err != nil {
				errMsg := fmt.Errorf("error fetching canaries: %w", err)
				failedEvents = append(failedEvents, addErrorToFailedEvents(cl.Events, errMsg)...)
			}

		case "checks":
			if err := ctx.DB().Where("id IN ?", cl.ItemIDs).Find(&upstreamMsg.Checks).Error; err != nil {
				errMsg := fmt.Errorf("error fetching checks: %w", err)
				failedEvents = append(failedEvents, addErrorToFailedEvents(cl.Events, errMsg)...)
			}

		case "config_scrapers":
			if err := ctx.DB().Where("id IN ?", cl.ItemIDs).Find(&upstreamMsg.ConfigScrapers).Error; err != nil {
				errMsg := fmt.Errorf("error fetching config_scrapers: %w", err)
				failedEvents = append(failedEvents, addErrorToFailedEvents(cl.Events, errMsg)...)
			}

		case "config_items":
			if err := ctx.DB().Where("id IN ?", cl.ItemIDs).Find(&upstreamMsg.ConfigItems).Error; err != nil {
				errMsg := fmt.Errorf("error fetching config_items: %w", err)
				failedEvents = append(failedEvents, addErrorToFailedEvents(cl.Events, errMsg)...)
			}

		case "config_component_relationships":
			if err := ctx.DB().Where("(component_id, config_id) IN ?", cl.ItemIDs).Find(&upstreamMsg.ConfigComponentRelationships).Error; err != nil {
				errMsg := fmt.Errorf("error fetching config_component_relationships: %w", err)
				failedEvents = append(failedEvents, addErrorToFailedEvents(cl.Events, errMsg)...)
			}

		case "component_relationships":
			if err := ctx.DB().Where("(component_id, relationship_id, selector_id) IN ?", cl.ItemIDs).Find(&upstreamMsg.ComponentRelationships).Error; err != nil {
				errMsg := fmt.Errorf("error fetching component_relationships: %w", err)
				failedEvents = append(failedEvents, addErrorToFailedEvents(cl.Events, errMsg)...)
			}

		case "config_relationships":
			if err := ctx.DB().Where("(related_id, config_id, selector_id) IN ?", cl.ItemIDs).Find(&upstreamMsg.ConfigRelationships).Error; err != nil {
				errMsg := fmt.Errorf("error fetching config_relationships: %w", err)
				failedEvents = append(failedEvents, addErrorToFailedEvents(cl.Events, errMsg)...)
			}
		}
	}

	upstreamMsg.ApplyLabels(config.LabelsMap())

	upstreamClient := NewUpstreamClient(config)
	err := upstreamClient.Push(ctx, upstreamMsg)
	if err == nil {
		return failedEvents
	}

	if len(events) == 1 {
		errMsg := fmt.Errorf("failed to push to upstream: %w", err)
		failedEvents = append(failedEvents, addErrorToFailedEvents(events, errMsg)...)
	} else {
		// Error encountered while pushing could be an SQL or Application error
		// Since we do not know which event in the bulk is failing
		// Process each event individually since upsteam.Push is idempotent

		for _, e := range events {
			failedEvents = append(failedEvents, PushToUpstream(ctx, config, []postq.Event{e})...)
		}
	}

	if len(events) > 0 || len(failedEvents) > 0 {
		ctx.Tracef("processed %d events, %d errors", len(events), len(failedEvents))
	}

	return failedEvents
}

func addErrorToFailedEvents(events []postq.Event, err error) []postq.Event {
	var failedEvents []postq.Event
	for _, e := range events {
		e.SetError(err.Error())
		failedEvents = append(failedEvents, e)
	}

	return failedEvents
}

// GroupChangelogsByTables groups the given events by the table they belong to.
func GroupChangelogsByTables(events []postq.Event) []GroupedPushEvents {
	type pushEvent struct {
		TableName string
		ItemIDs   []string
		Event     postq.Event
	}

	var pushEvents []pushEvent
	for _, cl := range events {
		tableName := cl.Properties["table"]
		var itemIDs []string
		switch tableName {
		case "component_relationships":
			itemIDs = []string{cl.Properties["component_id"], cl.Properties["relationship_id"], cl.Properties["selector_id"]}
		case "config_component_relationships":
			itemIDs = []string{cl.Properties["component_id"], cl.Properties["config_id"]}
		case "config_relationships":
			itemIDs = []string{cl.Properties["related_id"], cl.Properties["config_id"], cl.Properties["selector_id"]}
		default:
			itemIDs = []string{cl.Properties["id"]}
		}
		pe := pushEvent{
			TableName: tableName,
			ItemIDs:   itemIDs,
			Event:     cl,
		}
		pushEvents = append(pushEvents, pe)
	}

	tblPushMap := make(map[string]*GroupedPushEvents)
	var group []GroupedPushEvents
	for _, p := range pushEvents {
		if k, exists := tblPushMap[p.TableName]; exists {
			k.ItemIDs = append(k.ItemIDs, p.ItemIDs)
			k.Events = append(k.Events, p.Event)
		} else {
			gp := &GroupedPushEvents{
				TableName: p.TableName,
				ItemIDs:   [][]string{p.ItemIDs},
				Events:    []postq.Event{p.Event},
			}
			tblPushMap[p.TableName] = gp
		}
	}

	for _, v := range tblPushMap {
		group = append(group, *v)
	}
	return group
}

type GroupedPushEvents struct {
	TableName string
	ItemIDs   [][]string
	Events    postq.Events
}
