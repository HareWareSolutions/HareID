package services

import (
	"HareID/internal/enums"
	"HareID/internal/models"
	"HareID/internal/repository"
	"HareID/internal/validators"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JoinRequestServices struct {
	repo repository.Repository
	val  validators.Validations
	db   *pgxpool.Pool
}

// Serviço para criar um novo pedido de entrada em um time
func (s *JoinRequestServices) Create(ctx context.Context, requestUserID, teamID uint64) (models.JoinRequest, models.Notification, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return models.JoinRequest{}, models.Notification{}, err
	}
	defer tx.Rollback(ctx)

	team, err := s.repo.Teams.GetByID(ctx, teamID)
	if err != nil {
		log.Println(err)
		return models.JoinRequest{}, models.Notification{}, err
	}

	joinRequest := models.JoinRequest{
		SenderID:    requestUserID,
		TeamID:      teamID,
		TeamOwnerID: team.OwnerID,
		Status:      enums.PENDING, //PENDENTE
	}

	joinRequest, err = s.repo.JoinRequests.Create(ctx, tx, joinRequest)
	if err != nil {
		log.Println(err)
		return models.JoinRequest{}, models.Notification{}, err
	}

	//log.Println("created join request: ", joinRequest)

	createdNotification, err := s.repo.Notifications.CreateByJoinRequest(ctx, tx, joinRequest)
	if err != nil {
		return models.JoinRequest{}, models.Notification{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.JoinRequest{}, models.Notification{}, err
	}

	return joinRequest, createdNotification, nil

}

// Serviço para buscar todos os pedidos de entrada de um time
func (s *JoinRequestServices) GetAll(ctx context.Context, requestUserID, teamID uint64) ([]models.JoinRequest, error) {

	ok, err := s.val.JoinRequest.CanSee(ctx, requestUserID, teamID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("you dont have permission to see the team join requests")
	}

	requests, err := s.repo.JoinRequests.GetAll(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// Serviço para buscar um pedido de entrada de um time pelo ID do pedido
func (s *JoinRequestServices) GetByID(ctx context.Context, requestUserID, teamID, requestID uint64) (models.JoinRequest, error) {

	ok, err := s.val.JoinRequest.CanSee(ctx, requestID, teamID)
	if err != nil {
		return models.JoinRequest{}, err
	}

	if !ok {
		return models.JoinRequest{}, errors.New("you dont have permission to see the team join requests")
	}

	request, err := s.repo.JoinRequests.GetByID(ctx, requestID, teamID)
	if err != nil {
		return models.JoinRequest{}, err
	}

	return request, nil
}

// Serviço para deletar um pedido de entrada de um time pelo ID do pedido
func (s *JoinRequestServices) Delete(ctx context.Context, requestUserID, teamID, requestID uint64) (uint64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	ok, err := s.val.JoinRequest.CanSee(ctx, requestID, teamID)
	if err != nil {
		return 0, err
	}

	if !ok {
		return 0, errors.New("you dont have permission to delete the team join requests")
	}

	affectedRows, err := s.repo.JoinRequests.Delete(ctx, tx, requestID, teamID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (s *JoinRequestServices) Accept(ctx context.Context, requestUserID, teamID, requestID uint64) (uint64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	ok, err := s.val.JoinRequest.CanSee(ctx, requestID, teamID)
	if err != nil {
		return 0, err
	}

	if !ok {
		return 0, errors.New("you dont have permission to accept the team join requests")
	}

	request, err := s.repo.JoinRequests.GetByID(ctx, requestID, teamID)
	if err != nil {
		return 0, err
	}

	if request.Status != 0 {
		return 0, errors.New("request already accepted or rejected")
	}

	affectedRows, err := s.repo.JoinRequests.Accept(ctx, tx, requestUserID, teamID, requestID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (s *JoinRequestServices) Reject(ctx context.Context, requestUserID, teamID, requestID uint64) (uint64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	ok, err := s.val.JoinRequest.CanSee(ctx, requestID, teamID)
	if err != nil {
		return 0, err
	}

	if !ok {
		return 0, errors.New("you dont have permission to reject the requests")
	}

	request, err := s.repo.JoinRequests.GetByID(ctx, requestID, teamID)
	if err != nil {
		return 0, err
	}

	if request.Status != 0 {
		return 0, errors.New("request already accepted or rejected")
	}

	affectedRows, err := s.repo.JoinRequests.Reject(ctx, tx, requestUserID, teamID, requestID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affectedRows, nil
}
