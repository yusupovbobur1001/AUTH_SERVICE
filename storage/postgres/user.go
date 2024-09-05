package postgres

import (
	pb "auth_service/genproto/auth_service"
	"auth_service/pkg/logger"
	"auth_service/storage"
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserRepository struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewUserRepository(db *pgxpool.Pool, log logger.ILogger) storage.IUserStorage {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (repo *UserRepository) GetUser(ctx context.Context, request *pb.PrimaryKey) (*pb.User, error) {
	var (
		err      error
		response pb.User
		query    string
	)
	query = `SELECT 
    id,
    user_name,
    email,
    password_hash,
    role,
    created_at,
    updated_at FROM users WHERE id=$1`
	err = repo.db.QueryRow(ctx, query, request.Id).
		Scan(&response.Id,
			&response.UserName,
			&response.Email,
			&response.PasswordHash,
			&response.Role,
			&response.CreatedAt,
			&response.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context, request *pb.GetListRequest) (*pb.GetAllUsersResponse, error) {
	var (
		err        error
		users      []*pb.User
		offset     = (request.GetLimit() - 1) * request.GetPage()
		query      string
		count      = int64(0)
		countQuery string
		// where      string
	)
	countQuery = "select * from users"
	// if request.GetSearch() != "" {
	// 	where, err := helper.MakeWherePartOfQueryWithSearchFieldOfRequest(request.GetSearch())
	// 	if err != nil {
	// 		repo.log.Error("error while taking values from search field of request in storage layer", logger.Error(err))
	// 		return &pb.GetAllUsersResponse{}, err
	// 	}
	// 	countQuery += where
	// }

	if err = repo.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		repo.log.Error("error while selecting details count in storage layer", logger.Error(err))
		return &pb.GetAllUsersResponse{}, err
	}
	query = ` select 
        id,
        user_name,
        email,
        password_hash,
        role,
        created_at,
        updated_at from users  `
	// query += where
	query += ` order by created_at DESC limit $1 offset $2`

	rows, err := repo.db.Query(ctx, query, request.GetLimit(), offset)
	if err != nil {
		repo.log.Error("error while using Query method to take details in storage layer", logger.Error(err))
		return &pb.GetAllUsersResponse{}, err
	}
	for rows.Next() {
		var user pb.User
		err = rows.Scan(&user.Id,
			&user.UserName,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			repo.log.Error("this error get_all-~~~~~~~ERROR", logger.Error(err))
			return nil, err
		}
		users = append(users, &user)

	}
	return &pb.GetAllUsersResponse{
		Users: users,
	}, nil
}

func (repo *UserRepository) DeleteUser(ctx context.Context, request *pb.PrimaryKey) (*empty.Empty, error) {
	query := ""

	query = `DELETE FROM users WHERE id=$1`
	result, err := repo.db.Exec(ctx, query)
	if err != nil {
		repo.log.Error("this error is query -~~~~~~~`>ERROR", logger.Error(err))
		return nil, err
	}
	if result.RowsAffected() == 0 {
		repo.log.Error("this user is before deleted in the past ", logger.Error(err))
		return nil, fmt.Errorf("this error is before is deleted")
	}
	return &emptypb.Empty{}, nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, request *pb.User) (*pb.User, error) {
	var (
		err      error
		query    string
		response pb.User
	)

	query = `UPDATE users SET 
                 full_name,
                 user_name,
                 email,
                 password_hash,
                 phone,
                 img_url ,
                 role,updated_at  WHERE id=$1
                 RETTURNING id,
                 full_name,
                 user_name,
                 email,
                 password_hash,
                 phone,
                 img_url,
                 role,
                 created_at,
                 updated_at`
	err = repo.db.QueryRow(ctx, query, request.Id).
		Scan(response.Id,
			response.UserName,
			response.Email,
			response.PasswordHash,
			response.Role,
			response.CreatedAt,
			response.UpdatedAt,
		)
	if err != nil {
		repo.log.Error("this error is update password -~~~~~~>ERROR", logger.Error(err))
		return nil, err
	}
	return &response, nil
}
