package logger

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// mockWriter 模拟写入器
type mockWriter struct {
	buf *bytes.Buffer
}

func (m *mockWriter) Printf(format string, args ...interface{}) {
	m.buf.WriteString(fmt.Sprintf(format, args...))
}

func newMockWriter() *mockWriter {
	return &mockWriter{buf: &bytes.Buffer{}}
}

func TestLogger_LogMode(t *testing.T) {
	writer := newMockWriter()
	l := New(writer, Config{LogLevel: Info})

	// 测试设置日志级别
	newLogger := l.LogMode(Error)

	// 原日志器级别不变
	l.Info(context.Background(), "info message")
	if !strings.Contains(writer.buf.String(), "info message") {
		t.Error("Expected info message in original logger")
	}

	// 新日志器级别已改变
	writer.buf.Reset()
	newLogger.Info(context.Background(), "info message")
	if strings.Contains(writer.buf.String(), "info message") {
		t.Error("New logger should not log info messages")
	}
}

func TestLogger_Info(t *testing.T) {
	writer := newMockWriter()
	l := New(writer, Config{LogLevel: Info, Colorful: false})

	l.Info(context.Background(), "test info message")

	output := writer.buf.String()
	if !strings.Contains(output, "[info]") {
		t.Error("Expected [info] in output")
	}
	if !strings.Contains(output, "test info message") {
		t.Error("Expected info message in output")
	}
}

func TestLogger_Warn(t *testing.T) {
	writer := newMockWriter()
	l := New(writer, Config{LogLevel: Warn, Colorful: false})

	l.Warn(context.Background(), "test warn message")

	output := writer.buf.String()
	if !strings.Contains(output, "[warn]") {
		t.Error("Expected [warn] in output")
	}
	if !strings.Contains(output, "test warn message") {
		t.Error("Expected warn message in output")
	}
}

func TestLogger_Error(t *testing.T) {
	writer := newMockWriter()
	l := New(writer, Config{LogLevel: Error, Colorful: false})

	l.Error(context.Background(), "test error message")

	output := writer.buf.String()
	if !strings.Contains(output, "[error]") {
		t.Error("Expected [error] in output")
	}
	if !strings.Contains(output, "test error message") {
		t.Error("Expected error message in output")
	}
}

func TestLogger_Trace_Success(t *testing.T) {
	writer := newMockWriter()
	l := New(writer, Config{LogLevel: Info, Colorful: false})

	begin := time.Now()
	time.Sleep(1 * time.Millisecond) // 确保有执行时间

	l.Trace(context.Background(), begin, func() (string, int64) {
		return "SELECT * FROM users WHERE id = ?", 1
	}, nil)

	output := writer.buf.String()
	if !strings.Contains(output, "SELECT * FROM users WHERE id = ?") {
		t.Error("Expected SQL in output")
	}
	if !strings.Contains(output, "[rows:1]") {
		t.Error("Expected rows affected in output")
	}
}

func TestLogger_Trace_Error(t *testing.T) {
	writer := newMockWriter()
	l := New(writer, Config{LogLevel: Error, Colorful: false})

	begin := time.Now()
	testErr := errors.New("database connection failed")

	l.Trace(context.Background(), begin, func() (string, int64) {
		return "SELECT * FROM users WHERE id = ?", -1
	}, testErr)

	output := writer.buf.String()
	if !strings.Contains(output, "database connection failed") {
		t.Error("Expected error message in output")
	}
	if !strings.Contains(output, "SELECT * FROM users WHERE id = ?") {
		t.Error("Expected SQL in output")
	}
}

func TestLogger_Trace_SlowQuery(t *testing.T) {
	writer := newMockWriter()
	l := New(writer, Config{
		LogLevel:      Warn,
		SlowThreshold: 1 * time.Millisecond,
		Colorful:      false,
	})

	begin := time.Now()
	time.Sleep(2 * time.Millisecond) // 超过慢查询阈值

	l.Trace(context.Background(), begin, func() (string, int64) {
		return "SELECT * FROM users", 100
	}, nil)

	output := writer.buf.String()
	if !strings.Contains(output, "SLOW SQL") {
		t.Error("Expected SLOW SQL warning in output")
	}
	if !strings.Contains(output, "SELECT * FROM users") {
		t.Error("Expected SQL in output")
	}
}

func TestLogger_Trace_RecordNotFound_Ignored(t *testing.T) {
	writer := newMockWriter()
	l := New(writer, Config{
		LogLevel:                  Error,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})

	begin := time.Now()

	l.Trace(context.Background(), begin, func() (string, int64) {
		return "SELECT * FROM users WHERE id = ?", 0
	}, RecordNotFoundError)

	output := writer.buf.String()
	if strings.Contains(output, "record not found") {
		t.Error("Should ignore record not found error")
	}
}

func TestLogger_Trace_Silent(t *testing.T) {
	writer := newMockWriter()
	l := New(writer, Config{LogLevel: Silent})

	begin := time.Now()

	l.Trace(context.Background(), begin, func() (string, int64) {
		return "SELECT * FROM users", 1
	}, nil)

	output := writer.buf.String()
	if output != "" {
		t.Error("Silent mode should not output anything")
	}
}

func TestLogger_LogLevels(t *testing.T) {
	tests := []struct {
		name     string
		logLevel LogLevel
		method   func(Interface, context.Context)
		expected bool
	}{
		{
			name:     "Info level allows info",
			logLevel: Info,
			method: func(l Interface, ctx context.Context) {
				l.Info(ctx, "test")
			},
			expected: true,
		},
		{
			name:     "Warn level blocks info",
			logLevel: Warn,
			method: func(l Interface, ctx context.Context) {
				l.Info(ctx, "test")
			},
			expected: false,
		},
		{
			name:     "Error level blocks warn",
			logLevel: Error,
			method: func(l Interface, ctx context.Context) {
				l.Warn(ctx, "test")
			},
			expected: false,
		},
		{
			name:     "Silent level blocks error",
			logLevel: Silent,
			method: func(l Interface, ctx context.Context) {
				l.Error(ctx, "test")
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := newMockWriter()
			l := New(writer, Config{LogLevel: tt.logLevel, Colorful: false})

			tt.method(l, context.Background())

			hasOutput := writer.buf.Len() > 0
			if hasOutput != tt.expected {
				t.Errorf("Expected output: %v, got: %v", tt.expected, hasOutput)
			}
		})
	}
}
