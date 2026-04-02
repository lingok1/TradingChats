package repository

import (
	"context"
	"errors"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PromptTemplateRepository struct {
	collection *mongo.Collection
}

func NewPromptTemplateRepository(db *mongo.Database) *PromptTemplateRepository {
	return &PromptTemplateRepository{collection: db.Collection("prompt_templates")}
}

func (r *PromptTemplateRepository) Create(ctx context.Context, template *models.PromptTemplate) error {
	authCtx := models.GetAuthContext(ctx)
	tenantID := models.ResolveTenantID(authCtx, template.TenantID)
	if tenantID == "" {
		return errors.New("tenant_id is required")
	}
	template.TenantID = tenantID
	template.CreatedAt = utils.NowString()
	template.UpdatedAt = utils.NowString()
	result, err := r.collection.InsertOne(ctx, template)
	if err != nil {
		return err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		template.ID = id
	}
	return nil
}

func (r *PromptTemplateRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.PromptTemplate, error) {
	filter := bson.M{"_id": id}
	applyTenantFilter(ctx, filter)
	var template models.PromptTemplate
	if err := r.collection.FindOne(ctx, filter).Decode(&template); err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *PromptTemplateRepository) GetAll(ctx context.Context) ([]models.PromptTemplate, error) {
	filter := bson.M{}
	applyTenantFilter(ctx, filter)
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []models.PromptTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *PromptTemplateRepository) GetByTag(ctx context.Context, tag string) ([]models.PromptTemplate, error) {
	filter := bson.M{"tags": tag}
	applyTenantFilter(ctx, filter)
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []models.PromptTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *PromptTemplateRepository) Update(ctx context.Context, template *models.PromptTemplate) error {
	template.UpdatedAt = utils.NowString()
	filter := bson.M{"_id": template.ID}
	applyTenantFilter(ctx, filter)
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": template})
	return err
}

func (r *PromptTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	applyTenantFilter(ctx, filter)
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func applyTenantFilter(ctx context.Context, filter bson.M) {
	authCtx := models.GetAuthContext(ctx)
	if models.IsAdmin(authCtx) {
		return
	}
	if authCtx != nil && authCtx.TenantID != "" {
		filter["tenant_id"] = authCtx.TenantID
	}
}
