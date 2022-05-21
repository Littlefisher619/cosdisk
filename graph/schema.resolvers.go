package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/Littlefisher619/cosdisk/graph/auth"
	"github.com/Littlefisher619/cosdisk/graph/generated"
	"github.com/Littlefisher619/cosdisk/graph/model"
)

func (r *mutationResolver) SingleUpload(ctx context.Context, input model.FileInput) (*model.File, error) {
	userid := auth.ExtractUserIdFromContext(ctx)
	err := r.Service.UploadUserFileByReader(ctx, userid, input.Path, input.File.File)
	if err != nil {
		return nil, err
	}
	return &model.File{
		Name:        input.Path,
		ContentType: input.File.ContentType,
	}, nil
}

func (r *mutationResolver) DeleteFile(ctx context.Context, path string) (string, error) {
	userid := auth.ExtractUserIdFromContext(ctx)
	err := r.Service.DeleteFIle(userid, path)
	if err != nil {
		return "failed", err
	}
	return "success", nil
}

func (r *mutationResolver) CreateDir(ctx context.Context, path string) (string, error) {
	userid := auth.ExtractUserIdFromContext(ctx)
	err := r.Service.CreateDir(userid, path)
	if err != nil {
		return "failed", err
	}
	return "success", nil
}

func (r *mutationResolver) ListDir(ctx context.Context, path string) ([]*model.File, error) {
	userid := auth.ExtractUserIdFromContext(ctx)
	res := []*model.File{}
	files, err := r.Service.ListFiles(userid, path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			res = append(res, &model.File{
				Name:        f.Name(),
				ContentType: "dir",
			})
		} else {
			res = append(res, &model.File{
				Name:        f.Name(),
				ContentType: "file",
			})
		}
	}
	return res, nil
}

func (r *mutationResolver) DeleteDir(ctx context.Context, path string) (string, error) {
	userid := auth.ExtractUserIdFromContext(ctx)
	err := r.Service.DeleteDir(userid, path)
	if err != nil {
		return "failed", err
	}
	return "success", nil
}

func (r *mutationResolver) CreateShareFile(ctx context.Context, input model.ShareInput) (string, error) {
	userid := auth.ExtractUserIdFromContext(ctx)
	return r.Service.CreateShareFile(ctx, userid, input.Path, input.ExpireDays)
}

func (r *mutationResolver) GetSharedFile(ctx context.Context, shareID string) (string, error) {
	userid := auth.ExtractUserIdFromContext(ctx)
	err := r.Service.ShareFileToUser(ctx, userid, shareID)
	if err != nil {
		return "failed", err
	}
	return "success", nil
}

func (r *mutationResolver) GetDownloadURL(ctx context.Context, path string) (string, error) {
	userid := auth.ExtractUserIdFromContext(ctx)
	return r.Service.DownloadUserFileByUrl(ctx, userid, path)
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthResult, error) {
	user, err := r.Service.UserLogin(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	token, err := r.JwtManager.GenerateToken(ctx, fmt.Sprint(user.Id))
	if err != nil {
		return nil, fmt.Errorf("token generate: %w", err)
	}

	return &model.AuthResult{
		User:  user,
		Token: token,
	}, nil
}

func (r *mutationResolver) Register(ctx context.Context, input model.CreateUserInput) (*model.AuthResult, error) {
	user, err := r.Service.UserRegister(ctx, input.Name, input.Email, input.Password)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	token, err := r.JwtManager.GenerateToken(ctx, fmt.Sprint(user.Id))
	if err != nil {
		return nil, fmt.Errorf("token generate: %w", err)
	}

	return &model.AuthResult{
		User:  user,
		Token: token,
	}, nil
}

func (r *mutationResolver) MoveFile(ctx context.Context, path string, newpath string) (string, error) {
	userid := auth.ExtractUserIdFromContext(ctx)
	str, err := r.Service.MoveFile(ctx, userid, path, newpath)
	if err != nil {
		return "failed", err
	}
	return str, nil
}

func (r *queryResolver) Empty(ctx context.Context) (string, error) {
	return "Empty", nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
