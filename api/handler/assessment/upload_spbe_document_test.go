package assessment

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"image"
// 	"image/color"
// 	"image/png"
// 	"io"
// 	"log"
// 	"mime/multipart"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	userCtx "ta-spbe-backend/api/handler/context"
// 	"ta-spbe-backend/config"
// 	"ta-spbe-backend/repository"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/stretchr/testify/assert"
// 	waProto "go.mau.fi/whatsmeow/binary/proto"
// )

// type userRepoMock struct{}
// type messageQueueMock struct{}
// type mailerMock struct{}
// type fileSystemIOMock struct{}
// type jsonECMock struct{}
// type waClientMock struct{}

// var (
// 	findOneByEmail       func(ctx context.Context, email string) (*repository.User, error)
// 	findOneByID          func(ctx context.Context, id string) (*repository.User, error)
// 	findEmailAndPassword func(ctx context.Context, email string) (*repository.User, error)
// 	findAdmin            func(ctx context.Context, email string) (*repository.User, error)
// 	insert               func(ctx context.Context, user *repository.User) error
// 	produce              func(topic string, body []byte) error
// 	send                 func(subject, message []byte, receiver []string, templateName string, items interface{}) error
// 	create               func(name string) (*os.File, error)
// 	copy                 func(dst io.Writer, src io.Reader) (int64, error)
// 	marshal              func(v any) ([]byte, error)
// 	sendMessage          func(ctx context.Context, recipientNumber string, message *waProto.Message) error
// )

// func (repo userRepoMock) FindOneByEmail(ctx context.Context, email string) (*repository.User, error) {
// 	return findOneByEmail(ctx, email)
// }

// func (repo userRepoMock) FindOneByID(ctx context.Context, id string) (*repository.User, error) {
// 	return findOneByID(ctx, id)
// }

// func (repo userRepoMock) FindEmailAndPassword(ctx context.Context, email string) (*repository.User, error) {
// 	return findEmailAndPassword(ctx, email)
// }

// func (repo userRepoMock) FindAdmin(ctx context.Context, email string) (*repository.User, error) {
// 	return findAdmin(ctx, email)
// }

// func (repo userRepoMock) Insert(ctx context.Context, user *repository.User) error {
// 	return insert(ctx, user)
// }

// func (mq messageQueueMock) Produce(topic string, body []byte) error {
// 	return produce(topic, body)
// }

// func (mailer mailerMock) Send(subject, message []byte, receiver []string, templateName string, items interface{}) error {
// 	return send(subject, message, receiver, templateName, items)
// }

// func (fsIO fileSystemIOMock) Create(name string) (*os.File, error) {
// 	return create(name)
// }

// func (fsIO fileSystemIOMock) Copy(dst io.Writer, src io.Reader) (int64, error) {
// 	return copy(dst, src)
// }

// func (jsonEC jsonECMock) Marshal(v any) ([]byte, error) {
// 	return marshal(v)
// }

// func (waClient waClientMock) SendMessage(ctx context.Context, recipientNumber string, message *waProto.Message) error {
// 	return sendMessage(ctx, recipientNumber, message)
// }

