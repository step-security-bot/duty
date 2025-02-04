package context

import (
	gocontext "context"
	"time"

	commons "github.com/flanksource/commons/context"
	dutyGorm "github.com/flanksource/duty/gorm"
	"github.com/flanksource/duty/models"
	"github.com/flanksource/duty/types"
	"github.com/flanksource/kommons"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

type Poolable interface {
	Pool() *pgxpool.Pool
}

type Gormable interface {
	DB() *gorm.DB
}

type Context struct {
	commons.Context
}

func New(opts ...commons.ContextOptions) Context {
	return NewContext(gocontext.Background(), opts...)
}

func NewContext(baseCtx gocontext.Context, opts ...commons.ContextOptions) Context {
	baseOpts := []commons.ContextOptions{
		commons.WithDebugFn(func(ctx commons.Context) bool {
			annotations := getObjectMeta(ctx).Annotations
			return annotations != nil && (annotations["debug"] == "true" || annotations["trace"] == "true")
		}),
		commons.WithTraceFn(func(ctx commons.Context) bool {
			annotations := getObjectMeta(ctx).Annotations
			return annotations != nil && annotations["trace"] == "true"
		}),
	}
	baseOpts = append(baseOpts, opts...)
	return Context{
		Context: commons.NewContext(
			baseCtx,
			baseOpts...,
		),
	}
}

func (k Context) WithTimeout(timeout time.Duration) (Context, gocontext.CancelFunc) {
	ctx, cancelFunc := k.Context.WithTimeout(timeout)
	return Context{
		Context: ctx,
	}, cancelFunc
}

// WithAnyValue is a wrapper around WithValue
func (k Context) WithAnyValue(key, val any) Context {
	return Context{
		Context: k.WithValue(key, val),
	}
}

func (k Context) WithObject(object metav1.ObjectMeta) Context {
	return Context{
		Context: k.WithValue("object", object),
	}
}

func (k Context) WithTopology(topology any) Context {
	return Context{
		Context: k.WithValue("topology", topology),
	}
}

func (k Context) WithUser(user *models.Person) Context {
	k.GetSpan().SetAttributes(attribute.String("user-id", user.ID.String()))
	return Context{
		Context: k.WithValue("user", user),
	}
}

func (k Context) User() *models.Person {
	v := k.Value("user")
	if v == nil {
		return nil
	}
	return v.(*models.Person)
}

// WithAgent sets the current session's agent in the context
func (k Context) WithAgent(agent models.Agent) Context {
	k.GetSpan().SetAttributes(attribute.String("agent-id", agent.ID.String()))
	return Context{
		Context: k.WithValue("agent", agent),
	}
}

func (k Context) Agent() *models.Agent {
	v := k.Value("agent")
	if v == nil {
		return nil
	}
	return lo.ToPtr(v.(models.Agent))
}

func (k Context) WithTrace() Context {
	k.Context = k.Context.WithTrace()
	return k
}

func (k Context) WithDebug() Context {
	k.Context = k.Context.WithDebug()
	return k
}

func (k Context) WithKubernetes(client kubernetes.Interface) Context {
	return Context{
		Context: k.WithValue("kubernetes", client),
	}
}

func (k Context) WithKommons(client *kommons.Client) Context {
	return Context{
		Context: k.WithValue("kommons", client),
	}
}

func (k Context) WithNamespace(namespace string) Context {
	return Context{
		Context: k.WithValue("namespace", namespace),
	}
}

func (k Context) WithDB(db *gorm.DB, pool *pgxpool.Pool) Context {
	return Context{
		Context: k.WithValue("db", db).WithValue("pgxpool", pool),
	}
}

func (k Context) WithDBLogLevel(level string) Context {
	db := k.DB()
	db.Logger = dutyGorm.NewGormLogger(level)
	return Context{
		Context: k.WithValue("db", db),
	}
}

// FastDB returns a db suitable for high-performance usage, with limited logging
func (k Context) FastDB() *gorm.DB {
	db := k.DB().Session(&gorm.Session{
		NewDB: true,
	})
	db.Logger.LogMode(logger.Error)
	return db
}

func (k Context) DB() *gorm.DB {
	val := k.Value("db")
	if val == nil {
		return nil
	}

	v, ok := val.(*gorm.DB)
	if !ok || v == nil {
		return nil
	}
	return v.WithContext(k)
}

func (k Context) Pool() *pgxpool.Pool {
	val := k.Value("pgxpool")
	if val == nil {
		return nil
	}
	v, ok := val.(*pgxpool.Pool)
	if !ok || v == nil {
		return nil
	}
	return v

}

func (k *Context) Kubernetes() kubernetes.Interface {
	v, ok := k.Value("kubernetes").(kubernetes.Interface)
	if !ok || v == nil {
		return fake.NewSimpleClientset()
	}
	return v
}

func (k *Context) Kommons() *kommons.Client {
	v, ok := k.Value("kommons").(*kommons.Client)
	if !ok || v == nil {
		return nil
	}
	return v
}

func (k Context) Topology() any {
	return k.Value("topology")
}

func (k Context) StartSpan(name string) (Context, trace.Span) {
	ctx, span := k.Context.StartSpan(name)
	span.SetAttributes(
		attribute.String("name", k.GetName()),
		attribute.String("namespace", k.GetNamespace()),
	)

	return Context{
		Context: ctx,
	}, span
}

func getObjectMeta(ctx commons.Context) metav1.ObjectMeta {
	o := ctx.Value("object")
	if o == nil {
		return metav1.ObjectMeta{Annotations: map[string]string{}, Labels: map[string]string{}}
	}
	return o.(metav1.ObjectMeta)
}

func (k Context) GetObjectMeta() metav1.ObjectMeta {
	return getObjectMeta(k.Context)
}

func (k Context) GetNamespace() string {
	if k.Value("object") != nil {
		return k.GetObjectMeta().Namespace
	}
	if k.Value("namespace") != nil {
		return k.Value("namespace").(string)
	}
	return ""
}
func (k Context) GetName() string {
	return k.GetObjectMeta().Name
}

func (k Context) GetLabels() map[string]string {
	return k.GetObjectMeta().Labels
}

func (k Context) GetAnnotations() map[string]string {
	return k.GetObjectMeta().Annotations
}

func (k Context) GetEnvValueFromCache(input types.EnvVar, namespace ...string) (string, error) {
	if len(namespace) == 0 {
		namespace = []string{k.GetNamespace()}
	}
	return GetEnvValueFromCache(k, input, namespace[0])
}

func (k Context) GetEnvStringFromCache(env string, namespace string) (string, error) {
	return GetEnvStringFromCache(k, env, namespace)
}

func (k Context) GetSecretFromCache(namespace, name, key string) (string, error) {
	return GetSecretFromCache(k, namespace, name, key)
}

func (k Context) GetConfigMapFromCache(namespace, name, key string) (string, error) {
	return GetConfigMapFromCache(k, namespace, name, key)
}

func (k Context) HydrateConnectionByURL(connectionString string) (*models.Connection, error) {
	return HydrateConnectionByURL(k, connectionString)
}

func (k Context) HydrateConnection(connection *models.Connection) (*models.Connection, error) {
	return HydrateConnection(k, connection)
}

func (k Context) Wrap(ctx gocontext.Context) Context {
	return NewContext(ctx, commons.WithTracer(k.GetTracer())).
		WithDB(k.DB(), k.Pool()).
		WithKubernetes(k.Kubernetes()).
		WithKommons(k.Kommons()).
		WithNamespace(k.GetNamespace())
}
