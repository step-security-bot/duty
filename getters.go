package duty

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/flanksource/duty/context"
	"github.com/flanksource/duty/models"
	"github.com/flanksource/duty/types"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// getterCache caches the results for all the getters in this file.
var getterCache = cache.New(time.Second*90, time.Minute*5)

func cacheKey[T any](field, key string) string {
	var v T
	return fmt.Sprintf("%T:%s=%s", v, field, key)
}

type GetterOption uint8

const (
	GetterOptionNoCache GetterOption = iota + 1
)

type GetterOptions []GetterOption

func (t GetterOptions) IsSet(option GetterOption) bool {
	for _, opt := range t {
		if opt == option {
			return true
		}
	}

	return false
}

func FindCachedAgent(ctx DBContext, id string) (*models.Agent, error) {
	if id == uuid.Nil.String() {
		return nil, nil
	}

	agent, err := findCachedEntity[models.Agent](ctx, id)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func FindCachedCheck(ctx DBContext, id string) (*models.Check, error) {
	check, err := findCachedEntity[models.Check](ctx, id)
	if err != nil {
		return nil, err
	}

	return check, nil
}

func FindCachedCanary(ctx DBContext, id string) (*models.Canary, error) {
	canary, err := findCachedEntity[models.Canary](ctx, id)
	if err != nil {
		return nil, err
	}

	return canary, nil
}

// FindPerson looks up a person by the given identifier which can either be
//   - UUID
//   - email
func FindPerson(ctx DBContext, identifier string, opts ...GetterOption) (*models.Person, error) {
	var field string
	if _, err := uuid.Parse(identifier); err == nil {
		field = "id"
	} else {
		field = "email"
	}

	person, err := findEntityByField[models.Person](ctx, field, identifier, opts...)
	if err != nil {
		return nil, err
	}

	return person, nil
}

// FindTeam looks up a team by the given identifier which can either be
//   - UUID
//   - team name
func FindTeam(ctx DBContext, identifier string, opts ...GetterOption) (*models.Team, error) {
	var field string
	if _, err := uuid.Parse(identifier); err == nil {
		field = "id"
	} else {
		field = "name"
	}

	team, err := findEntityByField[models.Team](ctx, field, identifier, opts...)
	if err != nil {
		return nil, err
	}

	return team, nil
}

type FindOption func(db *gorm.DB)

var LocalFilter = "deleted_at is NULL AND agent_id = '00000000-0000-0000-0000-000000000000' OR agent_id IS NULL"

func PickColumns(columns ...string) FindOption {
	return func(db *gorm.DB) {
		if len(columns) == 0 {
			return
		}
		db.Select(columns)
	}
}

func apply(db *gorm.DB, opts ...FindOption) *gorm.DB {
	for _, opt := range opts {
		opt(db)
	}
	return db
}

func FindChecks(ctx context.Context, resourceSelectors types.ResourceSelectors, opts ...FindOption) (components []models.Check, err error) {
	var uniqueComponents []models.Check
	for _, resourceSelector := range resourceSelectors {
		if resourceSelector.LabelSelector != "" {
			labelComponents, err := FindChecksByLabel(ctx, resourceSelector.LabelSelector, opts...)
			if err != nil {
				return nil, fmt.Errorf("Error getting checks with label selectors[%s]: %v", resourceSelector.LabelSelector, err)
			}
			uniqueComponents = append(uniqueComponents, labelComponents...)
		}
		if resourceSelector.FieldSelector != "" {
			return nil, fmt.Errorf("fieldSelector not supported on checks")
		}
	}

	return lo.UniqBy(uniqueComponents, models.CheckID), nil
}

func FindComponents(ctx context.Context, resourceSelectors types.ResourceSelectors, opts ...FindOption) (components []models.Component, err error) {
	var uniqueComponents []models.Component
	for _, resourceSelector := range resourceSelectors {
		if resourceSelector.LabelSelector != "" {
			labelComponents, err := FindComponentsByLabel(ctx, resourceSelector.LabelSelector, opts...)
			if err != nil {
				return nil, fmt.Errorf("Error getting components with label selectors[%s]: %v", resourceSelector.LabelSelector, err)
			}
			uniqueComponents = append(uniqueComponents, labelComponents...)
		}
		if resourceSelector.FieldSelector != "" {
			fieldComponents, err := FindComponentsByField(ctx, resourceSelector.FieldSelector, opts...)
			if err != nil {
				return nil, fmt.Errorf("Error getting components with field selectors[%s]: %v", resourceSelector.FieldSelector, err)
			}
			uniqueComponents = append(uniqueComponents, fieldComponents...)
		}
	}

	return lo.UniqBy(uniqueComponents, models.ComponentID), nil
}

func getLabelsFromSelector(selector string) (matchLabels map[string]string) {
	matchLabels = make(types.JSONStringMap)
	labels := strings.Split(selector, ",")
	for _, label := range labels {
		if strings.Contains(label, "=") {
			kv := strings.Split(label, "=")
			if len(kv) == 2 {
				matchLabels[kv[0]] = kv[1]
			} else {
				matchLabels[kv[0]] = ""
			}
		}
	}
	return
}

func FindComponentsByLabel(ctx context.Context, labelSelector string, opts ...FindOption) (components []models.Component, err error) {
	if labelSelector == "" {
		return nil, nil
	}
	var items = make(map[string]models.Component)
	matchLabels := getLabelsFromSelector(labelSelector)
	var labels = make(map[string]string)
	var onlyKeys []string
	for k, v := range matchLabels {
		if v != "" {
			labels[k] = v
		} else {
			onlyKeys = append(onlyKeys, k)
		}
	}
	var comps []models.Component
	if err := apply(ctx.DB().Where(LocalFilter).
		Where("labels @> ?", types.JSONStringMap(labels)), opts...).
		Find(&comps).Error; err != nil {
		return nil, err
	}
	for _, c := range comps {
		items[c.ID.String()] = c
	}
	for _, k := range onlyKeys {
		var comps []models.Component
		if err := apply(ctx.DB().Where(LocalFilter).
			Where("labels ?? ?", k), opts...).
			Find(&comps).Error; err != nil {
			return nil, err
		}

		for _, c := range comps {
			items[c.ID.String()] = c
		}
	}
	return lo.Values(items), nil
}

func FindComponentsByField(ctx context.Context, fieldSelector string, opts ...FindOption) ([]models.Component, error) {
	if fieldSelector == "" {
		return nil, nil
	}
	var components = make(map[string]models.Component)
	matchLabels := getLabelsFromSelector(fieldSelector)
	for k, v := range matchLabels {
		var comp []models.Component
		//FIXME FindOptions not applied
		ctx.DB().Raw("select * from lookup_component_by_property(?, ?)", k, v).Scan(&comp)
		for _, c := range comp {
			components[c.ID.String()] = c
		}
	}
	return lo.Values(components), nil
}

func FindChecksByLabel(ctx context.Context, labelSelector string, opts ...FindOption) (components []models.Check, err error) {
	if labelSelector == "" {
		return nil, nil
	}
	var items = make(map[string]models.Check)
	matchLabels := getLabelsFromSelector(labelSelector)
	var labels = make(map[string]string)
	var onlyKeys []string
	for k, v := range matchLabels {
		if v != "" {
			labels[k] = v
		} else {
			onlyKeys = append(onlyKeys, k)
		}
	}
	var comps []models.Check
	if err := apply(ctx.DB().Where(LocalFilter).
		Where("labels @> ?", types.JSONStringMap(labels)), opts...).
		Find(&comps).Error; err != nil {
		return nil, err
	}
	for _, c := range comps {
		items[c.ID.String()] = c
	}
	for _, k := range onlyKeys {
		var comps []models.Check
		if err := apply(ctx.DB().Where(LocalFilter).
			Where("labels ?? ?", k), opts...).
			Find(&comps).Error; err != nil {
			return nil, err
		}

		for _, c := range comps {
			items[c.ID.String()] = c
		}
	}
	return lo.Values(items), nil
}

func FindCachedComponent(ctx DBContext, id string) (*models.Component, error) {
	component, err := findCachedEntity[models.Component](ctx, id)
	if err != nil {
		return nil, err
	}

	return component, nil
}

func FindCachedConfig(ctx DBContext, id string) (*models.ConfigItem, error) {
	config, err := findCachedEntity[models.ConfigItem](ctx, id)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func FindCachedIncident(ctx DBContext, id string) (*models.Incident, error) {
	incident, err := findCachedEntity[models.Incident](ctx, id)
	if err != nil {
		return nil, err
	}

	return incident, nil
}

func findCachedEntity[T any](ctx DBContext, id string) (*T, error) {
	return findEntityByField[T](ctx, "id", id)
}

func findEntityByField[T any](ctx DBContext, field, key string, opts ...GetterOption) (*T, error) {
	if !GetterOptions(opts).IsSet(GetterOptionNoCache) {
		if value, ok := getterCache.Get(cacheKey[T](field, key)); ok {
			if cache, ok := value.(*T); ok {
				return cache, nil
			} else {
				return nil, fmt.Errorf("unexpected cached value type: %T", value)
			}
		}
	}

	var resource T
	if err := ctx.DB().Where(fmt.Sprintf("%s = ?", field), key).First(&resource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	getterCache.SetDefault(cacheKey[T](field, key), &resource)
	return &resource, nil
}
