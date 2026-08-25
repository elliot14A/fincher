package uploads_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/elliot14A/fincher/internal/api/uploads"
	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Valid 1x1 transparent PNG bytes
var validPNGBytes = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
	0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
}

func setupTestServer(t *testing.T) (*echo.Echo, *ent.Client) {
	client, err := turso.Open(":memory:", "")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}

	ctx := context.Background()
	if err := turso.AutoMigrate(ctx, client); err != nil {
		t.Fatalf("failed to run schema automigrations: %v", err)
	}

	e := echo.New()
	g := e.Group("/api/uploads")
	uploads.RegisterRoutes(g, client)

	return e, client
}

func createMultipartRequest(t *testing.T, fieldName, fileName string, fileBytes []byte) (*http.Request, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := part.Write(fileBytes); err != nil {
		t.Fatalf("failed to write part bytes: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/uploads", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, writer.FormDataContentType()
}

func TestUploadsAPI_Lifecycle(t *testing.T) {
	e, client := setupTestServer(t)
	defer client.Close()

	// 1. Upload valid PNG
	req, _ := createMultipartRequest(t, "file", "avatar.png", validPNGBytes)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp models.UploadResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.ID == "" || resp.MimeType != "image/png" || resp.URL == "" {
		t.Fatalf("unexpected upload response: %+v", resp)
	}

	// 2. Stream uploaded image (GET)
	getReq := httptest.NewRequest(http.MethodGet, "/api/uploads/"+resp.ID, nil)
	getRec := httptest.NewRecorder()
	e.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GET, got %d", getRec.Code)
	}
	if getRec.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("expected image/png content type, got: %s", getRec.Header().Get("Content-Type"))
	}
	if getRec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("expected X-Content-Type-Options: nosniff, got: %s", getRec.Header().Get("X-Content-Type-Options"))
	}
	if !bytes.Equal(getRec.Body.Bytes(), validPNGBytes) {
		t.Fatalf("streamed bytes do not match uploaded bytes")
	}

	// 3. Delete upload
	delReq := httptest.NewRequest(http.MethodDelete, "/api/uploads/"+resp.ID, nil)
	delRec := httptest.NewRecorder()
	e.ServeHTTP(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for DELETE, got %d", delRec.Code)
	}

	// 4. Verify 404 after deletion
	get404Req := httptest.NewRequest(http.MethodGet, "/api/uploads/"+resp.ID, nil)
	get404Rec := httptest.NewRecorder()
	e.ServeHTTP(get404Rec, get404Req)

	if get404Rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 NotFound after deletion, got %d", get404Rec.Code)
	}
}

func TestUploadsAPI_SecurityRejections(t *testing.T) {
	e, client := setupTestServer(t)
	defer client.Close()

	// 1. Reject SVG payload (stored XSS prevention)
	svgPayload := []byte(`<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`)
	reqSVG, _ := createMultipartRequest(t, "file", "logo.svg", svgPayload)
	recSVG := httptest.NewRecorder()
	e.ServeHTTP(recSVG, reqSVG)

	if recSVG.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for SVG upload, got: %d", recSVG.Code)
	}

	// 2. Reject HTML disguised as PNG
	fakePNGPayload := []byte(`<html><body><script>alert(1)</script></body></html>`)
	reqFake, _ := createMultipartRequest(t, "file", "avatar.png", fakePNGPayload)
	recFake := httptest.NewRecorder()
	e.ServeHTTP(recFake, reqFake)

	if recFake.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for HTML disguised as PNG, got: %d", recFake.Code)
	}
}

func TestUploadsAPI_Exact1MBBoundary(t *testing.T) {
	e, client := setupTestServer(t)
	defer client.Close()

	// Construct a valid PNG file padded to exact 1,048,576 bytes
	exact1MB := make([]byte, models.MaxUploadSizeBytes)
	copy(exact1MB, validPNGBytes) // Valid PNG header

	req1MB, _ := createMultipartRequest(t, "file", "exact1mb.png", exact1MB)
	rec1MB := httptest.NewRecorder()
	e.ServeHTTP(rec1MB, req1MB)

	if rec1MB.Code != http.StatusCreated {
		t.Fatalf("expected exact 1MB file to succeed with 201 Created, got %d: %s", rec1MB.Code, rec1MB.Body.String())
	}

	// Construct a file with 1,048,577 bytes (1MB + 1 byte) -> MUST fail with 400
	over1MB := make([]byte, models.MaxUploadSizeBytes+1)
	copy(over1MB, validPNGBytes)

	reqOver, _ := createMultipartRequest(t, "file", "over1mb.png", over1MB)
	recOver := httptest.NewRecorder()
	e.ServeHTTP(recOver, reqOver)

	if recOver.Code != http.StatusBadRequest {
		t.Fatalf("expected 1MB+1 byte file to be rejected with 400 Bad Request, got %d", recOver.Code)
	}
}
