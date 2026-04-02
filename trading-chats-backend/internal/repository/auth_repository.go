package repository

import (
	"context"
	"fmt"
	"time"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthRepository struct {
	usersCollection    *mongo.Collection
	tenantsCollection  *mongo.Collection
	sessionsCollection *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) *AuthRepository {
	return &AuthRepository{
		usersCollection:    db.Collection("users"),
		tenantsCollection:  db.Collection("tenants"),
		sessionsCollection: db.Collection("user_sessions"),
	}
}

func (r *AuthRepository) EnsureIndexes(ctx context.Context) error {
	if _, err := r.usersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true).SetName("username_1"),
	}); err != nil {
		return err
	}

	if _, err := r.usersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"tenant_id": 1},
		Options: options.Index().SetName("tenant_id_1"),
	}); err != nil {
		return err
	}

	if _, err := r.tenantsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"code": 1},
		Options: options.Index().SetUnique(true).SetName("code_1"),
	}); err != nil {
		return err
	}

	_, err := r.sessionsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"access_token": 1},
			Options: options.Index().SetUnique(true).SetName("access_token_1"),
		},
		{
			Keys:    bson.M{"refresh_token": 1},
			Options: options.Index().SetUnique(true).SetName("refresh_token_1"),
		},
		{
			Keys:    bson.M{"user_id": 1},
			Options: options.Index().SetName("user_id_1"),
		},
		{
			Keys:    bson.M{"expires_at": 1},
			Options: options.Index().SetName("expires_at_1"),
		},
	})
	return err
}

func (r *AuthRepository) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	if tenant.CreatedAt == nil {
		tenant.CreatedAt = utils.NowString()
	}
	tenant.UpdatedAt = utils.NowString()
	_, err := r.tenantsCollection.InsertOne(ctx, tenant)
	return err
}

func (r *AuthRepository) UpsertTenant(ctx context.Context, tenant *models.Tenant) error {
	if tenant == nil {
		return fmt.Errorf("tenant is nil")
	}
	if tenant.ID == "" {
		tenant.ID = tenant.Code
	}
	if tenant.CreatedAt == nil {
		tenant.CreatedAt = utils.NowString()
	}
	tenant.UpdatedAt = utils.NowString()
	_, err := r.tenantsCollection.ReplaceOne(ctx, bson.M{"_id": tenant.ID}, tenant, options.Replace().SetUpsert(true))
	return err
}

func (r *AuthRepository) GetTenantByID(ctx context.Context, id string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.tenantsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&tenant); err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *AuthRepository) GetTenantByCode(ctx context.Context, code string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.tenantsCollection.FindOne(ctx, bson.M{"code": code}).Decode(&tenant); err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user.CreatedAt == nil {
		user.CreatedAt = utils.NowString()
	}
	user.UpdatedAt = utils.NowString()
	_, err := r.usersCollection.InsertOne(ctx, user)
	return err
}

func (r *AuthRepository) UpsertUser(ctx context.Context, user *models.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	if user.CreatedAt == nil {
		user.CreatedAt = utils.NowString()
	}
	user.UpdatedAt = utils.NowString()
	_, err := r.usersCollection.ReplaceOne(ctx, bson.M{"username": user.Username}, user, options.Replace().SetUpsert(true))
	return err
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.usersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.usersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UpdateUserLastLogin(ctx context.Context, userID string) error {
	_, err := r.usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"last_login_at": utils.NowString(), "updated_at": utils.NowString()}})
	return err
}

func (r *AuthRepository) SaveSession(ctx context.Context, session *models.UserSession) error {
	if session.CreatedAt == nil {
		session.CreatedAt = utils.NowString()
	}
	session.UpdatedAt = utils.NowString()
	_, err := r.sessionsCollection.InsertOne(ctx, session)
	return err
}

func (r *AuthRepository) GetSessionByAccessToken(ctx context.Context, token string) (*models.UserSession, error) {
	var session models.UserSession
	if err := r.sessionsCollection.FindOne(ctx, bson.M{"access_token": token}).Decode(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AuthRepository) GetSessionByRefreshToken(ctx context.Context, token string) (*models.UserSession, error) {
	var session models.UserSession
	if err := r.sessionsCollection.FindOne(ctx, bson.M{"refresh_token": token}).Decode(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AuthRepository) DeleteSessionByRefreshToken(ctx context.Context, token string) error {
	_, err := r.sessionsCollection.DeleteOne(ctx, bson.M{"refresh_token": token})
	return err
}

func (r *AuthRepository) DeleteExpiredSessions(ctx context.Context, now time.Time) error {
	_, err := r.sessionsCollection.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": now}})
	return err
}
