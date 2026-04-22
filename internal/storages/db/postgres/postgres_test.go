package postgres

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v5"
)

func TestClient_Close(t *testing.T) {
	t.Parallel()

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	client := &Client{Conn: mock}

	mock.ExpectClose()
	if err = client.Close(context.Background()); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}