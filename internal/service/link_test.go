package service

import (
	"context"
	"errors"
	"testing"

	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/apperror"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/service/mocks"

	"go.uber.org/mock/gomock"
)

func TestLinkService_Shorten(t *testing.T) {
	type args struct {
		url        string
		customCode string
	}
	
	tests := []struct {
		name         string
		args         args
		mockOpt      func(m *mocks.MockLinkStore)
		want         string
		wantErr      error
		checkErrText string
	}{
		{
			name: "Success with custom code",
			args: args{url: "https://google.com", customCode: "my-link"},
			mockOpt: func(m *mocks.MockLinkStore) {
				m.EXPECT().Get(gomock.Any(), "my-link").Return("", nil)
				m.EXPECT().Save(gomock.Any(), "my-link", "https://google.com").Return(nil)
			},
			want:    "my-link",
			wantErr: nil,
		},
		{
			name: "Error invalid custom code characters",
			args: args{url: "https://google.com", customCode: "invalid@code!"},
			mockOpt: func(m *mocks.MockLinkStore) {},
			want:    "",
			wantErr: apperror.ErrInvalidCustomCode,
		},
		{
			name: "Error custom code already exists",
			args: args{url: "https://google.com", customCode: "busy-code"},
			mockOpt: func(m *mocks.MockLinkStore) {
				m.EXPECT().Get(gomock.Any(), "busy-code").Return("https://old.com", nil)
			},
			want:    "",
			wantErr: apperror.ErrCodeAlreadyExists,
		},
		{
			name: "Error store Get failed and wrapped with %w",
			args: args{url: "https://google.com", customCode: "some-code"},
			mockOpt: func(m *mocks.MockLinkStore) {
				m.EXPECT().Get(gomock.Any(), "some-code").Return("", errors.New("db connection timeout"))
			},
			want:         "",
			checkErrText: "service failed to check existing code: db connection timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mocks.NewMockLinkStore(ctrl)
			tt.mockOpt(mockStore)

			svc := NewLinkService(mockStore)
			got, err := svc.Shorten(context.Background(), tt.args.url, tt.args.customCode)

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Shorten() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkErrText != "" {
				if err == nil || err.Error() != tt.checkErrText {
					t.Errorf("Shorten() error text = %v, expected text %v", err, tt.checkErrText)
					return
				}
			}

			if tt.wantErr == nil && tt.checkErrText == "" && err != nil {
				t.Fatalf("Shorten() unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("Shorten() got = %v, want %v", got, tt.want)
			}
		})
	}
}