// func TestUploadSPBEDocument_Success(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	insertUploadDocumentAssessmentRepoMock = func(ctx context.Context, assessmentUploadDetail *repository.AssessmentUploadDetail) error {
// 		return nil
// 	}
// 	findOneByID = func(ctx context.Context, id string) (*repository.User, error) {
// 		return &repository.User{
// 			ID:    "test-id",
// 			Name:  "test-name",
// 			Email: "test-email",
// 			Role:  "ADMIN",
// 		}, nil
// 	}
// 	create = func(name string) (*os.File, error) {
// 		return &os.File{}, nil
// 	}
// 	copy = func(dst io.Writer, src io.Reader) (int64, error) {
// 		return 0, nil
// 	}
// 	produce = func(topic string, body []byte) error {
// 		return nil
// 	}
// 	send = func(subject, message []byte, receiver []string, templateName string, items interface{}) error {
// 		return nil
// 	}
// 	marshal = func(v any) ([]byte, error) {
// 		return nil, nil
// 	}
// 	sendMessage = func(ctx context.Context, recipientNumber string, message *waProto.Message) error {
// 		return nil
// 	}

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockMailer := mailerMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	ctx := context.WithValue(req.Context(), userCtx.UserCtxKey, userCtx.UserCtx{
// 		ID:   "test-user-id",
// 		Role: repository.UserRole("ADMIN"),
// 	})
// 	req = req.WithContext(ctx)
// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusCreated, rr.Code)

// 	os.Remove("./image.png")
// }

// func TestUploadSPBEDocument_FailMissingSupportingDocument(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()
// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockMailer := mailerMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusUnprocessableEntity, rr.Code)
// }

// func TestUploadSPBEDocument_FailSupportingDocumentFileSize(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		f, err := os.Create("foo.bar")
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if err := f.Truncate(1e9); err != nil {
// 			log.Fatal(err)
// 		}
// 		part, err := writer.CreateFormFile("supporting_document", "foo.bar")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		if _, err := io.Copy(part, f); err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockMailer := mailerMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusUnprocessableEntity, rr.Code)

// 	os.Remove("./foo.bar")
// }

// func TestUploadSPBEDocument_FailMissingInstitutionName(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockMailer := mailerMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusUnprocessableEntity, rr.Code)

// 	os.Remove("./image.png")
// }

// func TestUploadSPBEDocument_FailMissingIndicatorNumber(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockMailer := mailerMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusUnprocessableEntity, rr.Code)

// 	os.Remove("./image.png")
// }

// func TestUploadSPBEDocument_FailInvalidIndicatorNumber(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("invalid_indicator_number"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockMailer := mailerMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusUnprocessableEntity, rr.Code)

// 	os.Remove("./image.png")
// }

// func TestUploadSPBEDocument_FailInvalidUserCredential(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockMailer := mailerMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)

// 	os.Remove("./image.png")
// }

// func TestUploadSPBEDocument_FailCreateUploadedFile(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	insertUploadDocumentAssessmentRepoMock = func(ctx context.Context, assessmentUploadDetail *repository.AssessmentUploadDetail) error {
// 		return nil
// 	}
// 	findOneByID = func(ctx context.Context, id string) (*repository.User, error) {
// 		return &repository.User{
// 			ID:    "test-id",
// 			Name:  "test-name",
// 			Email: "test-email",
// 			Role:  "ADMIN",
// 		}, nil
// 	}
// 	create = func(name string) (*os.File, error) {
// 		return nil, fmt.Errorf("test error")
// 	}

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockMailer := mailerMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	ctx := context.WithValue(req.Context(), userCtx.UserCtxKey, userCtx.UserCtx{
// 		ID:   "test-user-id",
// 		Role: repository.UserRole("ADMIN"),
// 	})
// 	req = req.WithContext(ctx)
// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)

// 	os.Remove("./image.png")
// }

// func TestUploadSPBEDocument_FailCopyUploadedFile(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	insertUploadDocumentAssessmentRepoMock = func(ctx context.Context, assessmentUploadDetail *repository.AssessmentUploadDetail) error {
// 		return nil
// 	}
// 	findOneByID = func(ctx context.Context, id string) (*repository.User, error) {
// 		return &repository.User{
// 			ID:    "test-id",
// 			Name:  "test-name",
// 			Email: "test-email",
// 			Role:  "ADMIN",
// 		}, nil
// 	}
// 	create = func(name string) (*os.File, error) {
// 		return &os.File{}, nil
// 	}
// 	copy = func(dst io.Writer, src io.Reader) (int64, error) {
// 		return 0, fmt.Errorf("test error")
// 	}

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockMailer := mailerMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	ctx := context.WithValue(req.Context(), userCtx.UserCtxKey, userCtx.UserCtx{
// 		ID:   "test-user-id",
// 		Role: repository.UserRole("ADMIN"),
// 	})
// 	req = req.WithContext(ctx)
// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)

// 	os.Remove("./image.png")
// }

// func TestUploadSPBEDocument_FailJsonMarshal(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	insertUploadDocumentAssessmentRepoMock = func(ctx context.Context, assessmentUploadDetail *repository.AssessmentUploadDetail) error {
// 		return nil
// 	}
// 	findOneByID = func(ctx context.Context, id string) (*repository.User, error) {
// 		return &repository.User{
// 			ID:    "test-id",
// 			Name:  "test-name",
// 			Email: "test-email",
// 			Role:  "ADMIN",
// 		}, nil
// 	}
// 	create = func(name string) (*os.File, error) {
// 		return &os.File{}, nil
// 	}
// 	copy = func(dst io.Writer, src io.Reader) (int64, error) {
// 		return 0, nil
// 	}
// 	produce = func(topic string, body []byte) error {
// 		return nil
// 	}
// 	send = func(subject, message []byte, receiver []string, templateName string, items interface{}) error {
// 		return nil
// 	}
// 	marshal = func(v any) ([]byte, error) {
// 		return nil, fmt.Errorf("test error")
// 	}

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockMailer := mailerMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	ctx := context.WithValue(req.Context(), userCtx.UserCtxKey, userCtx.UserCtx{
// 		ID:   "test-user-id",
// 		Role: repository.UserRole("ADMIN"),
// 	})
// 	req = req.WithContext(ctx)
// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)

// 	os.Remove("./image.png")
// }

// func TestUploadSPBEDocument_FailMessageQueueProduce(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	insertUploadDocumentAssessmentRepoMock = func(ctx context.Context, assessmentUploadDetail *repository.AssessmentUploadDetail) error {
// 		return nil
// 	}
// 	findOneByID = func(ctx context.Context, id string) (*repository.User, error) {
// 		return &repository.User{
// 			ID:    "test-id",
// 			Name:  "test-name",
// 			Email: "test-email",
// 			Role:  "ADMIN",
// 		}, nil
// 	}
// 	create = func(name string) (*os.File, error) {
// 		return &os.File{}, nil
// 	}
// 	copy = func(dst io.Writer, src io.Reader) (int64, error) {
// 		return 0, nil
// 	}
// 	produce = func(topic string, body []byte) error {
// 		return fmt.Errorf("test error")
// 	}
// 	send = func(subject, message []byte, receiver []string, templateName string, items interface{}) error {
// 		return nil
// 	}
// 	marshal = func(v any) ([]byte, error) {
// 		return nil, nil
// 	}

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockMailer := mailerMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	ctx := context.WithValue(req.Context(), userCtx.UserCtxKey, userCtx.UserCtx{
// 		ID:   "test-user-id",
// 		Role: repository.UserRole("ADMIN"),
// 	})
// 	req = req.WithContext(ctx)
// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)

// 	os.Remove("./image.png")
// }

// func TestUploadSPBEDocument_ErrorUserNotFound(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	insertUploadDocumentAssessmentRepoMock = func(ctx context.Context, assessmentUploadDetail *repository.AssessmentUploadDetail) error {
// 		return nil
// 	}
// 	findOneByID = func(ctx context.Context, id string) (*repository.User, error) {
// 		return nil, fmt.Errorf("test error")
// 	}
// 	create = func(name string) (*os.File, error) {
// 		return &os.File{}, nil
// 	}
// 	copy = func(dst io.Writer, src io.Reader) (int64, error) {
// 		return 0, nil
// 	}
// 	produce = func(topic string, body []byte) error {
// 		return nil
// 	}
// 	send = func(subject, message []byte, receiver []string, templateName string, items interface{}) error {
// 		return nil
// 	}
// 	marshal = func(v any) ([]byte, error) {
// 		return nil, nil
// 	}

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockMailer := mailerMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	ctx := context.WithValue(req.Context(), userCtx.UserCtxKey, userCtx.UserCtx{
// 		ID:   "test-user-id",
// 		Role: repository.UserRole("ADMIN"),
// 	})
// 	req = req.WithContext(ctx)
// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)

// 	os.Remove("./image.png")
// }

// func TestUploadSPBEDocument_FailSendEmail(t *testing.T) {
// 	pr, pw := io.Pipe()
// 	writer := multipart.NewWriter(pw)

// 	go func() {
// 		defer writer.Close()

// 		part, err := writer.CreateFormFile("supporting_document", "someimg.png")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		img := createImage()
// 		err = png.Encode(part, img)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2, err := writer.CreateFormField("institution_name")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part2.Write([]byte("Kabupaten Bandung"))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3, err := writer.CreateFormField("indicator_number")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		part3.Write([]byte("3"))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/documents/upload", pr)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	req.Header.Add("Content-Type", writer.FormDataContentType())

// 	insertUploadDocumentAssessmentRepoMock = func(ctx context.Context, assessmentUploadDetail *repository.AssessmentUploadDetail) error {
// 		return nil
// 	}
// 	findOneByID = func(ctx context.Context, id string) (*repository.User, error) {
// 		return &repository.User{
// 			ID:    "test-id",
// 			Name:  "test-name",
// 			Email: "test-email",
// 			Role:  "ADMIN",
// 		}, nil
// 	}
// 	create = func(name string) (*os.File, error) {
// 		return &os.File{}, nil
// 	}
// 	copy = func(dst io.Writer, src io.Reader) (int64, error) {
// 		return 0, nil
// 	}
// 	produce = func(topic string, body []byte) error {
// 		return nil
// 	}
// 	send = func(subject, message []byte, receiver []string, templateName string, items interface{}) error {
// 		return fmt.Errorf("test error")
// 	}
// 	marshal = func(v any) ([]byte, error) {
// 		return nil, nil
// 	}

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	mockUserRepo := userRepoMock{}
// 	mockMessageQueue := messageQueueMock{}
// 	mockFileSystemIO := fileSystemIOMock{}
// 	mockMailer := mailerMock{}
// 	mockJsonEC := jsonECMock{}
// 	mockWaClient := waClientMock{}
// 	mockConfig := config.Config{
// 		API: config.API{
// 			Host: "test-host",
// 			Port: 1,
// 		},
// 		SMTPClient: config.SMTPClient{
// 			Debug:         true,
// 			Host:          "test-host",
// 			Port:          1,
// 			AdminIdentity: "",
// 			AdminEmail:    "test@mail.com",
// 			AdminPassword: "test-password",
// 		},
// 	}

// 	ctx := context.WithValue(req.Context(), userCtx.UserCtxKey, userCtx.UserCtx{
// 		ID:   "test-user-id",
// 		Role: repository.UserRole("ADMIN"),
// 	})
// 	req = req.WithContext(ctx)
// 	r.Post("/assessments/documents/upload", UploadSPBEDocument(mockAssessmentRepo, mockUserRepo, mockMessageQueue, mockMailer, mockFileSystemIO, mockJsonEC, mockWaClient, mockConfig.API, mockConfig.SMTPClient))
// 	r.ServeHTTP(rr, req)

// 	var response UploadSpbeDocumentResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusCreated, rr.Code)

// 	os.Remove("./image.png")
// }

// func createImage() *image.RGBA {
// 	width := 200
// 	height := 100

// 	upLeft := image.Point{0, 0}
// 	lowRight := image.Point{width, height}

// 	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

// 	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
// 	cyan := color.RGBA{100, 200, 200, 0xff}

// 	// Set color for each pixel.
// 	for x := 0; x < width; x++ {
// 		for y := 0; y < height; y++ {
// 			switch {
// 			case x < width/2 && y < height/2: // upper left quadrant
// 				img.Set(x, y, cyan)
// 			case x >= width/2 && y >= height/2: // lower right quadrant
// 				img.Set(x, y, color.White)
// 			default:
// 				// Use zero value.
// 			}
// 		}
// 	}

// 	// Encode as PNG.
// 	f, _ := os.Create("image.png")
// 	png.Encode(f, img)

// 	return img

// }
