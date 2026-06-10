package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// --- Context Keys ---
type ctxKey string

const (
	keyAuditID   ctxKey = "AuditID"
	keyUserID    ctxKey = "UserID"
	keyEmail     ctxKey = "Email"
	keyRequestID ctxKey = "RequestID"
)

func startTimeKey(method string) ctxKey {
	return ctxKey(method + "StartTime")
}

// --- Logger Initialization ---
var base = func() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	l.SetLevel(logrus.TraceLevel)
	return l
}()

// --- Bootstrap: Create request context with correlation IDs ---
func NewRequestContext(parent context.Context, requestID string) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	ctx := context.WithValue(parent, keyAuditID, uuid.NewString())
	ctx = context.WithValue(ctx, keyRequestID, requestID)
	return ctx
}

// --- Context enrichment ---
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, keyUserID, userID)
}

func WithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, keyEmail, email)
}

// --- Extract fields from context ---
func fieldsFromContext(ctx context.Context, method string) logrus.Fields {
	f := logrus.Fields{}
	if ctx == nil {
		return f
	}

	// String values
	putString := func(k string, key ctxKey) {
		if v, ok := ctx.Value(key).(string); ok && v != "" {
			f[k] = v
		}
	}

	// Integer values
	putInt := func(k string, key ctxKey) {
		if v, ok := ctx.Value(key).(int); ok {
			f[k] = v
		}
	}

	putString("audit_id", keyAuditID)
	putString("request_id", keyRequestID)
	putInt("user_id", keyUserID)
	putString("email", keyEmail)

	if method != "" {
		f["method"] = method
	}

	return f
}

// --- Core logging with context enrichment ---
func logWithContext(ctx context.Context, level logrus.Level, method string, msg interface{}) {
	base.WithFields(fieldsFromContext(ctx, method)).Log(level, msg)
}

// --- Public API: Entry / Exit / Info / Error / Warn ---
func Entry(ctx context.Context, method string) {
	logWithContext(ctx, logrus.TraceLevel, method,
		fmt.Sprintf("Entered %s", method))
}

func Exit(ctx context.Context, method string, elapsedMs int) {
	logWithContext(ctx, logrus.TraceLevel, method,
		fmt.Sprintf("Exited %s (elapsed: %dms)", method, elapsedMs))
}

func Info(ctx context.Context, method string, message interface{}) {
	logWithContext(ctx, logrus.InfoLevel, method, message)
}

func Warn(ctx context.Context, method string, message interface{}) {
	logWithContext(ctx, logrus.WarnLevel, method, message)
}

func Error(ctx context.Context, method string, err error) {
	logWithContext(ctx, logrus.ErrorLevel, method,
		fmt.Sprintf("Error in %s: %v", method, err))
}

// --- Timing helpers ---
func contextWithStartTime(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, startTimeKey(method), time.Now())
}

func ElapsedMs(ctx context.Context, method string) int {
	if ctx == nil {
		return 0
	}
	start, ok := ctx.Value(startTimeKey(method)).(time.Time)
	if !ok {
		return 0
	}
	return int(time.Since(start).Milliseconds())
}

// --- Public API: FuncInitializer / FuncDisposer ---
// Use at start of function: ctx, method := logger.FuncInitializer(ctx, "FunctionName")
// Use at end of function: defer logger.FuncDisposer(ctx, method)
func FuncInitializer(ctx context.Context, method string) (context.Context, string) {
	ctx = contextWithStartTime(ctx, method)
	Entry(ctx, method)
	return ctx, method
}

func FuncDisposer(ctx context.Context, method string) {
	Exit(ctx, method, ElapsedMs(ctx, method))
}
